// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckPoSA.sol";


contract BSC_AmbBridge is CommonBridge, CheckPoSA {

    function initialize(CommonStructs.ConstructorArgs memory args) public initializer {
        __CommonBridge_init(args);
    }

    function submitTransferPoSA(PoSAProof calldata posaProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(posaProof.transfer.eventId);
        checkEventId(posaProof.transfer.eventId);
        CheckPoSA_(posaProof, sideBridgeAddress);
        lockTransfers(posaProof.transfer.transfers, posaProof.transfer.eventId);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }

}
