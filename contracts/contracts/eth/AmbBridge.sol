// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../common/CommonStructs.sol";
import "../checks/CheckPoW.sol";
import "../tokens/IwAMB.sol";


contract AmbBridge is CommonBridge, CheckPoW {
    function initialize(CommonStructs.ConstructorArgs memory args) public initializer {
        __CommonBridge_init(args);
    }

    function submitTransferPoW(PoWProof memory powProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(powProof.transfer.event_id);

        checkEventId(powProof.transfer.event_id);

        CheckPoW_(powProof, sideBridgeAddress);

        lockTransfers(powProof.transfer.transfers, powProof.transfer.event_id);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }

}
