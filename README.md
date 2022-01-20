
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

### bridge withdaw

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

3. пока вызовы transfer происходят в одном таймфрейме - трансферы просто добавляются в очередь.  
как только очередной вызов `withdraw` происходит в новом таймфрейме:
   - создается эвент `Withdraw(event_id, withdraw_queue)`
   - withdraw_queue очищается
   - event_id инкрементируется на 1

    _таймфрейм = block.timestamp / timeframe_


### relay

1. relay получает ивент `Withdraw(event_id, withdraw_queue)` c AmbBridge
2. сверяет что `event_id == EthBridge.inputEventId + 1`, иначе ищет Withdraw c подходящим event_id
3. ждет N следующих блоков (safety blocks)
4. создает receipts proof (см ниже)
5. кодирует блоки (блок с эвентом и safety) в зависимости от консенсуса сети: BlockPoA или BlockPoW (см ниже)
6. вызывает у EthBridge метод `submitTransfer`
    

### bridge submitTransfer
```
submitTransfer(
    uint event_id,
    BlockPoW[] memory blocks,
    Transfer[] memory events,
    bytes[] memory proof
)
```

1. `require(event_id == inputEventId + 1);` - проверка что эвенты приходят последовательно, без пропусков.  
`inputEventId++;`

2. считается receiptsRoot, в следующей функции проверяется что блок с таким receiptsRoot валидный 

3. вызывается _CheckPoW или _CheckPoA в зависимости от консенсуса сети с которой пришли блоки:
   - для PoA проверяется что подписан именно этот блок и подпись сделана правильным адресом (определяется по timestamp блока) 
   - для PoW проверяется что хеш блока < block.Difficulty
   
    всегда проверяется что `hash(blocks[i]) == blocks[i+1].prev_hash`

4. трансферы сохраняются в пул залоченных транзакций

### bridge unlockTransfers

1. достает из пула залоченных транзакций те, которые сохранены раньше чем `block.timestamp - lockTime`
2. переводит токены
3. удаляет из пула выполненные транзакции


## Extra


### Receipts proof


simplified example of merkle patricia tree:
```
root = hash(rlpEncode(childrens)
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
│   │   ├── eventData
│   │   └── someMoreReceiptInfo
│   ├── receipt...
│   └── receipt16
└── path3 = hash(rlpEncode(childrens)
    ├── receipt1
    ├── receipt...
    └── receipt16
```

smartcontract have `eventData` and `root` as `blocks[0].prevHashOrReceiptRoot`,  
so we need to provide other data to check if `eventData` was in trie.

simplify trie

```
root
├── p3 == path1
│   ├── p1 = rlpEncode(receipt1-receip5) + rlpEncode(someReceipt6Info)
│   │   └── eventData
│   └── p2 = rlpEncode(someMoreReceipt6Info) + rlpEncode(receipt7-receipt16)
└── p4 == path3

where

path2 = hash(p1 + eventData + p2)
root = hash(p3 + path2 + p4)

```

so

```bytes[] proof = [p1, p2, p3, p4, ...]```

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



### Block pre-encoding

```
BlockHeader = {
    ParentHash,
    UncleHash, Coinbase, Root, TxHash,
    ReceiptHash,
    Bloom, Difficulty, Number, GasLimit, GasUsed,
    Time,
    Extra, MixDigest, Nonce
}
```

```BlockHash = keccak256(rlpEncode(BlockHeader))```

rlpEncode in solidity cost much gas, so pre-encode block in relay:

```
encodedBlock = {
    p1: rlpPrefix(ParentHash)
    parentHash: ParentHash
    p2: rlpEncode(UncleHash, Coinbase, Root, TxHash) + rlpPrefix(ReceiptHash)
    receiptHash: ReceiptHash
    p3: rlpEncode(Bloom, Difficulty, Number, GasLimit, GasUsed) + rlpPrefix(Time)
    time: Time
    p4: rlpEncode(Extra, MixDigest, Nonce)
}

// parentHash, receiptHash - not encoded bytes32 values
// time - int2bytes witout leading zeros

mainBlock = {
    (p1+parentHash+p2) - bytes
    receiptHash - bytes32
    p3 - bytes
    time - bytes
    p4 - bytes
}

safetyBlock = {
    p1 - bytes
    prevHash - bytes32
    (p2+receiptHash+p3) - bytes
    time - bytes
    p3 - bytes
}

```

=>

```
struct Block {
    bytes p1;
    bytes32 prevHashOrReceiptRoot;  // receipt for main block, prevHash for others
    bytes p2;
    bytes timestamp;
    bytes p3;

    bytes signature;
}

bytes32 hash = keccak256(abi.encodePacked(
    blocks[i].p1,
    blocks[i].prevHashOrReceiptRoot,
    blocks[i].p2,
    blocks[i].timestamp,
    blocks[i].p3
));

```

