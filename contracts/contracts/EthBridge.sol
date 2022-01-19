// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";
import "./helpers/CommonBridge.sol";
import "./helpers/CheckPoA.sol";


contract EthBridge is CommonBridge, CheckPoA {
    bytes1 constant b1 = bytes1(0x01);

    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timestampSeconds)
    CommonBridge(_sideBridgeAddress, relayAddress, tokenThisAddresses, tokenSideAddresses, fee_, timestampSeconds) {}
}
