// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./CommonStructs.sol";

contract CheckPoW {
    struct BlockPoW {
        bytes p1;
        bytes32 prevHashOrReceiptRoot;  // receipt for main block, prevHash for safety blocks
        bytes p2;
        bytes difficulty;
        bytes p3;
    }

    function CheckPoW_(
        BlockPoW[] memory blocks,
        CommonStructs.Transfer[] memory events,
        bytes[] memory proof) public
    {
//        bytes32 hash = calcReceiptsRoot(proof, abi.encode(events));

        for (uint i = 0; i < blocks.length; i++) {
//            require(blocks[i].prevHashOrReceiptRoot == hash, "prevHash or receiptRoot wrong");
            bytes32 hash = keccak256(abi.encodePacked(blocks[i].p1, blocks[i].prevHashOrReceiptRoot, blocks[i].p2, blocks[i].difficulty, blocks[i].p3));
        }
    }

    function testCheckPow(BlockPoW memory block) public {
        bytes32 hash = keccak256(abi.encodePacked(block.p1, block.prevHashOrReceiptRoot, block.p2, block.difficulty, block.p3));
        require(uint(hash) < bytesToUint(block.difficulty), "hash must be less than difficulty");
    }

    function bytesToUint(bytes memory b) public view returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }
}
