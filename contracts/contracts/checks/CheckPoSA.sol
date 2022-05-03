// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonStructs.sol";
import "./CheckReceiptsProof.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";


contract CheckPoSA is Initializable, CheckReceiptsProof {
    uint256 private constant ADDRESS_LENGTH = 20;
    uint256 private constant EXTRA_VANITY_LENGTH = 32;
    uint256 private constant EXTRA_SEAL_LENGTH = 65;
    uint256 private constant EPOCH_LENGTH = 200;
    bytes1 constant PARENT_HASH_PREFIX = 0xA0;

    mapping(uint => mapping(address => bool)) private allValidators;
    uint currentValidatorSet;
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
        currentValidatorSet = _initialEpoch;
        currentValidatorSetSize = _initialValidators.length;

        for (uint i = 0; i < _initialValidators.length; i++) {
            allValidators[currentValidatorSet][_initialValidators[i]] = true;
        }
    }

    function checkPoSA_(PoSAProof calldata posaProof, address sideBridgeAddress) internal {
        bytes32 bareHash;
        bytes32 sealHash;
        uint finalizeVsBlock;
        uint nextVsSize;

        bytes32 receiptHash = calcTransferReceiptsHash(posaProof.transfer, sideBridgeAddress);
        require(posaProof.blocks[posaProof.transferEventBlock].receiptHash == receiptHash, "Transfer event validation failed");

        for (uint i = 0; i < posaProof.blocks.length; i++) {
            BlockPoSA calldata block_ = posaProof.blocks[i];
            (bareHash, sealHash) = calcBlockHash(block_);

            require(verifySignature(bareHash, getSignature(block_.extraData)), "invalid signature");


            uint blockNumber = bytesToUint(block_.number);

            if (blockNumber % EPOCH_LENGTH == 0) {
                nextVsSize = newValidatorSet(block_.extraData);
                finalizeVsBlock = blockNumber + currentValidatorSetSize / 2;
            }
            else if (blockNumber == finalizeVsBlock) {
                currentValidatorSet++;
                currentValidatorSetSize = nextVsSize;
            }

            if (i + 1 != posaProof.blocks.length) {
                require(sealHash == posaProof.blocks[i + 1].parentHash, "wrong parent hash");
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

    function getSignature(bytes calldata extraData) private pure returns (bytes memory) {
        uint start = extraData.length - EXTRA_SEAL_LENGTH;
        return extraData[start : start + EXTRA_SEAL_LENGTH];
    }

    function getExtraDataUnsigned(bytes calldata extraData) private pure returns (bytes memory) {
        return extraData[0 : extraData.length - EXTRA_SEAL_LENGTH];
    }

    function newValidatorSet(bytes calldata extraData) private returns(uint) {
        uint nextValidatorSet = currentValidatorSet + 1;
        uint endPos = extraData.length - EXTRA_SEAL_LENGTH;

        uint nextValidatorSetSize;
        for (uint pos = EXTRA_VANITY_LENGTH; pos < endPos; pos += ADDRESS_LENGTH) {
            address validator = address(bytes20(extraData[pos : pos + ADDRESS_LENGTH]));
            allValidators[nextValidatorSet][validator] = true;
            nextValidatorSetSize++;
        }

        return nextValidatorSetSize;
    }

    function verifySignature(bytes32 hash, bytes memory signature) private view returns (bool) {
        address signer = getSigner(hash, signature);
        return allValidators[currentValidatorSet][signer];
    }

    function bytesToUint(bytes memory b) private pure returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }

    function getSigner(bytes32 messageHash, bytes memory signature) internal pure returns (address) {
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
            if lt(v, 27) {v := add(v, 27)}
        }
        return ecrecover(messageHash, v, r, s);
    }
}
