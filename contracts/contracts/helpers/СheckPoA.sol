// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract  CheckPoA {

    struct BlockPoA {
        bytes p0_seal;
        bytes p0_bare;

        bytes p1;
        bytes32 prevHashOrReceiptRoot;  // receipt for main block, prevHash for safety blocks
        bytes p2;
        bytes timestamp;
        bytes p3;

        bytes s1;
        bytes signature;
        bytes s2;
    }


//    function TestAll(BlockPoA[] memory blocks, Withdraw[] memory events, bytes[] memory proof) public {
//        bytes32 hash = calcReceiptsRoot(proof, abi.encode(events));
//
//        for (uint i = 0; i < blocks.length; i++) {
//            require(blocks[i].prevHashOrReceiptRoot == hash, "prevHash or receiptRoot wrong");
//
//            bytes rlp = abi.encodePacked(blocks[i].p1, blocks[i].prevHashOrReceiptRoot, blocks[i].p2, blocks[i].timestamp, blocks[i].p3);
//
//            // hash without seal for signature check
//            bytes32 bare_hash = keccak256(abi.encodePacked(blocks[i].p0_bare, rlp));
//            TestSignature(validator, bare_hash, blocks[i].signature);
//
//            // hash with seal, for prev_hash check
//            hash = keccak256(abi.encodePacked(blocks[i].p0_seal, rlp, blocks[i].s1, blocks[i].signature, blocks[i].s2));
//        }
//    }
}
