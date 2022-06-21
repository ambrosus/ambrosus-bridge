// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../checks/CheckAura.sol";

contract CheckAuraTest is CheckAura {
    constructor(address validatorSetAddress_) {
        validatorSetAddress = validatorSetAddress_;
    }

    function checkAuraTest(AuraProof calldata auraProof, uint minSafetyBlocks, address sideBridgeAddress, address[] memory initialValidators_) public {
        validatorSet = initialValidators_;
        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress);
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
