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

    function checkPoSATest(PoSAProof calldata posaProof, address sideBridgeAddress,
        address[] memory _initialValidators, uint _initialEpoch, bytes1 _chainId) public {
        // TODO: we can't use the `__CheckPoSA_init` twice, but copy-paste is also not good
        chainId = _chainId;
        currentEpoch = _initialEpoch;
        currentValidatorSetSize = _initialValidators.length;

        for (uint i = 0; i < _initialValidators.length; i++) {
            allValidators[currentEpoch][_initialValidators[i]] = true;
        }

        checkPoSA_(posaProof, sideBridgeAddress);
    }

    function blockHashTest(BlockPoSA calldata block_) public view returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function blockHashTestPaid(BlockPoSA calldata block_) public returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function checkSignatureTest(bytes32 hash, bytes memory signature) public view returns(address) {
        return ecdsaRecover(hash, signature);
    }
}
