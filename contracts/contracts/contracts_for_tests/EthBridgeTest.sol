// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/EthBridge.sol";
import "../common/CommonStructs.sol";

contract EthBridgeTest is EthBridge {
    constructor(
        CommonStructs.ConstructorArgs memory args,
        address[] memory initialValidators,
        address validatorSetAddress_
    )
    EthBridge(args, initialValidators, validatorSetAddress_) {}
}
