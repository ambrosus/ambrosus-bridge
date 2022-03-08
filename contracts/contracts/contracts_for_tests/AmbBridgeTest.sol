// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/AmbBridge.sol";
import "../common/CommonStructs.sol";

contract AmbBridgeTest is AmbBridge {
    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timeframeSeconds_, uint lockTime_, uint minSafetyBlocks_
    )
    AmbBridge(_sideBridgeAddress, relayAddress,
        tokenThisAddresses, tokenSideAddresses,
        fee_, timeframeSeconds_, lockTime_, minSafetyBlocks_
    ) {}

    function getLockedTransferTest(uint event_id, uint index) public view returns (address, address, uint) {
        CommonStructs.Transfer storage t = lockedTransfers[event_id].transfers[index];
        return (t.tokenAddress, t.toAddress, t.amount);
    }

    function lockTransfersTest(CommonStructs.Transfer[] memory events, uint event_id) public {
        lockTransfers(events, event_id);
    }

    function unlockTransfersTest(uint event_id) public {
        unlockTransfers(event_id);
    }
}
