// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./CheckReceiptsProof.sol";
import "./CheckPoW_Ethash.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";


contract CheckPoW is Initializable, Ethash {
    struct BlockPoW {
        bytes3 p0WithNonce;
        bytes3 p0WithoutNonce;

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

    uint256 minimumDifficulty;

    function __CheckPoW_init(
        uint256 minimumDifficulty_
    ) internal initializer {
        minimumDifficulty = minimumDifficulty_;
    }

    /*
     PoWProof.blocks contains:
      - block with transfer event;
      - safety blocks for transfer event

      Function will check all blocks, checking it pow hash.
      Each block parentHash must be equal to the hash of the previous block.
      If there are no errors, the transfer is considered valid
    */
    function checkPoW_(PoWProof calldata powProof, address sideBridgeAddress) internal view
    {
        bytes32 hash = calcTransferReceiptsHash(powProof.transfer, sideBridgeAddress);
        for (uint i = 0; i < powProof.blocks.length; i++) {
            require(powProof.blocks[i].parentOrReceiptHash == hash, "parentHash or receiptHash wrong");
            hash = blockHash(powProof.blocks[i]);

            verifyEthash(powProof.blocks[i]);
        }
    }


    function verifyEthash(BlockPoW calldata block_) internal view {
        uint difficulty = bytesToUint(block_.difficulty);
        require(difficulty >= minimumDifficulty, "difficulty too low");
        verifyPoW(
            bytesToUint(block_.number),
            blockHashWithoutNonce(block_),
            bytesToUint(block_.nonce),
            difficulty,
            block_.dataSetLookup,
            block_.witnessForLookup
        );
    }

    function blockHash(BlockPoW calldata block_) internal pure returns (bytes32) {
        // Note: too much arguments in abi.encodePacked() function cause CompilerError: Stack too deep...
        return keccak256(abi.encodePacked(
                abi.encodePacked(
                    block_.p0WithNonce,
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

    function blockHashWithoutNonce(BlockPoW calldata block_) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(
            abi.encodePacked(
                block_.p0WithoutNonce,
                block_.p1,
                block_.parentOrReceiptHash,
                block_.p2
            ),
            abi.encodePacked(
                block_.difficulty,
                block_.p3,
                block_.number,
                block_.p4,
                block_.p6
            )
        ));
    }


    function bytesToUint(bytes calldata b) private pure returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }

    uint256[15] private ___gap;
}
