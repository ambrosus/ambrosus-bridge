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
        for (uint i = 0; i < _initialValidators.length; i++) {
            allValidators[currentValidatorSet][_initialValidators[i]] = true;
        }

    }

    function CheckPoSA_(PoSAProof calldata posaProof, address sideBridgeAddress) internal {
        bytes32 receiptHash = calcTransferReceiptsHash(posaProof.transfer, sideBridgeAddress);
        require(posaProof.blocks[posaProof.transferEventBlock].receiptHash == receiptHash, "Transfer event validation failed");

        for (uint i = 0; i < posaProof.blocks.length; i++) {
            CheckBlock(posaProof.blocks[i]);
            // todo check parentHash
        }
    }

    function CheckBlock(BlockPoSA calldata block_) private {
        (bytes32 bareHash, bytes32 sealHash) = calcBlockHash(block_);

        require(verifySignature(bareHash, getSignature(block_.extraData)), "invalid signature");

        uint blockNumber = bytesToUint(block_.number);


        if (blockNumber % EPOCH_LENGTH == 0) {
            // todo verifySignature will fail if we do this now
            currentValidatorSet++;

            address[] memory newValidators = getValidatorSet(block_.extraData);
            for (uint i = 0; i < newValidators.length; i++) {
                allValidators[currentValidatorSet][newValidators[i]] = true;
            }
        }
        // todo finalize ValidatorSet
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
        return extraData[start:start + EXTRA_SEAL_LENGTH];
    }

    function getExtraDataUnsigned(bytes calldata extraData) private pure returns(bytes memory) {
        return extraData[0:EXTRA_VANITY_LENGTH];
    }

    function getValidatorSet(bytes calldata extraData) private pure returns (address[] memory) {
        uint256 currentPosition = EXTRA_VANITY_LENGTH;
        uint256 endPosition = extraData.length - EXTRA_SEAL_LENGTH;
        uint256 numValidators = (endPosition - currentPosition) / ADDRESS_LENGTH;

        address[] memory validators = new address[](numValidators);

        for (uint256 i = 0; i < numValidators; i++) {
            validators[i] = bytesToAddress(extraData[currentPosition:currentPosition + ADDRESS_LENGTH]);

            currentPosition += ADDRESS_LENGTH;
        }

        return validators;
    }

    function verifySignature(bytes32 hash, bytes memory signature) private view returns (bool) {
        address signer = getSigner(hash, signature);
        return allValidators[currentValidatorSet][signer];
    }

    function bytesToUint(bytes memory b) private pure returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }

    function bytesToAddress(bytes memory _bytes) private pure returns (address) {
        address tempAddress;
        assembly {
            tempAddress := div(mload(add(_bytes, 0x20)), 0x1000000000000000000000000)
        }
        return tempAddress;
    }

    function getSigner(bytes32 messageHash, bytes memory signature) internal pure returns(address) {
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
