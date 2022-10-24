// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../tokens/BridgeERC20.sol";

contract BridgeERC20Test is BridgeERC20 {
    constructor(string memory name_, string memory symbol_, uint8 decimals_, address bridgeAddress)
    BridgeERC20(name_, symbol_, decimals_, bridgeAddress) {}

    function mint(address to,  uint256 amount) public {
        _mint(to, amount);
    }

}
