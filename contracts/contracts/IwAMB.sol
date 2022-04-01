// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface IwAMB {
    event Wrap(address from, uint amount);
    event Unwrap(address from, uint amount);

    function wrap() external payable;

    function unwrap(uint amount) external payable;

}