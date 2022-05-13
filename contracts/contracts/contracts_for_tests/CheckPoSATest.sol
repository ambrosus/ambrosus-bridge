// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../checks/CheckPoSA.sol";


contract CheckPoSATest is CheckPoSA {
    constructor(
        address[] memory initialValidators,
        uint initialEpoch,
        bytes1 chainId
    ) {
        __CheckPoSA_init(initialValidators, initialEpoch, chainId);
    }

    function blockHashTest(BlockPoSA calldata block_) public view returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function blockHashTestPaid(BlockPoSA calldata block_) public returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function checkSignatureTest(bytes32 hash, bytes memory signature) public view returns(address) {
        return getSigner(hash, signature);
    }
}
