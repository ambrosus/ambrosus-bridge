// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../common/checks/CheckPoW.sol";
import "../common/CommonStructs.sol";


contract AmbBridge is CommonBridge, CheckPoW {
    address ambWrapperAddress;

    constructor(
        CommonStructs.ConstructorArgs memory args,
        address ambWrapper_
    )
    CommonBridge(args)
    {

        // relay uses this event to know from what moment to synchronize the validator set;
        // side bridge contract must be deployed with validator set actual at the time this event was emitted.
        emit Transfer(0, queue);


        emitTestEvent(address(this), msg.sender, 10, true);

        ambWrapperAddress = ambWrapper_;
    }

    function submitTransfer(PoWProof memory powProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(powProof.transfer.event_id);

        checkEventId(powProof.transfer.event_id);

        CheckPoW_(powProof, sideBridgeAddress);

        lockTransfers(powProof.transfer.transfers, powProof.transfer.event_id);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(ADMIN_ROLE) {
        sideBridgeAddress = _sideBridgeAddress;
    }

    function setAmbWrapper(address wrapper) public onlyRole(ADMIN_ROLE) {
        ambWrapperAddress = wrapper;
    }
}
