// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";
import "./common/CommonBridge.sol";
import "./common/CheckPoA.sol";


contract EthBridge is CommonBridge, CheckPoA {
    bytes1 constant b1 = bytes1(0x01);

    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timeframeSeconds_, uint lockTime_, uint minSafetyBlocks_
    )
    CommonBridge(
        _sideBridgeAddress, relayAddress,
        tokenThisAddresses, tokenSideAddresses,
        fee_, timeframeSeconds_, lockTime_, minSafetyBlocks_
    )
    {}

    function submitTransfer(
        BlockPoA[] memory blocks, Transfer_Event memory transfer, ValidatorSet_Event[] memory vs_changes
    ) public onlyRole(RELAY_ROLE) {

        require(transfer.event_id == inputEventId + 1);
        inputEventId++;

        CheckPoA_(blocks, transfer, vs_changes, minSafetyBlocks);

//        lockTransfers(events, event_id);
    }
}
