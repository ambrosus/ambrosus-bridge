// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../common/CommonStructs.sol";
import "../checks/CheckPoW.sol";


contract AmbBridge is CommonBridge, CheckPoW {
    constructor(CommonStructs.ConstructorArgs memory args) CommonBridge(args) {}

    function submitTransferPoW(PoWProof memory powProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(powProof.transfer.eventId);

        checkEventId(powProof.transfer.eventId);

        checkPoW_(powProof, sideBridgeAddress);

        lockTransfers(powProof.transfer.transfers, powProof.transfer.eventId);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }

}
