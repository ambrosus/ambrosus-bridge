// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonStructs.sol";
import "./CheckReceiptsProof.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";


contract CheckPoSA is Initializable, CheckReceiptsProof {
    mapping(uint => mapping(address => bool)) private allValidators;
    uint currentValidatorSet;

    uint256 private currentHeight;
    bytes32 private genesisBlockHash;

    uint256 private constant ADDRESS_LENGTH = 20;
    uint256 private constant EXTRA_VANITY_LENGTH = 32;
    uint256 private constant EXTRA_SEAL_LENGTH = 65;
    uint256 private constant EPOCH_LENGTH = 200;


    struct BlockPoSA {
        bytes p0Signed;
        bytes p0Unsigned;

        bytes32 parentHash;
        bytes p1;
        bytes32 receiptHash;
        bytes p2;

        bytes number;
        bytes p3;

        bytes p4Signed;  // rlp
        bytes p4Unsigned;  // rlp
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
        bytes32 _genesisBlockHash,
        uint256 _currentHeight
    ) internal initializer {
        require(_initialValidators.length > 0, "Length of _initialValidators must be bigger than 0");

        for (uint i = 0; i < _initialValidators.length; i++) {
            allValidators[currentValidatorSet][_initialValidators[i]] = true;
        }

        genesisBlockHash = _genesisBlockHash;
        currentHeight = _currentHeight;
    }

    function CheckPoSA_(PoSAProof memory posaProof, address sideBridgeAddress) external {
        bytes32 receiptHash = calcTransferReceiptsHash(posaProof.transfer, sideBridgeAddress);
        require(posaProof.blocks[posaProof.transferEventBlock].receiptHash == receiptHash, "Transfer event validation failed");

        for (uint i = 0; i < posaProof.blocks.length; i++) {
                CheckBlock(posaProof.blocks[i]);
            }
        }

    function CheckBlock(BlockPoSA memory block_) private {
        require(verifySignature(getUnsignedHeaderHash(block_), getSignature(block_.extraData)), "invalid signature");

        uint blockNumber = bytesToUint(block_.number);


        if (blockNumber % EPOCH_LENGTH == 0) {
            currentValidatorSet++;

            address[] memory newValidators = getValidatorSet(block_.extraData);
            for (uint i = 0; i < newValidators.length; i++) {
                allValidators[currentValidatorSet][newValidators[i]] = true;
            }
        }
    }

    // todo remove?
    function getSignedHeaderHash(BlockPoSA memory block_) private pure returns(bytes32) {
        return keccak256(abi.encodePacked(
                block_.p0Signed,

                block_.parentHash,
                block_.p1,
                block_.receiptHash,

                block_.p2,
                block_.number,
                block_.p3,

                block_.p4Signed,
                block_.extraData,

                block_.p5
            ));
    }

    function getUnsignedHeaderHash(BlockPoSA memory block_) private pure returns(bytes32) {
        return keccak256(abi.encodePacked(
                block_.p0Unsigned,

                block_.parentHash,
                block_.p1,
                block_.receiptHash,

                block_.p2,
                block_.number,
                block_.p3,

                block_.p4Unsigned,
                getExtraDataUnsigned(block_.extraData),

                block_.p5
            ));
    }

    function getSignature(bytes memory extraData) private pure returns (bytes memory) {
        return sliceBytes(extraData, extraData.length - EXTRA_SEAL_LENGTH, EXTRA_SEAL_LENGTH);
    }

    function getExtraDataUnsigned(bytes memory extraData) private pure returns(bytes memory) {
        return sliceBytes(extraData, 0, EXTRA_VANITY_LENGTH);
    }

    function getValidatorSet(bytes memory extraData) private pure returns (address[] memory) {
        uint256 currentPosition = EXTRA_VANITY_LENGTH;
        uint256 endPosition = extraData.length - EXTRA_SEAL_LENGTH;
        uint256 numValidators = (endPosition - currentPosition) / ADDRESS_LENGTH;

        address[] memory validators = new address[](numValidators);

        for (uint256 i = 0; i < numValidators; i++) {
            validators[i] = bytesToAddress(sliceBytes(extraData, currentPosition, ADDRESS_LENGTH), 0);

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

    function bytesToAddress(bytes memory _bytes, uint256 _start) private pure returns (address) {
        require(_bytes.length >= _start + 20, "toAddress_outOfBounds");
        address tempAddress;

        assembly {
            tempAddress := div(mload(add(add(_bytes, 0x20), _start)), 0x1000000000000000000000000)
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

    function sliceBytes(
        bytes memory _bytes,
        uint256 _start,
        uint256 _length
    ) private pure returns (bytes memory) {
        require(_length + 31 >= _length, "slice_overflow");
        require(_bytes.length >= _start + _length, "slice_outOfBounds");

        bytes memory tempBytes;

        assembly {
            switch iszero(_length)
            case 0 {
                tempBytes := mload(0x40)
                let lengthmod := and(_length, 31)
                let mc := add(add(tempBytes, lengthmod), mul(0x20, iszero(lengthmod)))
                let end := add(mc, _length)

                for {
                    let cc := add(add(add(_bytes, lengthmod), mul(0x20, iszero(lengthmod))), _start)
                } lt(mc, end) {
                    mc := add(mc, 0x20)
                    cc := add(cc, 0x20)
                } {
                    mstore(mc, mload(cc))
                }

                mstore(tempBytes, _length)
                mstore(0x40, and(add(mc, 31), not(31)))
            }
            default {
                tempBytes := mload(0x40)
                mstore(tempBytes, 0)

                mstore(0x40, add(tempBytes, 0x20))
            }
        }

        return tempBytes;
    }
}
