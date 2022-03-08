// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../CommonStructs.sol";
import "./CheckReceiptsProof.sol";
import "./Ethash.sol";

contract CheckPoW is CheckReceiptsProof, Ethash {
    struct BlockPoW {
        bytes p1;
        bytes32 parentOrReceiptHash;
        bytes p2;
        bytes difficulty;
        bytes p3;
        bytes number;
        bytes p4;  // end when extra end
        bytes p5;  // after extra
        bytes nonce;
        bytes p6;

        uint[] dataSetLookup;
        uint[] witnessForLookup;
    }

    struct PoWProof {
        BlockPoW[] blocks;
        CommonStructs.TransferProof transfer;
    }

    function CheckPoW_(PoWProof memory powProof) public
    {

        address sideBridgeAddress = address(this);
        // todo

        bytes32 hash = CalcTransferReceiptsHash(powProof.transfer, sideBridgeAddress);
        for (uint i = 0; i < powProof.blocks.length; i++) {
            require(powProof.blocks[i].parentOrReceiptHash == hash, "parentHash or receiptHash wrong");
            hash = blockHash(powProof.blocks[i]);

            verifyEthash(powProof.blocks[i]);
        }
    }


    function blockHash(BlockPoW memory block) public view returns (bytes32) {
        // Note: too much arguments in abi.encodePacked() function cause CompilerError: Stack too deep...
        return keccak256(abi.encodePacked(
                abi.encodePacked(
                    block.p1,
                    block.parentOrReceiptHash,
                    block.p2,
                    block.difficulty,
                    block.p3
                ),
                abi.encodePacked(
                    block.number,
                    block.p4,
                    block.p5,
                    block.nonce,
                    block.p6
                )
            ));
    }

    function verifyEthash(BlockPoW memory block) public view {
        verifyPoW(
            bytesToUint(block.number),
            blockHashWithoutNonce(block),
            bytesToUint(block.nonce),
            bytesToUint(block.difficulty),
            block.dataSetLookup,
            block.witnessForLookup
        );

    }

    function blockHashWithoutNonce(BlockPoW memory block) private pure returns (bytes32) {
        bytes memory rlpHeaderHashWithoutNonce = abi.encodePacked(
            block.p1,
            block.parentOrReceiptHash,
            block.p2,
            block.difficulty,
            block.p3,
            block.number,
            block.p4,
            block.p6
        );
        rlpHeaderHashWithoutNonce[1] = rlpHeaderHashWithoutNonce[0];
        rlpHeaderHashWithoutNonce[2] = rlpHeaderHashWithoutNonce[1];

        return keccak256(rlpHeaderHashWithoutNonce);
    }


    function bytesToUint(bytes memory b) public view returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }
}
