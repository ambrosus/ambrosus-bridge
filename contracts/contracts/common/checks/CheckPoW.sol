// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../CommonStructs.sol";
import "./CheckReceiptsProof.sol";
import "./Ethash.sol";
import "hardhat/console.sol";

contract CheckPoW is CheckReceiptsProof, Ethash  {
    struct BlockPoW {
        bytes p0withNonce;
        bytes P0withoutNonce;

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

    function CheckPoW_(PoWProof memory powProof, address sideBridgeAddress) public view
    {
        bytes32 hash = CalcTransferReceiptsHash(powProof.transfer, sideBridgeAddress);
        for (uint i = 0; i < powProof.blocks.length; i++) {
            require(powProof.blocks[i].parentOrReceiptHash == hash, "parentHash or receiptHash wrong");
            hash = blockHash(powProof.blocks[i]);

            verifyEthash(powProof.blocks[i]);
        }
    }


    function blockHash(BlockPoW memory block_) private pure returns (bytes32) {
        // Note: too much arguments in abi.encodePacked() function cause CompilerError: Stack too deep...
        return keccak256(abi.encodePacked(
                abi.encodePacked(
                    block_.p0withNonce,
                    block_.p1,
                    block_.parentOrReceiptHash,
                    block_.p2,
                    block_.difficulty,
                    block_.p3
                ),
                abi.encodePacked(
                    block_.number,
                    block_.p4,
                    block_.p5,
                    block_.nonce,
                    block_.p6
                )
            ));
    }

    function verifyEthash(BlockPoW memory block_) public view {
        verifyPoW(
            bytesToUint(block_.number),
            blockHashWithoutNonce(block_),
            bytesToUint(block_.nonce),
            bytesToUint(block_.difficulty),
            block_.dataSetLookup,
            block_.witnessForLookup
        );

    }

    function blockHashWithoutNonce(BlockPoW memory block_) private pure returns (bytes32) {
        bytes memory rlpHeaderHashWithoutNonce = abi.encodePacked(
            block_.P0withoutNonce,
            block_.p1,
            block_.parentOrReceiptHash,
            block_.p2,
            block_.difficulty,
            block_.p3,
            block_.number,
            block_.p4,
            block_.p6
        );

        return keccak256(rlpHeaderHashWithoutNonce);
    }


    function bytesToUint(bytes memory b) private pure returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }
}
