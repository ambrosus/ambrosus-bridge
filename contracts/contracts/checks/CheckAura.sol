// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonStructs.sol";
import "./CheckReceiptsProof.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract CheckAura is Initializable, CheckReceiptsProof {
    bytes1 constant parentHashPrefix = 0xA0;
    bytes1 constant stepPrefix = 0x84;
    bytes2 constant signaturePrefix = 0xB841;

    bytes32 public lastProcessedBlock;
    address[] public validatorSet;


    struct BlockAura {
        bytes3 p0_seal;
        bytes3 p0_bare;

        bytes32 parent_hash;
        bytes p2;
        bytes32 receipt_hash;
        bytes p3;

        bytes4 step;
        bytes signature;  // todo maybe pass s r v values?

        uint64 finalized_vs;
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
        uint64 transfer_event_block;
    }



    function __CheckAura_init(address[] memory _initialValidators) internal initializer {
        require(_initialValidators.length > 0, "Length of _initialValidators must be bigger than 0");
        validatorSet = _initialValidators;
    }

    function CheckAura_(AuraProof memory auraProof, uint minSafetyBlocks,
        address sideBridgeAddress, address validatorSetAddress) internal {

        uint safetyChainLength;
        bytes32 block_hash;
        uint last_finalized_vs;

        bytes32 receiptHash = CalcTransferReceiptsHash(auraProof.transfer, sideBridgeAddress);
        require(auraProof.blocks[auraProof.transfer_event_block].receipt_hash == receiptHash, "Transfer event validation failed");


        for (uint i = 0; i < auraProof.blocks.length; i++) {
            BlockAura memory block_ = auraProof.blocks[i];

            if (block_.finalized_vs != 0) {  // 0 means no events should be finalized; so indexes are shifted by 1
                for (uint j = last_finalized_vs; j < block_.finalized_vs; j++) {
                    ValidatorSetProof memory vs_change = auraProof.vs_changes[j];

                    handleVS(vs_change);
                    if (vs_change.receipt_proof.length != 0) {
                        bytes32 receiptHash = CalcValidatorSetReceiptHash(vs_change.receipt_proof, validatorSetAddress, validatorSet);

                        // event finalize always happened on block one after the block with event
                        // so event_block is finalized_block - 2
                        require(auraProof.blocks[i - 2].receipt_hash == receiptHash, "Wrong VS receipt hash");
                        safetyChainLength = 2;
                    }
                }

                last_finalized_vs = block_.finalized_vs - 1;
            }

            block_hash = CheckBlock(block_);


            if (i+1 != auraProof.blocks.length && block_hash == auraProof.blocks[i + 1].parent_hash) {
                safetyChainLength++;
            } else if (i == auraProof.transfer_event_block) {
                safetyChainLength == 0;
            } else {
                require(safetyChainLength >= minSafetyBlocks, "wrong parent hash");
            }

        }

        lastProcessedBlock = block_hash;
    }

    function handleVS(ValidatorSetProof memory vsEvent) private {
        if (vsEvent.delta_index < 0) {
            uint index = uint(int(vsEvent.delta_index * (- 1) - 1));
            validatorSet[index] = validatorSet[validatorSet.length - 1];
            validatorSet.pop();
        }
        else {
            uint index = uint(int((vsEvent.delta_index)));
            validatorSet.push(validatorSet[index]);
            validatorSet[index] = vsEvent.delta_address;
        }
    }

    function CheckBlock(BlockAura memory block_) internal view returns (bytes32) {
        (bytes32 bare_hash, bytes32 seal_hash) = blockHash(block_);

        address validator = validatorSet[bytesToUint(block_.step) % validatorSet.length];
        CheckSignature(validator, bare_hash, block_.signature);

        return seal_hash;
    }

    function blockHash(BlockAura memory block_) internal pure returns (bytes32, bytes32) {
        bytes memory common_rlp = abi.encodePacked(parentHashPrefix, block_.parent_hash, block_.p2, block_.receipt_hash, block_.p3);
        return (
        // hash without seal (bare), for signature check
        keccak256(abi.encodePacked(block_.p0_bare, common_rlp)),
        // hash with seal, for prev_hash check
        keccak256(abi.encodePacked(block_.p0_seal, common_rlp, stepPrefix, block_.step, signaturePrefix, block_.signature))
        );
    }

    function GetValidatorSet() public view returns (address[] memory) {
        return validatorSet;
    }


    function CheckSignature(address signer, bytes32 message_hash, bytes memory signature) internal pure {
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
            if lt(v, 27) {v := add(v, 27)}
        }
        require(ecrecover(message_hash, v, r, s) == signer, "Failed to verify sign");
    }

    function CalcValidatorSetReceiptHash(bytes[] memory receipt_proof,
        address validatorSetAddress,
        address[] memory vSet) private pure returns (bytes32) {

        bytes32 el = keccak256(abi.encodePacked(
                receipt_proof[0],
                validatorSetAddress,
                receipt_proof[1],
                abi.encode(vSet),
                receipt_proof[2]
            ));
        return CalcReceiptsHash(receipt_proof, el, 3);
    }

    function bytesToUint(bytes4 b) internal pure returns (uint){
        return uint(uint32(b));
    }
}
