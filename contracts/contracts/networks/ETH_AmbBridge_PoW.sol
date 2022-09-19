// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckPoW.sol";


contract ETH_AmbBridge_PoW is CommonBridge, CheckPoW {

    function initialize(
        CommonStructs.ConstructorArgs calldata args,
        uint minimumDifficulty
    ) public initializer {
        __CommonBridge_init(args);
        __CheckPoW_init(minimumDifficulty);
    }

    function submitTransferPoW(PoWProof calldata powProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(powProof.transfer.eventId);
        checkEventId(powProof.transfer.eventId);
        checkPoW_(powProof, sideBridgeAddress);
        // checkPoW_(powProof, sideBridgeAddress);
        lockTransfers(powProof.transfer.transfers, powProof.transfer.eventId);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }
}

