// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./IwAMB.sol";

contract wAMB is IwAMB, ERC20 {
    constructor(string memory name_, string memory symbol_) ERC20(name_, symbol_) {}

    function wrap() public virtual override payable {
        _mint(msg.sender, msg.value);
    }

    function unwrap(uint amount) public virtual override payable {
        _burn(msg.sender, amount);
        payable(msg.sender).transfer(amount);
    }

}