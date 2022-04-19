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
    )
    EthBridge(args, initialValidators, validatorSetAddress_, lastProcessedBlock) {}

    function checkAuraTest(AuraProof memory auraProof, uint minSafetyBlocks, address sideBridgeAddress, address validatorSetAddress) public {
        checkAura_(auraProof, minSafetyBlocks, sideBridgeAddress, validatorSetAddress);
    }

    function checkSignatureTest(address signer, bytes32 message, bytes memory signature) public pure {
        checkSignature(signer, message, signature);
    }

    function blockHashTest(BlockAura memory block_) public pure returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function blockHashTestPaid(BlockAura memory block_) public returns (bytes32, bytes32) {
        return calcBlockHash(block_);
    }

    function bytesToUintTest(bytes4 b) public pure returns (uint) {
        return bytesToUint(b);
    }

}
