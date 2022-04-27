// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../eth/AmbBridge.sol";

contract CommonBridgeTest is CommonBridge {
    constructor(
        CommonStructs.ConstructorArgs memory args
    ) {
        __CommonBridge_init(args);
    }
}
