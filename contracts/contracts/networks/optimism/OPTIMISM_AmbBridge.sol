// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;


import "../../common/CommonBridge.sol";
import "../../checks/CheckUntrustless2.sol";


contract OPTIMISM_AmbBridge is CommonBridge, CheckUntrustless2 {

    function initialize(
        CommonStructs.ConstructorArgs calldata args
    ) public initializer {
        __CommonBridge_init(args);
    }

    function submitTransferUntrustless(uint eventId, CommonStructs.Transfer[] calldata transfers) public onlyRole(RELAY_ROLE) whenNotPaused {
        checkEventId(eventId);
        emit TransferSubmit(eventId);
        lockTransfers(transfers, eventId);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }

}
