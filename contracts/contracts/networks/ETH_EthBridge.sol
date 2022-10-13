// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckUntrustless2.sol";


contract ETH_EthBridge is CommonBridge, CheckUntrustless2 {

    function initialize(
        CommonStructs.ConstructorArgs calldata args
    ) public initializer {
        __CommonBridge_init(args);
    }

    function upgrade(
        address[] calldata _watchdogs,
        address _fee_provider
    ) public {
        require(msg.sender == address(this), "This method require multisig");
        _setupRoles(WATCHDOG_ROLE, _watchdogs);
        _setupRole(FEE_PROVIDER_ROLE, _fee_provider);
    }

    function submitTransferUntrustless(uint eventId, CommonStructs.Transfer[] calldata transfers) public onlyRole(RELAY_ROLE) whenNotPaused {
        checkEventId(eventId);
        emit TransferSubmit(eventId);
        lockTransfers(transfers, eventId);
    }
}
