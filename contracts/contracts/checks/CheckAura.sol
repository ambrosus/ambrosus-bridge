// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./CheckReceiptsProof.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";

contract CheckAura is Initializable, CheckReceiptsProof {
    bytes1 constant PARENT_HASH_PREFIX = 0xA0;
    bytes1 constant STEP_PREFIX = 0x84;
    bytes2 constant SIGNATURE_PREFIX = 0xB841;

    address[] public validatorSet;
    address validatorSetAddress;


    struct BlockAura {
        bytes3 p0Seal;
        bytes3 p0Bare;

        bytes32 parentHash;
        bytes p2;
        bytes32 receiptHash;
        bytes p3;

        bytes4 step;
        bytes signature;  // todo maybe pass s r v values?

        uint64 finalizedVs;
    }


    struct ValidatorSetProof {
        bytes[] receiptProof;
        address deltaAddress;
        int64 deltaIndex; // < 0 ? remove : add
    }

    struct AuraProof {
        BlockAura[] blocks;
        CommonStructs.TransferProof transfer;
        ValidatorSetProof[] vsChanges;
        uint64 transferEventBlock;
    }



    function __CheckAura_init(
        address[] memory initialValidators_,
        address validatorSetAddress_
    ) internal initializer {
        require(initialValidators_.length > 0, "Length of _initialValidators must be bigger than 0");

        validatorSet = initialValidators_;
        validatorSetAddress = validatorSetAddress_;

    }

    function checkAura_(AuraProof memory auraProof, uint minSafetyBlocks, address sideBridgeAddress) internal {

        uint safetyChainLength;
        bytes32 blockHash;
        uint lastFinalizedVs;

        bytes32 receiptHash = calcTransferReceiptsHash(auraProof.transfer, sideBridgeAddress);
        require(auraProof.blocks[auraProof.transferEventBlock].receiptHash == receiptHash, "Transfer event validation failed");


        for (uint i = 0; i < auraProof.blocks.length; i++) {
            BlockAura memory block_ = auraProof.blocks[i];

            if (block_.finalizedVs != 0) {// 0 means no events should be finalized; so indexes are shifted by 1
                for (uint j = lastFinalizedVs; j < block_.finalizedVs; j++) {
                    ValidatorSetProof memory vsChange = auraProof.vsChanges[j];

                    handleVS(vsChange);
                    if (vsChange.receiptProof.length != 0) {
                        receiptHash = calcValidatorSetReceiptHash(vsChange.receiptProof, validatorSetAddress, validatorSet);

                        // event finalize always happened on block one after the block with event
                        // so event_block is finalized_block - 2
                        require(auraProof.blocks[i - 2].receiptHash == receiptHash, "Wrong VS receipt hash");
                        safetyChainLength = 2;
                    }
                }

                lastFinalizedVs = block_.finalizedVs - 1;
            }

            blockHash = checkBlock(block_);


            if (i + 1 != auraProof.blocks.length && blockHash == auraProof.blocks[i + 1].parentHash) {
                safetyChainLength++;
            } else if (i == auraProof.transferEventBlock) {
                safetyChainLength == 0;
            } else {
                require(safetyChainLength >= minSafetyBlocks, "wrong parent hash");
            }

        }

    }

    function getValidatorSet() public view returns (address[] memory) {
        return validatorSet;
    }

    function handleVS(ValidatorSetProof memory vsEvent) internal {
        if (vsEvent.deltaIndex < 0) {
            uint index = uint(int(vsEvent.deltaIndex * (- 1) - 1));
            validatorSet[index] = validatorSet[validatorSet.length - 1];
            validatorSet.pop();
        }
        else {
            uint index = uint(int((vsEvent.deltaIndex)));

            // logic if validatorSet contract will be updated
            // validatorSet.push(validatorSet[index]);
            // validatorSet[index] = vsEvent.deltaAddress;

            // old (current) validatorSet contract logic
            validatorSet.push(vsEvent.deltaAddress);
        }
    }

    function checkBlock(BlockAura memory block_) internal view returns (bytes32) {
        (bytes32 bareHash, bytes32 sealHash) = calcBlockHash(block_);

        address validator = validatorSet[bytesToUint(block_.step) % validatorSet.length];
        checkSignature(validator, bareHash, block_.signature);

        return sealHash;
    }

    function calcBlockHash(BlockAura memory block_) internal pure returns (bytes32, bytes32) {
        bytes memory commonRlp = abi.encodePacked(PARENT_HASH_PREFIX, block_.parentHash, block_.p2, block_.receiptHash, block_.p3);
        return (
        // hash without seal (bare), for signature check
        keccak256(abi.encodePacked(block_.p0Bare, commonRlp)),
        // hash with seal, for prev_hash check
        keccak256(abi.encodePacked(block_.p0Seal, commonRlp, STEP_PREFIX, block_.step, SIGNATURE_PREFIX, block_.signature))
        );
    }


    function checkSignature(address signer, bytes32 messageHash, bytes memory signature) internal pure {
        bytes32 r;
        bytes32 s;
        uint8 v;
        assembly {
            r := mload(add(signature, 32))
            s := mload(add(signature, 64))
            v := byte(0, mload(add(signature, 96)))
            if lt(v, 27) {v := add(v, 27)}
        }
        require(ecrecover(messageHash, v, r, s) == signer, "Failed to verify sign");
    }

    function calcValidatorSetReceiptHash(bytes[] memory receipt_proof, address validatorSetAddress, address[] memory vSet) private pure returns (bytes32) {
        bytes32 el = keccak256(abi.encodePacked(
                receipt_proof[0],
                validatorSetAddress,
                receipt_proof[1],
                abi.encode(vSet),
                receipt_proof[2]
            ));
        return calcReceiptsHash(receipt_proof, el, 3);
    }

    function bytesToUint(bytes4 b) internal pure returns (uint){
        return uint(uint32(b));
    }
}
