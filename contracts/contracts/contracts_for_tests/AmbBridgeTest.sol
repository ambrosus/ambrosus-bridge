// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/AmbBridge.sol";

contract AmbBridgeTest is AmbBridge {
    function initialize_(CommonStructs.ConstructorArgs memory args, uint minimumDifficulty) public initializer {
        AmbBridge.initialize(args, minimumDifficulty);
    }

    function blockHashTest(BlockPoW calldata block_) public pure returns (bytes32) {
        return blockHash(block_);
    }

    function verifyEthashTest(BlockPoW calldata block_) public view {
        verifyEthash(block_);
    }

    function checkPoWTest(PoWProof calldata powProof, address sideBridgeAddress) public {
        checkPoW_(powProof, sideBridgeAddress);
    }

}
