// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../checks/CheckPoW.sol";

contract CheckPoWTest is CheckPoW {

    constructor(
        uint minimumDifficulty_
    ) {
        minimumDifficulty = minimumDifficulty_;
    }

    function checkPoWTest(PoWProof calldata powProof, address sideBridgeAddress) public {
        checkPoW_(powProof, sideBridgeAddress);
    }

    function verifyEthashTest(BlockPoW calldata block_) public view {
        verifyEthash(block_);
    }

    function blockHashTest(BlockPoW calldata block_) public pure returns (bytes32) {
        return blockHash(block_);
    }

}
