// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./helpers/CommonBridge.sol";
import "./helpers/CheckPoW.sol";

contract AmbBridge is CommonBridge, CheckPoW {
    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timestampSeconds)
    CommonBridge(_sideBridgeAddress, relayAddress, tokenThisAddresses, tokenSideAddresses, fee_, timestampSeconds) {}
}
