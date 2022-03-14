// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../CommonStructs.sol";
import "./CheckReceiptsProof.sol";


contract CheckAura is CheckReceiptsProof {
    // bitmask
    uint8 constant BlTypeSafetyEnd = 1;
    uint8 constant BlTypeSafety = 2;
    uint8 constant BlTypeTransfer = 4;
    uint8 constant BlTypeVSChange = 8;

    bytes32 public lastProcessedBlock;

    struct BlockAura {
        bytes p0_seal;
        bytes p0_bare;

        bytes p1;
        bytes32 parent_hash;
        bytes p2;
        bytes32 receipt_hash;
        bytes p3;

        bytes s1;
        bytes step;
        bytes s2;
        bytes signature;

        uint8 type_;
        int64 delta_index;
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

    function CheckAura_(AuraProof memory auraProof, uint minSafetyBlocks,
        address sideBridgeAddress, address validatorSetAddress) public {

        // validator set change event
        uint n = uint(int(auraProof.blocks[0].delta_index));
        for (uint i = 0; i < n; i++) {
            handleVS(auraProof.vs_changes[i]);
        }

        uint safetyChainLength;
        bytes32 block_hash;

        for (uint i = 0; i < auraProof.blocks.length; i++) {
            BlockAura memory block_ = auraProof.blocks[i];
            // check signature, calc hash
            block_hash = CheckBlock(block_);

            if (block_.type_ & BlTypeSafetyEnd != 0) { // end of safety chain
                require(safetyChainLength >= minSafetyBlocks, "safety chain too short");
                safetyChainLength = 0;
            } else {
                require(block_hash == auraProof.blocks[i + 1].parent_hash, "wrong parent hash");
                safetyChainLength++;
            }

            if (block_.type_ & BlTypeVSChange != 0) {// validator set change event
                ValidatorSetProof memory vsEvent = auraProof.vs_changes[i];
                handleVS(vsEvent);

                if (vsEvent.receipt_proof.length != 0) {
                    bytes32 receiptHash = CalcValidatorSetReceiptHash(auraProof, validatorSetAddress, validatorSet);
                    require(block_.receipt_hash == receiptHash, "Wrong Hash");
                }
            }

            // transfer event
            if (block_.type_ & BlTypeTransfer != 0) {
                bytes32 receiptHash = CalcTransferReceiptsHash(auraProof.transfer, sideBridgeAddress);
                require(block_.receipt_hash == receiptHash, "Transfer event validation failed");

                lastProcessedBlock = block_hash;
            }
        }
    }

    function handleVS(ValidatorSetProof memory vsEvent) private {
        if (vsEvent.delta_index < 0) {
            uint index = uint(int(vsEvent.delta_index * (-1) - 1));
            validatorSet[index] = validatorSet[validatorSet.length - 1];
            validatorSet.pop();
        }
        else {
            uint index = uint(int((vsEvent.delta_index)));
            validatorSet.push(validatorSet[index]);
            validatorSet[index] = vsEvent.delta_address;
        }
    }

    function CheckBlock(BlockAura memory block) internal view returns (bytes32) {
        bytes memory common_rlp = abi.encodePacked(block.p1, block.parent_hash, block.p2, block.receipt_hash, block.p3);

        // hash without seal for signature check
        bytes32 bare_hash = keccak256(abi.encodePacked(block.p0_bare, common_rlp));
        address validator = GetValidator(bytesToUint(block.step));
        CheckSignature(validator, bare_hash, block.signature);
        // revert if wrong

        // hash with seal, for prev_hash check
        return keccak256(abi.encodePacked(block.p0_seal, common_rlp, block.s1, block.step, block.s2, block.signature));

    }


    function GetValidator(uint step) internal view returns (address) {
        return validatorSet[step % validatorSet.length];
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

    function CalcValidatorSetReceiptHash(AuraProof memory auraProof,
                                         address validatorSetAddress,
                                         address[] memory vSet) private returns(bytes32) {

        bytes32 el = keccak256(abi.encodePacked(
            auraProof.transfer.receipt_proof[0],
            validatorSetAddress,
            auraProof.transfer.receipt_proof[1],
            abi.encode(vSet),
            auraProof.transfer.receipt_proof[2]
        ));
        return CalcReceiptsHash(auraProof.transfer.receipt_proof, el, 3);
    }

    function bytesToUint(bytes memory b) private pure returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }
}
