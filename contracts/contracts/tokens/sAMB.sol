// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "./IWrapper.sol";

contract sAMB is IWrapper, ERC20 {
    constructor(string memory name_, string memory symbol_) ERC20(name_, symbol_) {}

    function deposit() public override payable {
        _mint(msg.sender, msg.value);

        emit Deposit(msg.sender, msg.value);
    }

    function withdraw(uint amount) public override payable {
        _burn(msg.sender, amount);
        payable(msg.sender).transfer(amount);

        emit Withdrawal(msg.sender, msg.value);
    }

}
