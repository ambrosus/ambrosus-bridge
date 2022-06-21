// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract ProxyMultisigTest {
    bytes4 public value;

    function changeValue(bytes4 value_) public {
        value = value_;
    }

    receive() external payable {}
}
