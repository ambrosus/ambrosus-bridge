// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract BridgeERC20 is ERC20, AccessControl {
    bytes32 public constant BRIDGE_ROLE = keccak256("BRIDGE_ROLE");
    uint8 _decimals;

    constructor(string memory name_, string memory symbol_, uint8 decimals_, address[] memory bridgeAddresses)
    ERC20(name_, symbol_) {
        _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _setBridgeAddressesRole(bridgeAddresses);
        _decimals = decimals_;
    }

    function decimals() public view override returns (uint8) {
        return _decimals;
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
