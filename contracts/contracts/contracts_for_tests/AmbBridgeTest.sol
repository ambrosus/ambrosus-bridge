// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/AmbBridge.sol";
import "../common/CommonStructs.sol";

contract AmbBridgeTest is AmbBridge {
    constructor(
        CommonStructs.ConstructorArgs memory args,
        address ambWrapper
    )
    AmbBridge(args, ambWrapper) {}

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

    function blockHashTest(BlockPoW memory block_) public pure returns (bytes32) {
        return blockHash(block_);
    }
}
