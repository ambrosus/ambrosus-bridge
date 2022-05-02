// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckPoSA.sol";


contract AmbBridge is CommonBridge, CheckPoSA {

    function initialize(CommonStructs.ConstructorArgs memory args) public initializer {
        __CommonBridge_init(args);
    }

    function submitTransferPoW(PoSAProof memory posaProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(posaProof.transfer.eventId);
        checkEventId(posaProof.transfer.eventId);
        checkPoW_(posaProof, sideBridgeAddress);
        lockTransfers(posaProof.transfer.transfers, posaProof.transfer.eventId);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }

}
