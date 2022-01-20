// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";
import "./helpers/CommonBridge.sol";
import "./helpers/CheckPoA.sol";


contract EthBridge is CommonBridge, CheckPoA {
    bytes1 constant b1 = bytes1(0x01);

    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timestampSeconds_, uint lockTime_)
    CommonBridge(_sideBridgeAddress, relayAddress, tokenThisAddresses, tokenSideAddresses, fee_, timestampSeconds_, lockTime_) {}

    function submitTransfer(
        uint event_id,
        BlockPoA[] memory blocks,
        CommonStructs.Transfer[] memory events,
        bytes[] memory proof) public onlyRole(RELAY_ROLE) {

        require(event_id == inputEventId + 1);
        inputEventId++;

        CheckPoA_(blocks, events, proof);

        for (uint i = 0; i < events.length; i++) {
            lockedTransfers.push(events[i]);
        }
    }
}
