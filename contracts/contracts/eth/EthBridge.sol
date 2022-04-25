// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../common/CommonStructs.sol";
import "../checks/CheckAura.sol";


contract EthBridge is CommonBridge, CheckAura {
    address validatorSetAddress;

    constructor(
        CommonStructs.ConstructorArgs memory args,
        address[] memory initialValidators,
        address validatorSetAddress_,
        bytes32 lastProcessedBlock_
    )
    CommonBridge(args)
    CheckAura(initialValidators)
    {
        validatorSetAddress = validatorSetAddress_;
        lastProcessedBlock = lastProcessedBlock_;
    }

    function submitTransferAura(AuraProof memory auraProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(auraProof.transfer.eventId);

        checkEventId(auraProof.transfer.eventId);

//        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress, validatorSetAddress);

        lockTransfers(auraProof.transfer.transfers, auraProof.transfer.eventId);
    }
}
