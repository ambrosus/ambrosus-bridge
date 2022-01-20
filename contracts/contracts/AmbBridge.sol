// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./helpers/CommonBridge.sol";
import "./helpers/CheckPoW.sol";
import "./helpers/CommonStructs.sol";

contract AmbBridge is CommonBridge, CheckPoW {
    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timeframeSeconds_, uint lockTime_)
    CommonBridge(_sideBridgeAddress, relayAddress, tokenThisAddresses, tokenSideAddresses, fee_, timeframeSeconds_, lockTime_) {}

    function submitTransfer(
        uint event_id,
        BlockPoW[] memory blocks,
        CommonStructs.Transfer[] memory events,
        bytes[] memory proof) public onlyRole(RELAY_ROLE) {

        require(event_id == inputEventId + 1);
        inputEventId++;

        CheckPoW_(blocks, events, proof);

        for (uint i = 0; i < events.length; i++) {
            lockedTransfers.transfers.push(events[i]);
        }
        lockedTransfers.endTimestamp = block.timestamp + lockTime;
    }
}
