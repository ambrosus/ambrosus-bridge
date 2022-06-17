// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonStructs.sol";
import "./CheckReceiptsProof.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "./SignatureCheck.sol";


contract CheckPoSA is Initializable {
    uint256 private constant ADDRESS_LENGTH = 20;
    uint256 private constant EXTRA_VANITY_LENGTH = 32;
    uint256 private constant EXTRA_SEAL_LENGTH = 65;
    uint256 private constant EPOCH_LENGTH = 200;
    bytes1 constant PARENT_HASH_PREFIX = 0xA0;

    mapping(uint => mapping(address => bool)) internal allValidators;
    uint public currentEpoch;
    uint currentValidatorSetSize;

    bytes1 chainId;


    struct BlockPoSA {
        bytes3 p0Signed;
        bytes3 p0Unsigned;

        bytes32 parentHash;
        bytes p1;
        bytes32 receiptHash;
        bytes p2;
        bytes number;
        bytes p3;

        bytes p4Signed;
        bytes p4Unsigned;
        bytes extraData;

        bytes p5;
    }

    struct PoSAProof {
        BlockPoSA[] blocks;
        CommonStructs.TransferProof transfer;
        uint64 transferEventBlock;
    }


    function __CheckPoSA_init(
        address[] memory _initialValidators,
        uint _initialEpoch,
        bytes1 _chainId
    ) internal initializer {
        require(_initialValidators.length > 0, "Length of _initialValidators must be bigger than 0");

        chainId = _chainId;
        currentEpoch = _initialEpoch;
        currentValidatorSetSize = _initialValidators.length;

        for (uint i = 0; i < _initialValidators.length; i++) {
            allValidators[currentEpoch][_initialValidators[i]] = true;
        }
    }

    function checkPoSA_(PoSAProof calldata posaProof, uint minSafetyBlocks, address sideBridgeAddress) internal {
        bytes32 bareHash;
        bytes32 parentHash;
        uint finalizeVsBlock;
        uint nextVsSize;

        // posaProof can be without transfer event when we have to many vsChanges and transfer doesn't fit into proof
        if (posaProof.transferEventBlock != 0) {
            bytes32 receiptHash = calcTransferReceiptsHash(posaProof.transfer, sideBridgeAddress);
            require(posaProof.blocks[posaProof.transferEventBlock].receiptHash == receiptHash, "Transfer event validation failed");
            require(posaProof.blocks.length - posaProof.transferEventBlock >= minSafetyBlocks, "Not enough safety blocks");
        }

        for (uint i = 0; i < posaProof.blocks.length; i++) {
            BlockPoSA calldata block_ = posaProof.blocks[i];

            if (parentHash != bytes32(0))
                require(block_.parentHash == parentHash, "Wrong parent hash");

            (bareHash, parentHash) = calcBlockHash(block_);

            require(verifySignature(bareHash, getSignature(block_.extraData)), "invalid signature");

            // change validator set

            uint blockNumber = bytesToUint(block_.number);

            if (blockNumber % EPOCH_LENGTH == 0) {
                require(blockNumber / EPOCH_LENGTH == currentEpoch + 1, "invalid epoch");

                nextVsSize = newValidatorSet(block_.extraData);
                finalizeVsBlock = blockNumber + currentValidatorSetSize / 2;
            } else if (blockNumber == finalizeVsBlock) {
                currentEpoch++;
                currentValidatorSetSize = nextVsSize;

                // after finalizing vs change, next block in posaProof.blocks can have any parentHash (skipping some blocks)
                // but only if it's not the safety blocks for transfer event
                if (i < posaProof.transferEventBlock)
                    parentHash = bytes32(0);
            }
        }
    }


    function calcBlockHash(BlockPoSA calldata block_) internal view returns (bytes32, bytes32) {
        bytes memory commonRlp = abi.encodePacked(PARENT_HASH_PREFIX, block_.parentHash, block_.p1, block_.receiptHash, block_.p2, block_.number, block_.p3);
        return (
        // hash without seal (bare), for signature check
        keccak256(abi.encodePacked(block_.p0Unsigned, chainId, commonRlp, block_.p4Unsigned, getExtraDataUnsigned(block_.extraData), block_.p5)),
        // hash with seal, for prev_hash check
        keccak256(abi.encodePacked(block_.p0Signed, commonRlp, block_.p4Signed, block_.extraData, block_.p5))
        );
    }

    function getSignature(bytes calldata extraData) private pure returns (bytes calldata) {
        uint start = extraData.length - EXTRA_SEAL_LENGTH;
        return extraData[start : start + EXTRA_SEAL_LENGTH];
    }

    function getExtraDataUnsigned(bytes calldata extraData) private pure returns (bytes memory) {
        return extraData[0 : extraData.length - EXTRA_SEAL_LENGTH];
    }

    function newValidatorSet(bytes calldata extraData) private returns (uint) {
        uint nextValidatorSet = currentEpoch + 1;
        uint endPos = extraData.length - EXTRA_SEAL_LENGTH;

        uint nextValidatorSetSize;
        for (uint pos = EXTRA_VANITY_LENGTH; pos < endPos; pos += ADDRESS_LENGTH) {
            address validator = address(bytes20(extraData[pos : pos + ADDRESS_LENGTH]));
            allValidators[nextValidatorSet][validator] = true;
            nextValidatorSetSize++;
        }

        return nextValidatorSetSize;
    }

    function verifySignature(bytes32 hash, bytes calldata signature) private view returns (bool) {
        address signer = ecdsaRecover(hash, signature);
        return allValidators[currentEpoch][signer];
    }

    function bytesToUint(bytes calldata b) private pure returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }
}
