pragma solidity ^0.8.6;

library CommonStructs {
    struct Transfer {
        address tokenAddress;
        address toAddress;
        uint amount;
    }
}
