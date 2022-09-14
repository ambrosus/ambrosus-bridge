// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../tokens/BridgeERC20_Amb.sol";

contract BridgeERC20_AmbTest is BridgeERC20_Amb {
    constructor(string memory name_, string memory symbol_, uint8 decimals_,
        address[] memory bridgeAddresses_, uint8[] memory sideTokenDecimals_
    )
    BridgeERC20_Amb(name_, symbol_, decimals_, bridgeAddresses_, sideTokenDecimals_) {}

    function mint(address to, uint256 amount) public {
        _mint(to, amount);
    }

    function changeBridgeBalance(address bridge, uint balance) public {
        bridgeBalances[bridge] = balance;
    }
}
