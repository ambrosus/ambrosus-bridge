pragma solidity ^0.8.0;

contract ProxyMultisigTest {
    bytes4 public value;

    function changeValue(bytes4 value_) public {
        value = value_;
    }

    receive() external payable {}
}
