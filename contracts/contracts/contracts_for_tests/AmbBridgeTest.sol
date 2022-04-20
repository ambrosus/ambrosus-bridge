// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/AmbBridge.sol";
import "../common/CommonStructs.sol";

contract AmbBridgeTest is AmbBridge {
    constructor(CommonStructs.ConstructorArgs memory args) AmbBridge(args) {}

    function getLockedTransferTest(uint eventId) public view returns (CommonStructs.LockedTransfers memory) {
        return lockedTransfers[eventId];
    }

    function lockTransfersTest(CommonStructs.Transfer[] memory events, uint eventId) public {
        lockTransfers(events, eventId);
    }

    function blockHashTest(BlockPoW memory block_) public pure returns (bytes32) {
        return blockHash(block_);
    }

    function verifyEthashTest(BlockPoW memory block_) public view {
        verifyEthash(block_);
    }

    function checkPoWTest(PoWProof memory powProof, address sideBridgeAddress) public {
        checkPoW_(powProof, sideBridgeAddress);
    }

    function calcTransferReceiptsHashTest(CommonStructs.TransferProof memory p, address eventContractAddress) public pure returns (bytes32) {
        return calcTransferReceiptsHash(p, eventContractAddress);
    }

}
