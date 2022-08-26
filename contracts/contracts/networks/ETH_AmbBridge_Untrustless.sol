// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckUntrustless.sol";


contract ETH_AmbBridge_Untrustless is CommonBridge, CheckUntrustless {

    function initialize(
        CommonStructs.ConstructorArgs calldata args,
        uint _confirmations,
        address[] calldata _relays
    ) public initializer {
        __CommonBridge_init(args);
        _setRelaysAndConfirmations(new address[](0), _relays, _confirmations);
    }

    function submitTransferUntrustless(uint eventId, CommonStructs.Transfer[] calldata transfers) public whenNotPaused {
        // relay "role" checked at CheckUntrustless contract
        require(eventId == inputEventId + 1, "EventId out of order");

        bool confirm = checkUntrustless_(eventId, transfers);
        if (confirm) {// required count of confirmations reached
            ++inputEventId;
            emit TransferSubmit(eventId);
            lockTransfers(transfers, eventId);
            // todo need lock?
        }
    }

    function setRelaysAndConfirmations(address[] calldata toRemove, address[] calldata toAdd, uint _confirmations) public {
        require(msg.sender == address(this), "This method require multisig");
        _setRelaysAndConfirmations(toRemove, toAdd, _confirmations);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }


}
