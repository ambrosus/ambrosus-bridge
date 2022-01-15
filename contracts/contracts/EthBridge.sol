// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";


contract EthBridge {
    struct Block {
        bytes p1;
        bytes32 prevHashOrReceiptRoot;  // receipt for main block, prevHash for others
        bytes p2;
        bytes timestamp;
        bytes p3;

        bytes signature;
    }

    struct Withdraw {
        address fromAddress;
        address toAddress;
        uint amount;
    }

    bytes1 constant b1 = bytes1(0x01);

    address validator;

    event WithdrawEvent(address indexed from, address indexed to, uint amount);

    constructor(address validator_) {
        validator = validator_;
    }


    function TestAll(Block[] memory blocks, Withdraw[] memory events, bytes[] memory proof) public {
        TestReceiptsProof(proof, abi.encode(events), blocks[0].prevHashOrReceiptRoot);

        bytes32 hash = calcReceiptsRoot(proof, abi.encode(events));

        for (uint i = 0; i < blocks.length; i++) {
            require(blocks[i].prevHashOrReceiptRoot == hash, "prevHash or receiptRoot wrong");
            hash = keccak256(abi.encodePacked(blocks[i].p1, blocks[i].prevHashOrReceiptRoot, blocks[i].p2, blocks[i].timestamp, blocks[i].p3));

            TestSignature(validator, hash, blocks[i].signature);
        }


        //        require(!TestBloom(bloom, abi.encode(events_hash)), "Failed to verify bloom");

        for (uint i = 0; i < events.length; i++) {
            emit WithdrawEvent(events[i].fromAddress, events[i].toAddress, events[i].amount);
        }

    }

    function withdraw(address toAddress, uint amount) public {
        emit WithdrawEvent(msg.sender, toAddress, amount);

        // ...
    }

    function TestReceiptsProof(bytes[] memory proof, bytes memory eventToSearch, bytes32 receiptsRoot) public {
        require(calcReceiptsRoot(proof, eventToSearch) == receiptsRoot, "Failed to verify receipts proof");
    }

    function calcReceiptsRoot(bytes[] memory proof, bytes memory eventToSearch) public view returns (bytes32){
        bytes32 el = keccak256(abi.encodePacked(proof[0], eventToSearch, proof[1]));
        bytes memory s;

        for (uint i = 2; i < proof.length - 1; i += 2) {
            s = abi.encodePacked(proof[i], el, proof[i + 1]);
            el = (s.length > 32) ? keccak256(s) : bytes32(s);
        }

        return el;
    }


    function TestSignature(address signer, bytes32 message, bytes memory signature) internal view {
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
        }
        require(
            ecrecover(keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message)), v, r, s) == signer,
            "Failed to verify sign");
    }


    function TestBloom(bytes memory bloom, bytes memory topicHash) public returns (bool) {
        bytes32 hashbuf = keccak256(topicHash);

        // todo asm
        bytes1 v1 = b1 << uint8(hashbuf[1] & 0x07);
        bytes1 v2 = b1 << uint8(hashbuf[3] & 0x07);
        bytes1 v3 = b1 << uint8(hashbuf[5] & 0x07);

        uint i1 = 256 - uint((Uint16(hashbuf[0], hashbuf[1]) & 0x7ff) >> 3) - 1;
        uint i2 = 256 - uint((Uint16(hashbuf[2], hashbuf[3]) & 0x7ff) >> 3) - 1;
        uint i3 = 256 - uint((Uint16(hashbuf[4], hashbuf[5]) & 0x7ff) >> 3) - 1;

        return
        v1 == v1 & bloom[i1] &&
        v2 == v2 & bloom[i2] &&
        v3 == v3 & bloom[i3];
    }


    function Uint16(bytes2 a, bytes1 b) internal view returns (uint16) {
        return uint16(a) + uint16(uint8(b));
    }

    function bytesToUint(bytes memory b) public view returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }
}
