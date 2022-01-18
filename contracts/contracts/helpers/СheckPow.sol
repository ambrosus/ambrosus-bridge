// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract CheckPoW {

    struct BlockPoW {
        bytes p1;
        bytes32 prevHashOrReceiptRoot;  // receipt for main block, prevHash for safety blocks
        bytes p2;
        bytes difficulty;
        bytes p3;
    }

    function TestPoW(
        Block[] memory blocks,
        Transfer[] memory events,
        bytes[] memory proof) public
    {
        bytes32 hash = calcReceiptsRoot(proof, abi.encode(events));

        for (uint i = 0; i < blocks.length; i++) {
            require(blocks[i].prevHashOrReceiptRoot == hash, "prevHash or receiptRoot wrong");
            hash = keccak256(abi.encodePacked(blocks[i].p1, blocks[i].prevHashOrReceiptRoot, blocks[i].p2, blocks[i].difficulty, blocks[i].p3));

            TestPoW(hash, blocks[i].difficulty);
        }

    }






}
