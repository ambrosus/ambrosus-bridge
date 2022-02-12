// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./CommonStructs.sol";


contract CheckPoA {
    struct BlockPoA {
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

        uint type_;
    }


    struct ValidatorSet_Event {
        bytes[] receipt_proof;
        address delta_address;
        uint96 delta_index;  // 12байт шоб спаковать с 20байт адресом
    }


    struct Transfer_Event {
        bytes[] receipt_proof;
        uint event_id;
        CommonStructs.Transfer[] transfers;
    }


    function CheckPoA_(BlockPoA[] memory blocks, Transfer_Event memory transfer, ValidatorSet_Event[] memory vs_changes) public {
        uint safetyChainLength;

        for (uint i = 0; i < blocks.length; i++) {
            // check signature, calc hash
            bytes32 block_hash = CheckBlock(blocks[i]);

            if (blocks[i].type_ == -3) {  // end of safety chain
                require(safetyChainLength >= minSafetyBlocks, "safety chain too short");
                safetyChainLength = 0;
            } else {
                require(block_hash == blocks[i+1].parent_hash, "wrong parent hash");
                safetyChainLength++;
            }

            if (block.type_ >= 0) { // validator set change event
                // todo check vs event
            }
            else if (block.type_ == -1) { // transfer event
                // todo check transfers event
            }

        }

    }

    function CheckBlock(BlockPoA memory block) internal view returns (bytes32) {
        bytes memory common_rlp = abi.encodePacked(blocks[i].p1, block.parent_hash, blocks[i].p2, blocks[i].receipts_hash, blocks[i].p3);

        // hash without seal for signature check
        bytes32 bare_hash = keccak256(abi.encodePacked(blocks[i].p0_bare, common_rlp));
        address validator = GetValidator(bytesToUint(blocks[i].step));
        CheckSignature(validator, bare_hash, blocks[i].signature);  // revert if wrong

        // hash with seal, for prev_hash check
        return keccak256(abi.encodePacked(blocks[i].p0_seal, common_rlp, blocks[i].s1, blocks[i].step, blocks[i].s2, blocks[i].signature));

    }



    function CheckReceiptsProof(bytes[] memory proof, address eventContractAddress, bytes topic, bytes data, bytes32 receciptsRoot) public {
        require(calcReceiptsRoot(proof, eventContractAddress, topic,  data) == receiptsRoot, "Failed to verify receipts proof");
    }

    // check readme for focs
    function calcReceiptsRoot(bytes[] memory proof, address eventContractAddress, bytes topic, bytes data) public view returns (bytes32){
        bytes32 el = keccak256(abi.encodePacked(proof[0], eventContractAddress, proof[1], topic, proof[2], data, proof[3]));
        bytes memory s;

        for (uint i = 4; i < proof.length; i += 2) {
            s = abi.encodePacked(proof[i], el, proof[i + 1]);
            el = (s.length > 32) ? keccak256(s) : bytes32(s);
        }

        return el;
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
