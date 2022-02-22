// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../CommonStructs.sol";

contract CheckPoW {
    struct BlockPoW {
        bytes p1;
        bytes32 prevHashOrReceiptRoot;  // receipt for main block, prevHash for safety blocks
        bytes p2;
        bytes difficulty;
        bytes p3;
    }

    struct PoWProof {
        BlockPoW[] blocks;
        CommonStructs.TransferProof transfer;
    }

    function CheckPoW_(PoWProof memory powProof) public
    {

    }

//    function testCheckPow(BlockPoW memory block) public {
//        bytes32 hash = keccak256(abi.encodePacked(block.p1, block.prevHashOrReceiptRoot, block.p2, block.difficulty, block.p3));
//        require(uint(hash) < bytesToUint(block.difficulty), "hash must be less than difficulty");
//    }
//
//    function bytesToUint(bytes memory b) public view returns (uint){
//        return uint(bytes32(b)) >> (256 - b.length * 8);
//    }
}
