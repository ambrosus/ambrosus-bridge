// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/access/AccessControl.sol";

contract Faucet is AccessControl {
    event Faucet(address indexed to, uint256 indexed eventId, uint256 amount);

    constructor(address[] memory admins) payable {
        for (uint i = 0; i < admins.length; i++)
            _setupRole(DEFAULT_ADMIN_ROLE, admins[i]);
    }

    function faucet(address toAddress, uint256 eventId, uint256 amount) public onlyRole(DEFAULT_ADMIN_ROLE) {
        withdraw(toAddress, amount);
        emit Faucet(toAddress, eventId, amount);
    }

    function withdraw(address toAddress, uint256 amount) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(address(this).balance >= amount, "not enough funds");
        payable(toAddress).transfer(amount);
    }

    receive() external payable {}


}
