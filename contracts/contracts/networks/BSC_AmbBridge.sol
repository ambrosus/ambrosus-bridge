// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckPoSA.sol";


contract BSC_AmbBridge is CommonBridge, CheckPoSA {

    function initialize(
        CommonStructs.ConstructorArgs calldata args,
        address[] calldata initialValidators,
        uint initialEpoch,
        bytes1 chainId
    ) public initializer {
        __CommonBridge_init(args);
        __CheckPoSA_init(initialValidators, initialEpoch, chainId);
    }

    function submitTransferPoSA(PoSAProof calldata posaProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(posaProof.transfer.eventId);
        checkEventId(posaProof.transfer.eventId);
        checkPoSA_(posaProof, minSafetyBlocks, sideBridgeAddress);
        lockTransfers(posaProof.transfer.transfers, posaProof.transfer.eventId);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }

}
