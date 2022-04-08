// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/token/ERC20/extensions/ERC20Burnable.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract ERC20Token is ERC20, ERC20Burnable, AccessControl {
    bytes32 public constant BRIDGE_ROLE = keccak256("BRIDGE_ROLE");
    address private _bridgeAddress;

    constructor(string memory name_, string memory symbol_, address[] memory bridgeAddresses)
    ERC20(name_, symbol_) {
        _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _setBridgeAddressesRole(bridgeAddresses);
    }

    function setBridgeAddressesRole(address[] memory bridgeAddresses) public onlyRole(DEFAULT_ADMIN_ROLE) {
        _setBridgeAddressesRole(bridgeAddresses);
    }

    function _setBridgeAddressesRole(address[] memory bridgeAddresses) private {
        for (uint i = 0; i < bridgeAddresses.length; i++) {
            _setupRole(BRIDGE_ROLE, bridgeAddresses[i]);
        }
    }

    function _transfer(
        address sender,
        address recipient,
        uint256 amount
    ) internal virtual override {
        if (hasRole(BRIDGE_ROLE, sender)) {
            super._mint(recipient, amount);
        } else if (hasRole(BRIDGE_ROLE, recipient)) {
            super._burn(sender, amount);
        } else {
            super._transfer(sender, recipient, amount);
        }
    }
}