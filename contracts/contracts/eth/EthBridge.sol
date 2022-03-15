// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";
import "../common/CommonBridge.sol";
import "../common/checks/CheckAura.sol";
import "../common/CommonStructs.sol";


contract EthBridge is CommonBridge, CheckAura {
    address validatorSetAddress;

    constructor(
        CommonStructs.ConstructorArgs memory args,
        address[] memory initialValidators,
        address validatorSetAddress_,
        bytes32 lastProcessedBlock_
    )
    CommonBridge(
        args.sideBridgeAddress, args.relayAddress,
        args.tokenThisAddresses, args.tokenSideAddresses,
        args.fee, args.feeRecipient,
        args.timeframeSeconds, args.lockTime, args.minSafetyBlocks
    )
    CheckAura(initialValidators)
    {
        emitTestEvent(address(this), msg.sender, 10, true);

        validatorSetAddress = validatorSetAddress_;
        lastProcessedBlock = lastProcessedBlock_;
    }

    function submitTransfer(AuraProof memory auraProof) public onlyRole(RELAY_ROLE) {

        require(auraProof.transfer.event_id == inputEventId + 1);
        inputEventId++;

        CheckAura_(auraProof, minSafetyBlocks, sideBridgeAddress, validatorSetAddress);

        //        lockTransfers(events, event_id);
    }
}
