// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "./BridgeERC20.sol";

contract MockERC20 is ERC20Token {
    constructor(address[] memory bridgeAddresses)
    ERC20Token("SuperToken", "ST", bridgeAddresses) {}

    function mint(address to,  uint256 amount) public {
        _mint(to, amount);
    }
}
