
![uml](./docs/output/classes.png)
Smart contracts structure


## Flow

![uml](./docs/output/flow.png)

### frontend

1. юзер заходит на фронт
2. фронт берет список всех токенов, их иконки и т.д.  
(avax берет отсюда https://raw.githubusercontent.com/ava-labs/avalanche-bridge-resources/main/token_list.json)
3. фронт делает запросы на смарт контракты бриджей, проверяя что токен по адресу существует и не выключен
4. юзер вызывает `withdraw(tokenAddress, toAddress, amount, {value: fee})` у контракта бриджа в сети, из которой он хочет вывести деньги

### bridge withdraw

`withdraw(address tokenAmbAddress, address toAddress, uint amount)`

1. `require(msg.value == fee)`

2. информация о выводе добавляется в очередь `Transfer[] withdraw_queue`
    ```
    struct Transfer {
        address tokenAddress;
        address toAddress;
        uint amount;
    }
    ```

3. пока вызовы `withdraw` происходят в одном таймфрейме - трансферы просто добавляются в очередь.  
как только очередной вызов `withdraw` происходит в новом таймфрейме:
   - создается эвент `Transfer(event_id, withdraw_queue)`
   - withdraw_queue очищается
   - event_id инкрементируется на 1

    _таймфрейм = block.timestamp / timeframe_


### relay

1. relay получает ивент `Transfer(event_id, withdraw_queue)` c AmbBridge
2. сверяет что `event_id == EthBridge.inputEventId + 1`, иначе ищет Transfer c подходящим event_id
3. ждет N следующих блоков (safety blocks)
4. создает receipts proof (см ниже)
5. кодирует блоки (блок с эвентом и safety) в зависимости от консенсуса сети: BlockPoA или BlockPoW (см ниже)
6. вызывает у EthBridge метод `submitTransfers`
    

### bridge submitTransfers
```
submitTransfers(
    uint event_id,
    BlockPoA[] memory blocks,
    Transfer[] memory events,
    bytes[] memory proof
)
```

1. `require(event_id == inputEventId + 1);` - проверка что эвенты приходят последовательно, без пропусков.  
`inputEventId++;`

2. считается receiptsRoot, в следующей функции проверяется что блок с таким receiptsRoot валидный 

3. вызывается _CheckPoW или _CheckPoA в зависимости от консенсуса сети с которой пришли блоки:
   - для PoA проверяется что подписан именно этот блок и подпись сделана правильным адресом (определяется по полю step) 
   - PoW проверяется с использованием ethash (_coming soon_)
   
    неявно проверяется что `hash(blocks[i]) == blocks[i+1].prev_hash`

4. трансферы сохраняются в пул залоченных транзакций

### bridge unlockTransfers

1. достает из пула залоченных транзакций те, которые сохранены раньше чем `block.timestamp - lockTime`
2. переводит токены
3. удаляет из пула выполненные транзакции


## Extra


### Block pre-encoding

Что бы доказать что блок правильный нужно считать его хеш.

```
blockHeader = {
    ParentHash, UncleHash, Coinbase, Root, TxHash,
    ReceiptHash, 
    Bloom, Difficulty, Number, GasLimit, GasUsed,
    Time, Extra, MixDigest, Nonce
}

blockHash = keccak256(rlpEncode(blockHeader))
assert blockHash == needHash
```

Хешируемое значение почти всегда сначала кодируется в RLP (Recursive Length Prefix).   
Эта кодировка только добавляет к входному значению префикс, а значит входное значение с каким то смещением содержится в выходном.  
=> `rlpEncode(value) = rlpPrefix(value) + value`, где + означает конкатенацию байт.


Для экономии газа relay будет подготавливать блоки для смарт контракта таким образом, 
что бы вместо `rlpEncode` использовать конкатенацию (`abi.encodePacked`).

Например, relay может разбить `rlpEncode(header)` по разделителю `receiptRoot`  
`rlpParts := bytes.Split(rlpHeader, receiptRoot)`, тогда  
`keccak256(abi.encodePacked(rlpParts[0], receiptRoot, rlpParts[1]) == needHash`


### Receipts proof (todo)


simplified example of merkle patricia tree:
```
receiptsRoot = hash(rlpEncode(childrens)
├── path1 = hash(rlpEncode(childrens)
│   ├── receipt1
│   ├── receipt2
│   ├── receipt...
│   └── receipt16
├── path2 = hash(rlpEncode(childrens)
│   ├── receipt1
│   ├── receipt2
│   ├── receipt...
│   ├── receipt6
│   │   ├── someReceiptInfo
│   │   ├── **eventData**
│   │   └── someMoreReceiptInfo
│   ├── receipt...
│   └── receipt16
└── path3 = hash(rlpEncode(childrens)
    ├── receipt1
    ├── receipt...
    └── receipt16
```

smart contract have `eventData` as function argument and `receiptsRoot` as `blocks[0].prevHashOrReceiptRoot`,  
so we need to provide other data to check if `eventData` was in trie.

simplify trie

```
receiptsRoot = hash(path1 + path2 + path3)
├── path1 = hash(rlpEncode(childrens)
├── path2 = hash(rlpEncode(p1 + eventData + p2)
│   ├── p1 = rlpEncode(receipt1-receip5) + someReceipt6InfoRLP
│   ├───── eventData
│   └── p2 = someMoreReceipt6InfoRLP + rlpEncode(receipt7-receipt16)
└── path3 = hash(rlpEncode(childrens)

```

=>

path2 = hash(p1 + eventData + p2)  
receiptsRoot = hash(path1 + path2 + path3)

=>

```bytes[] proof = [p1, p2, path1, path3, ...]```

check

```
hash(abi.encodePacked(
    p5,
    hash(abi.encodePacked(
        p3, 
        hash(abi.encodePacked(
            p1,
            eventData,
            p2
        )),
        p4
    )),
    p6
) == header.receiptsRoot
```


### Fees API
#### Endpoint: /fees
#### Method: POST

#### Request body params:
- `tokenAddress` - string, hex address of token in **from** network
- `isAmb` - bool, is the **from** network is **AMB**
- `amount` - string, amount of tokens in hex
- `isAmountWithFees` - bool, is the amount includes fees (used with "set max" button on the frontend)

#### Response params:
- `bridgeFee` - string, bridge fee in hex
- `transferFee` - string, transfer fee in hex
- `signature` - string, signature of the data

#### Examples:
Request body:
- **URL**: http://localhost:8080/fees
- **Body**:
```json
{
  "tokenAddress": "0xc778417E063141139Fce010982780140Aa0cD5Ab",
  "isAmb": true,
  "amount": "0xDE0B6B3A7640000",
  "isAmountWithFees": false
}
```

- Success request:
    - **Status code**: 200
    - **Result**:
  ```json
  {
    "bridgeFee": "0xbc4b381d188000",
    "transferFee": "0xe8d4a51000",
    "signature": "0x6105ca999d43b1f1182d4955f5706e8bf27097b8cb80da35c04016238e2adff91e38f275d640911a46b208d0bb7239d83d5e427b7b93b0cd7390034d724bdb0500"
  }
  ```

- Failure request (wrong request body):
    - **Status code**: 400
    - **Result**:
  ```json
  {
    "message": "error when decoding request body",
    "developerMessage": "якась помилка"
  }
  ```

- Failure request (internal error):
    - **Status code**: 500
    - **Result**:
  ```json
  {
    "message": "error when getting bridge fee",
    "developerMessage": "якась помилка"
  }
  ```
  
- Failure request (when `isAmountWithFees` is true and `amount` is too small):
    - **Status code**: 500
    - **Result**:
  ```json
  {
    "message": "amount is too small"
  }
  ```
