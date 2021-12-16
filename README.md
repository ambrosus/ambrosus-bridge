## Block pre-encoding

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


## Receipts proof


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
│   ├── p1 = rlpEncode(receipt1-receip5) + rlpEncode(someReceiptInfo)
│   │   └── eventData
│   └── p2 = rlpEncode(someMoreReceiptInfo) + rlpEncode(7-receip16)
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
