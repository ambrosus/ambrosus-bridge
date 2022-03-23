// SPDX-License-Identifier: MIT
pragma solidity >=0.7.3 <0.9.0;

import "../CommonStructs.sol";
import "./libs/RLPReader.sol";
import "./libs/BytesLib.sol";
import "./libs/ECDSA.sol";
import "./CheckReceiptsProof.sol";

contract CheckPoSA is CheckReceiptsProof {
    using RLPReader for RLPReader.RLPItem;
    using RLPReader for RLPReader.Iterator;
    using RLPReader for bytes;

    address[] private validatorSet;

    bytes32 private longestChainEndpoint;
    uint256 private currentHeight;

    uint256 private constant MIX_HASH = 0;
    uint256 private constant DIFF_NO_TURN = 1;
    uint256 private constant DIFF_IN_TURN = 2;
    uint256 private constant ADDRESS_LENGTH = 20;
    uint256 private constant EXTRA_VANITY_LENGTH = 32;
    uint256 private constant EXTRA_SEAL_LENGTH = 65;
    uint256 private constant EPOCH_LENGTH = 200;

    constructor(address[] memory _initialValidators) {
        require(_initialValidators.length > 0, "Length of _initialValidators must be bigger than 0");

        validatorSet = _initialValidators;
    }

    function CheckPoSA_(bytes[] memory unsignedHeader, bytes[] memory signedHeader) public {
        require(unsignedHeader.length == signedHeader.length, "not same amount of headers");

        for (uint256 i = 0; i < unsignedHeader.length; i++) {
            CheckBlockHeader(unsignedHeader[i], signedHeader[i]);
        }
    }

    function CheckBlockHeader(bytes[] memory unsignedHeader, bytes[] memory signedHeader) public {
        RLPReader.RLPItem[] memory unsignedHeaderItems = unsignedHeader
            .toRlpItem()
            .toList();

        RLPReader.RLPItem[] memory signedHeaderItems = signedHeader
            .toRlpItem()
            .toList();
        
        require(signedHeaderItems[13].toUint() == MIX_HASH, "mixHash not 0");

        require(signedHeaderItems[7].toUint() == DIFF_IN_TURN || signedHeaderItems[7].toUint() == DIFF_NO_TURN,
            "difficulty not 1 or 2"
        );

        require(compareBlockHeader(unsignedHeaderItems, signedHeaderItems), "unsigned not equals signed");

        bytes memory signature = extractSignature(signedHeaderItems[12].toBytes());

        require(verifySignature(keccak256(unsignedHeader), signature), "invalid signature");

        bytes32 blockHash = keccak256(signedHeader);
        uint256 blockNumber = signedHeaderItems[8].toUint();

        require(blockHashes[blockHash] == 0, "block already exists");

        bytes32 parentHash = bytes32(signedHeaderItems[0].toUint());

        require(blockHashes[parentHash] != 0, "no parent");

        if (blockNumber % EPOCH_LENGTH == 0) {
            validatorSets[blockHash] = extractValidatorSet(signedHeaderItems[12].toBytes());
        } else if (blockNumber % EPOCH_LENGTH == validatorSet.length / 2) {
            bytes32 epochBlockHash = parentHash;

            for (uint256 i = 0; i < validatorSet.length / 2 - 1; i++) {
                epochBlockHash = blockHashes[epochBlockHash];
            }

            validatorSet = validatorSets[epochBlockHash];
        }

        blockHashes[blockHash] = parentHash;

        if (parentHash == longestChainEndpoint || currentHeight == blockNumber) {
            longestChainEndpoint = blockHash;
            currentHeight = blockNumber;
        }
    }

    function compareBlockHeader(
        RLPReader.RLPItem[] memory unsignedHeader,
        RLPReader.RLPItem[] memory signedHeader
    ) private pure returns (bool) {
        if (unsignedHeader.length != signedHeader.length + 1) return false;

        for (uint256 i = 0; i < signedHeader.length; i++) {
            if (i == 12) {
                bytes memory extraDataSignedHeader = signedHeader[i].toBytes();
                uint256 signatureStart = extraDataSignedHeader.length - EXTRA_SEAL_LENGTH;
                bytes memory extraDataUnsigned = BytesLib.slice(extraDataSignedHeader, 0, signatureStart);
                bytes memory extraDataUnSignedHeader = unsignedHeader[i + 1].toBytes();

                if (keccak256(extraDataUnsigned) != keccak256(extraDataUnSignedHeader)) return false;
            } else {
                if (unsignedHeader[i + 1].rlpBytesKeccak256() != signedHeader[i].rlpBytesKeccak256()) return false;
            }
        }

        return true;
    }

    function extractSignature(bytes memory extraData) private pure returns (bytes memory) {
        return BytesLib.slice(extraData, extraData.length - EXTRA_SEAL_LENGTH, EXTRA_SEAL_LENGTH);
    }

    function extractValidatorSet(bytes memory extraData) private pure returns (address[] memory) {
        uint256 currentPosition = EXTRA_VANITY_LENGTH;
        uint256 endPosition = extraData.length - EXTRA_SEAL_LENGTH;
        uint256 numValidators = (endPosition - currentPosition) / ADDRESS_LENGTH;

        address[] memory validators = new address[](numValidators);

        for (uint256 i = 0; i < numValidators; i++) {
            validators[i] = BytesLib.toAddress(BytesLib.slice(extraData, currentPosition, ADDRESS_LENGTH), 0);

            currentPosition += ADDRESS_LENGTH;
        }

        return validators;
    }

    function verifySignature(bytes32 hash, bytes memory signature) private view returns (bool) {
        address signer = ECDSA.recover(hash, signature);

        for (uint256 i = 0; i < validatorSet.length; i++) {
            if (signer == validatorSet[i]) {
                return true;
            }
        }

        return false;
    }

    function getCurrentHeight() public view returns (uint256) {
        return currentHeight;
    }
}
