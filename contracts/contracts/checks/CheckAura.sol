// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./CheckReceiptsProof.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "./SignatureCheck.sol";


contract CheckAura is Initializable {
    bytes1 constant PARENT_HASH_PREFIX = 0xA0;
    bytes1 constant STEP_PREFIX = 0x84;
    bytes2 constant SIGNATURE_PREFIX = 0xB841;

    address[] public validatorSet;
    address validatorSetAddress;
    bytes32 public lastProcessedBlock;

    // todo editable
    uint minSafetyBlocksValidators;


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


    struct ValidatorSetChange {
        address deltaAddress;
        int64 deltaIndex; // < 0 ? remove : add
    }

    struct ValidatorSetProof {
        bytes[] receiptProof;
        ValidatorSetChange[] changes;
        uint eventBlock;  // todo add in relay
    }

    struct AuraProof {
        BlockAura[] blocks;
        CommonStructs.TransferProof transfer;
        ValidatorSetProof[] vsChanges;
        uint64 transferEventBlock;
    }


    function __CheckAura_init(
        address[] memory initialValidators_,
        address validatorSetAddress_,
        bytes32 lastProcessedBlock_
    ) internal initializer {
        require(initialValidators_.length > 0, "Length of _initialValidators must be bigger than 0");

        validatorSet = initialValidators_;
        validatorSetAddress = validatorSetAddress_;
        lastProcessedBlock = lastProcessedBlock_;

    }

    function checkAura_(AuraProof calldata auraProof, uint minSafetyBlocks, address sideBridgeAddress) internal {

        bytes32 parentHash;
        bytes32 receiptHash;

        // auraProof can be without transfer event when we have to many vsChanges and transfer doesn't fit into proof
        if (auraProof.transferEventBlock != 0) {
            bytes32 receiptHash = calcTransferReceiptsHash(auraProof.transfer, sideBridgeAddress);
            require(auraProof.blocks[auraProof.transferEventBlock].receiptHash == receiptHash, "Transfer event validation failed");
            require(auraProof.blocks.length - auraProof.transferEventBlock >= minSafetyBlocks, "Not enough safety blocks");
        }

        for (uint i = 0; i < auraProof.blocks.length; i++) {
            BlockAura calldata block_ = auraProof.blocks[i];

            if (block_.finalizedVs != 0) {// 0 means no events should be finalized, so indexes are shifted by 1
                // vs changes in that block
                ValidatorSetProof memory vsProof = auraProof.vsChanges[block_.finalizedVs - 1];

                // apply vs changes
                for (uint k = 0; k < vsProof.changes.length; k++)
                    applyVsChange(vsProof.changes[k]);

                // check proof
                receiptHash = calcValidatorSetReceiptHash(vsProof.receiptProof, validatorSetAddress, validatorSet);

                // eventBlockNum = finalizedBlockNum - validatorSet.length / 2 - 1
                // eventBlockIndex = finalizedBlockIndex - minSafetyBlocksValidators
                require(vsProof.eventBlock - i > minSafetyBlocksValidators, "Few safety blocks validators");
                        require(auraProof.blocks[vsProof.eventBlock].receiptHash == receiptHash, "Wrong VS receipt hash");


                // there is gap BEFORE finalizing block, so disable parentHash check for it
                // but only if it's not the safety blocks for transfer event
                if (i < auraProof.transferEventBlock)
                    parentHash = bytes32(0);
            }

            // don't check parentHash for first block and for block after finalizing vs
            if (parentHash != bytes32(0))
                require(block_.parentHash == parentHash, "Wrong parent hash");

            parentHash = checkBlock(block_);

            // there is gap AFTER finalizing block, so disable parentHash check for it
            // but only if it's not the safety blocks for transfer event
            if (block_.finalizedVs != 0 && i < auraProof.transferEventBlock)
                parentHash = bytes32(0);

        }

        lastProcessedBlock = parentHash;
    }

    function getValidatorSet() public view returns (address[] memory) {
        return validatorSet;
    }

    function applyVsChange(ValidatorSetChange memory vsEvent) internal {
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

    function checkBlock(BlockAura calldata block_) internal view returns (bytes32) {
        (bytes32 bareHash, bytes32 sealHash) = calcBlockHash(block_);

        address validator = validatorSet[bytesToUint(block_.step) % validatorSet.length];
        require(ecdsaRecover(bareHash, block_.signature) == validator, "Failed to verify sign");

        return sealHash;
    }

    function calcBlockHash(BlockAura calldata block_) internal pure returns (bytes32, bytes32) {
        bytes memory commonRlp = abi.encodePacked(PARENT_HASH_PREFIX, block_.parentHash, block_.p2, block_.receiptHash, block_.p3);
        return (
        // hash without seal (bare), for signature check
        keccak256(abi.encodePacked(block_.p0Bare, commonRlp)),
        // hash with seal, for prev_hash check
        keccak256(abi.encodePacked(block_.p0Seal, commonRlp, STEP_PREFIX, block_.step, SIGNATURE_PREFIX, block_.signature))
        );
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

    uint256[15] private ___gap;
}
