
![uml](./docs/output/classes.png)
Smart contracts structure


## Flow

![uml](./docs/output/flow.png)

### Frontend

1. The user goes to the frontend
2. The front-end takes a list of all tokens, their icons, etc. 
(avax берет отсюда https://raw.githubusercontent.com/ava-labs/avalanche-bridge-resources/main/token_list.json)
3. The front-end makes queries to bridge smart contracts, checking that the token at the address exists and is not disabled
4. the user calls `withdraw(tokenAddress, toAddress, amount, {value: fee})` the bridge contract in the network from which he wants to withdraw money

### Bridge withdraw

`withdraw(address tokenAmbAddress, address toAddress, uint amount)`

1. `require(msg.value == fee)`

2. the output information is added to the queue `Transfer[] withdraw_queue`
    ```
    struct Transfer {
        address tokenAddress;
        address toAddress;
        uint amount;
    }
    ```

3. As long as `withdraw` calls occur in one timeframe - transfers are simply added to the queue.  
As soon as the next call of `withdraw` occurs in a new timeframe:
   - Event creation `Transfer(event_id, withdraw_queue)`
   - withdraw_queue cleared
   - event_id incremented by 1

    _timeframe = block.timestamp / timeframe_


### Relay

1. Relay get event `Transfer(event_id, withdraw_queue)` from AmbBridge
2. checks that `event_id == EthBridge.inputEventId + 1`, otherwise look for Transfer with a matching event_id
3. waits for N next safety blocks
4. creates receipts proof (see below)
5. encodes the blocks (event and safety block) depending on the network consensus: BlockPoA or BlockPoW (see below)
6. calls `submitTransfers` method of EthBridge
    

### Bridge submitTransfers
```
submitTransfers(
    uint event_id,
    BlockPoA[] memory blocks,
    Transfer[] memory events,
    bytes[] memory proof
)
```

1. `require(event_id == inputEventId + 1);` - check that the events come consistently, without skipping.  
`inputEventId++;`

2. is considered receiptsRoot, the next function checks that the block with this receiptsRoot is valid  

3. is called _CheckPoW or _CheckPoA, depending on the consensus of the network from which the blocks came:
   - for PoA, it checks that this particular block is signed and that the signature is made with the correct address (determined by the step field)
   - PoW is checked using ethash (_coming soon_)
   
    it is implicitly checked that `hash(blocks[i]) == blocks[i+1].prev_hash`

4. transfers are saved to the pool of blocked transactions

### Bridge unlockTransfers

1. retrieves from the pool of blocked transactions those saved before `block.timestamp - lockTime`.
2. transfers tokens
3. deletes executed transactions from the pool


## Extra


### Block pre-encoding

To prove that the block is correct, you need to read its hash.

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

The hashed value is almost always encoded in RLP (Recursive Length Prefix) first.   
This encoding only adds a prefix to the input value, which means that the input value with some offset is contained in the output value.  
=> `rlpEncode(value) = rlpPrefix(value) + value`, where + means concatenation of bytes.


To save gas, relay will prepare blocks for a smart contract in this way, 
to use concatenation (`abi.encodePacked`) instead of `rlpEncode`.

For example, relay can split `rlpEncode(header)` by the delimiter `receiptRoot`  
`rlpParts := bytes.Split(rlpHeader, receiptRoot)`, then  
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
- `amount` - string, used amount in calculation fees (useful when param `isAmountWithFees` has used)
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
    "amount": "0xDE0B6B3A7640000",
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
