// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./CommonStructs.sol";
import "./CheckReceiptsProof.sol";


contract CheckAura {
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
        int64 delta_index;
    }

    struct AuraProof {
        BlockAura[] blocks;
        CommonStructs.TransferProof transfer;
        ValidatorSetProof[] vs_changes;
    }

    function CheckAura_(AuraProof memory auraProof, uint minSafetyBlocks) public {
        uint safetyChainLength;

        for (uint i = 0; i < auraProof.blocks.length; i++) {
            BlockAura memory block = auraProof.blocks[i];
            // check signature, calc hash
            bytes32 block_hash = CheckBlock(block);

            if (block.type_ == - 3) { // end of safety chain
                require(safetyChainLength >= minSafetyBlocks, "safety chain too short");
                safetyChainLength = 0;
            } else {
                require(block_hash == auraProof.blocks[i + 1].parent_hash, "wrong parent hash");
                safetyChainLength++;
            }

            if (block.type_ >= 0) {// validator set change event
                // todo check vs event
            }
            else if (block.type_ == - 1) {// transfer event
                // todo check transfers event
            }

        }

    }

    function CheckBlock(BlockAura memory block) internal view returns (bytes32) {
        bytes memory common_rlp = abi.encodePacked(block.p1, block.parent_hash, block.p2, block.receipts_hash, block.p3);

        // hash without seal for signature check
        bytes32 bare_hash = keccak256(abi.encodePacked(block.p0_bare, common_rlp));
        address validator = GetValidator(bytesToUint(block.step));
        CheckSignature(validator, bare_hash, block.signature);
        // revert if wrong

        // hash with seal, for prev_hash check
        return keccak256(abi.encodePacked(block.p0_seal, common_rlp, block.s1, block.step, block.s2, block.signature));

    }


    function GetValidator(uint step) internal view returns (address) {
        // todo
        return address(this);
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
