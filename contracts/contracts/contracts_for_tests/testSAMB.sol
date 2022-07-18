// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../tokens/sAMB.sol";

contract testSAMB is sAMB {
    constructor(string memory name_, string memory symbol_) sAMB(name_, symbol_) {}

    function mint(address to,  uint256 amount) public {
        _mint(to, amount);
    }
}
