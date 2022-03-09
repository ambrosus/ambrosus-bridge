// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../CommonStructs.sol";
import "./CheckReceiptsProof.sol";


contract CheckAura is CheckReceiptsProof {
    struct BlockAura {
        bytes p0_seal;
        bytes p0_bare;

        bytes p1;
        bytes32 parent_hash;
        bytes p2;
        bytes32 receipts_hash;
        bytes p3;

        bytes s1;
        bytes step;
        bytes s2;
        bytes signature;

        int type_;
    }


    struct ValidatorSetProof {
        bytes[] receipt_proof;
        address delta_address;
        int64 delta_index; // < 0 ? remove : add
    }

    struct AuraProof {
        BlockAura[] blocks;
        CommonStructs.TransferProof transfer;
        ValidatorSetProof[] vs_changes;
    }

    address[] public validatorSet;


    constructor(address[] memory _initialValidators) {
        require(_initialValidators.length > 0, "Length of _initialValidators must be bigger than 0");
        validatorSet = _initialValidators;
    }

    function addValidator(uint index, address validator) internal {
        validatorSet.push(validatorSet[index]);
        validatorSet[index] = validator;
    }

    function removeValidator(uint index) internal {
        validatorSet[index] = validatorSet[validatorSet.length - 1];
        validatorSet.pop();
    }

    function CheckAura_(AuraProof memory auraProof, uint minSafetyBlocks) public {
        address sideBridgeAddress = address (this);  // todo



    uint safetyChainLength;

        for (uint i = 0; i < auraProof.blocks.length; i++) {
            BlockAura memory block_ = auraProof.blocks[i];
            // check signature, calc hash
            bytes32 block_hash = CheckBlock(block_);

            if (block_.type_ == - 3) { // end of safety chain
                require(safetyChainLength >= minSafetyBlocks, "safety chain too short");
                safetyChainLength = 0;
            } else {
                require(block_hash == auraProof.blocks[i + 1].parent_hash, "wrong parent hash");
                safetyChainLength++;
            }

            if (block_.type_ >= 0) {// validator set change event
                // todo check vs event
                // val set
            }
            else if (block_.type_ == - 1) {// transfer event
                bytes32 receiptHash = CalcTransferReceiptsHash(auraProof.transfer, sideBridgeAddress);
                require(block_.receipts_hash == receiptHash, "Transfer event validation failed");
            }

        }

    }

    function CheckBlock(BlockAura memory block_) internal view returns (bytes32) {
        bytes memory common_rlp = abi.encodePacked(block_.p1, block_.parent_hash, block_.p2, block_.receipts_hash, block_.p3);

        // hash without seal for signature check
        bytes32 bare_hash = keccak256(abi.encodePacked(block_.p0_bare, common_rlp));
        address validator = GetValidator(bytesToUint(block_.step));
        CheckSignature(validator, bare_hash, block_.signature);
        // revert if wrong

        // hash with seal, for prev_hash check
        return keccak256(abi.encodePacked(block_.p0_seal, common_rlp, block_.s1, block_.step, block_.s2, block_.signature));

    }


    function GetValidator(uint step) internal view returns (address) {
        // todo
        return address(this);
    }

    function GetValidatorSet() public view returns (address[] memory) {
        return validatorSet;
    }

    function CheckSignature(address signer, bytes32 message, bytes memory signature) internal view {
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
        }
        require(
            ecrecover(keccak256(abi.encodePacked("\x19Ethereum Signed Message:\n32", message)), v, r, s) == signer,
            "Failed to verify sign");
    }


    function bytesToUint(bytes memory b) public view returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }
}
