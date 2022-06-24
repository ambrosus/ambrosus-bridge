// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckAura.sol";


contract BSC_BscBridge is CommonBridge, CheckAura {

    function initialize(
        CommonStructs.ConstructorArgs calldata args,
        address[] calldata initialValidators,
        address validatorSetAddress,
        bytes32 lastProcessedBlock,
        uint minSafetyBlocksValidators
    ) public initializer {
        __CommonBridge_init(args);
        __CheckAura_init(initialValidators, validatorSetAddress, lastProcessedBlock, minSafetyBlocksValidators);
    }


    function changeMinSafetyBlocksValidators(uint minSafetyBlocksValidators_) public onlyRole(ADMIN_ROLE) {
        minSafetyBlocksValidators = minSafetyBlocksValidators_;
    }

    function submitTransferAura(AuraProof calldata auraProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(auraProof.transfer.eventId);
        checkEventId(auraProof.transfer.eventId);
        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress);
        lockTransfers(auraProof.transfer.transfers, auraProof.transfer.eventId);
    }

    function submitValidatorSetChangesAura(AuraProof calldata auraProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        require(auraProof.transfer.eventId == 0, "Event id must be 0");
        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress);
    }
}
