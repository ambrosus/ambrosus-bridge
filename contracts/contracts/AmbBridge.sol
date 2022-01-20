// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./helpers/CommonBridge.sol";
import "./helpers/CheckPoW.sol";
import "./helpers/CommonStructs.sol";

contract AmbBridge is CommonBridge, CheckPoW {
    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timestampSeconds_, uint lockTime_)
    CommonBridge(_sideBridgeAddress, relayAddress, tokenThisAddresses, tokenSideAddresses, fee_, timestampSeconds_, lockTime_) {}

    function submitTransfer(
        uint event_id,
        BlockPoW[] memory blocks,
        CommonStructs.Transfer[] memory events,
        bytes[] memory proof) public onlyRole(RELAY_ROLE) {

        // todo lock transfers

        require(event_id == inputEventId + 1);
        inputEventId++;

        CheckPoW_(blocks, events, proof);

        for (uint i = 0; i < events.length; i++) {
            lockedTransfers.push(events[i]);
        }
    }
}
