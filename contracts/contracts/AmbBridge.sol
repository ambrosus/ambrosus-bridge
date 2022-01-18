// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./helpers/CommonBridge.sol";
//import "./helpers/CheckPoW.sol";

contract AmbBridge is CommonBridge {
    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_)
    CommonBridge(_sideBridgeAddress, relayAddress, tokenThisAddresses, tokenSideAddresses, fee_)
    {
        uint arrayLength = tokenThisAddresses.length;
        for (uint i = 0; i < arrayLength; i++) {
            tokenAddresses[tokenThisAddresses[i]] = tokenSideAddresses[i];
        }
    }
}
