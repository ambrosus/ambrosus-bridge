// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckUntrastless.sol";


contract ETH_AmbBridge_Untrastless is CommonBridge, CheckUntrastless {

    function initialize(
        CommonStructs.ConstructorArgs calldata args,
        uint _confirmations,
        address[] calldata _relays
    ) public initializer {
        __CommonBridge_init(args);
        _setRolesAndConfirmations(new address[](0), _relays, _confirmations);
    }

    function submitTransferUntrastless(uint eventId, CommonStructs.Transfer[] calldata transfers) public onlyRole(RELAY_ROLE) whenNotPaused {
        checkEventId(eventId);

        bool confirm = checkUntrastless_(eventId, transfers);
        if (confirm) {// required count of confirmations reached
            emit TransferSubmit(powProof.transfer.eventId);
            lockTransfers(transfers, eventId);
            // todo need lock?
        }
    }

    function setRolesAndConfirmations(address[] calldata toRemove, address[] calldata toAdd, uint _confirmations) public {
        require(msg.sender == address(this), "This method require multisig");
        _setRolesAndConfirmations(toRemove, toAdd, _confirmations);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }


    function _setRolesAndConfirmations(address[] calldata toRemove, address[] calldata toAdd, uint _confirmations) internal {
        for (uint i = 0; i < toRemove.length; i++)
            _removeRole(RELAY_ROLE, toRemove[i]);
        for (uint i = 0; i < toAdd.length; i++)
            _grantRole(RELAY_ROLE, toAdd[i]);
        confirmationsThreshold = _confirmations;
    }

}
