// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/AmbBridge.sol";
import "../common/CommonStructs.sol";

contract AmbBridgeTest is AmbBridge {
    constructor(CommonStructs.ConstructorArgs memory args) {
        AmbBridge.initialize(args);
    }

    function getLockedTransferTest(uint event_id) public view returns (CommonStructs.LockedTransfers memory) {
        return lockedTransfers[event_id];
    }

    function lockTransfersTest(CommonStructs.Transfer[] memory events, uint event_id) public {
        lockTransfers(events, event_id);
    }

    function blockHashTest(BlockPoW memory block_) public pure returns (bytes32) {
        return blockHash(block_);
    }
}
