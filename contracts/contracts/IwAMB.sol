// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

interface IwAMB {
    function wrap() external payable;

    function unwrap(uint amount) external payable;

}