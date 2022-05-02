// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/EthBridge.sol";

contract EthBridgeTest is EthBridge {
    constructor(
        CommonStructs.ConstructorArgs memory args,
        address[] memory initialValidators,
        address validatorSetAddress_
    ) {
        EthBridge.initialize(args, initialValidators, validatorSetAddress_);
    }

    function checkAuraTest(AuraProof calldata auraProof, uint minSafetyBlocks, address sideBridgeAddress) public {
        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress);
    }

    function checkSignatureTest(address signer, bytes32 message, bytes memory signature) public pure {
        checkSignature(signer, message, signature);
    }

    function blockHashTest(BlockAura calldata block_) public pure returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function blockHashTestPaid(BlockAura calldata block_) public returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function bytesToUintTest(bytes4 b) public pure returns (uint) {
        return bytesToUint(b);
    }

}
