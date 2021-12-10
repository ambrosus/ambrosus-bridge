// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";


contract EthBridge {
    struct Block {
        bytes header;
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



    function TestAll(Block memory block, Withdraw[] memory events, bytes[] memory proof, bytes32 receiptsRoot) public {
        for (int i=0; i<10; i++) {
            bytes32 block_hash = keccak256(block.header);
            TestSignature(validator, block_hash, block.signature);
        }

        bytes32 events_hash = keccak256(abi.encode(events));
        TestReceiptsProof(proof, events_hash, receiptsRoot);

//        require(!TestBloom(bloom, abi.encode(events_hash)), "Failed to verify bloom");

        for (uint i=0; i<events.length; i++) {
            emit WithdrawEvent(events[i].fromAddress, events[i].toAddress, events[i].amount);
        }

    }



    function TestReceiptsProof(bytes[] memory proof, bytes32 eventToSearch, bytes32 receiptsRoot) public {
        bytes32 el = eventToSearch;
        bytes memory s;

        for (uint i = 0; i < proof.length-1; i += 2) {
            s = abi.encodePacked(proof[i], el, proof[i+1]);
            el = (s.length > 32) ? keccak256(s) : bytes32(s);
        }

        require(el == receiptsRoot, "Failed to verify receipts proof");
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


}
