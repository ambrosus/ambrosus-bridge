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

    uint256 private constant MIX_HASH = 0;
    uint256 private constant DIFF_NO_TURN = 1;
    uint256 private constant DIFF_IN_TURN = 2;
    uint256 private constant EXTRA_SEAL_LENGTH = 65;

    constructor(address[] memory _initialValidators) {
        require(_initialValidators.length > 0, "Length of _initialValidators must be bigger than 0");

        validatorSet = _initialValidators;
    }

    function CheckPoSA_(bytes[] memory unsignedHeader, bytes[] memory signedHeader) public {
        require(unsignedHeader.length == signedHeader.length, "not same amount of headers");

        for (uint256 i = 0; i < unsignedHeader.length; i++) {
            block_hash = CheckBlockHeader(unsignedHeader[i], signedHeader[i]);
        }
    }

    function CheckBlockHeader(
        bytes[] memory unsignedHeader,
        bytes[] memory signedHeader
    ) internal view returns (bytes32) {
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

        return keccak256(signedHeader);
    }

    function extractSignature(bytes memory extraData) private pure returns (bytes memory) {
        return BytesLib.slice(extraData, extraData.length - EXTRA_SEAL_LENGTH, EXTRA_SEAL_LENGTH);
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
}
