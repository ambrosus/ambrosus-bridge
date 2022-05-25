// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckAura.sol";


contract BSC_BscBridge is CommonBridge, CheckAura {

    function initialize(
        CommonStructs.ConstructorArgs memory args,
        address[] memory initialValidators,
        address validatorSetAddress,
        bytes32 lastProcessedBlock
    ) public initializer {
        __CommonBridge_init(args);
        __CheckAura_init(initialValidators, validatorSetAddress, lastProcessedBlock);
        emit Transfer(0, queue);
    }

    function submitTransferAura(AuraProof calldata auraProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(auraProof.transfer.eventId);
        checkEventId(auraProof.transfer.eventId);
        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress);
        lockTransfers(auraProof.transfer.transfers, auraProof.transfer.eventId);
    }
}
