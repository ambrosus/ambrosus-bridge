// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../checks/CheckAura.sol";

contract CheckAuraTest is CheckAura {
    constructor(
        address[] memory initialValidators_,
        address validatorSetAddress_,
        bytes32 lastProcessedBlock
    ) {
        __CheckAura_init(initialValidators_, validatorSetAddress_, lastProcessedBlock);
    }

    function checkAuraTest(AuraProof calldata auraProof, uint minSafetyBlocks, address sideBridgeAddress) public {
        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress);
    }

    function checkAuraTestVS(AuraProof calldata auraProof, uint minSafetyBlocks, address sideBridgeAddress, address[] memory initialValidators_) public {
        validatorSet = initialValidators_;
        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress);
    }

    function checkSignatureTest(address signer, bytes32 message, bytes memory signature) public pure {
        checkSignature(signer, message, signature);
    }

    function blockHashTest(BlockAura calldata block_) public pure returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function blockHashTestPaid(BlockAura calldata block_) public returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function bytesToUintTest(bytes4 b) public pure returns (uint) {
        return bytesToUint(b);
    }

}
