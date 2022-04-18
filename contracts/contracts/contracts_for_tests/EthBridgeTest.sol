// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/EthBridge.sol";
import "../common/CommonStructs.sol";

contract EthBridgeTest is EthBridge {
    constructor(
        CommonStructs.ConstructorArgs memory args,
        address[] memory initialValidators,
        address validatorSetAddress_,
        bytes32 lastProcessedBlock
    ) {
        EthBridge.initialize(args, initialValidators, validatorSetAddress_, lastProcessedBlock);
    }

    function CheckAuraTest(AuraProof memory auraProof, uint minSafetyBlocks, address sideBridgeAddress, address validatorSetAddress) public {
        CheckAura_(auraProof, minSafetyBlocks, sideBridgeAddress, validatorSetAddress);
    }

    function CheckSignatureTest(address signer, bytes32 message, bytes memory signature) public pure {
        CheckSignature(signer, message, signature);
    }

    function blockHashTest(BlockAura memory block_) public view returns (bytes32, bytes32) {
        return blockHash(block_);
    }

    function blockHashTestPaid(BlockAura memory block_) public returns (bytes32, bytes32) {
        return blockHash(block_);
    }

    function bytesToUintTest(bytes4 b) public view returns (uint) {
        return bytesToUint(b);
    }

}
