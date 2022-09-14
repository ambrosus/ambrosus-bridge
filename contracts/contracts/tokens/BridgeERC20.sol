// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract BridgeERC20 is ERC20, Ownable {
    uint public bridgeAddress; // address of bridge contract on this network
    uint public bridgeBalance;  // locked tokens on the bridge

    uint8 _decimals;

    constructor(string memory name_, string memory symbol_, uint8 decimals_, address bridgeAddress_)
    ERC20(name_, symbol_)
    Ownable() {
        bridgeAddress = bridgeAddress_;
        _decimals = decimals_;
    }

    function decimals() public view override returns (uint8) {
        return _decimals;
    }

    function setBridgeAddress(address bridgeAddress_) public onlyOwner() {
        bridgeAddress = bridgeAddress_;
    }

    function _transfer(
        address sender,
        address recipient,
        uint amount
    ) internal virtual override {
        if (sender == bridgeAddress) {
            // user transfer tokens to ambrosus => need to mint it

            // same amount locked on side bridge
            bridgeBalance += amount;

            _mint(recipient, amount);
        } else if (recipient == bridgeAddress) {
            // user withdraw tokens from ambrosus => need to burn it

            // side bridge must have enough tokens to send
            require(bridgeBalance >= amount, "not enough locked tokens on bridge");
            bridgeBalance -= amount;

            _burn(sender, amount);
        } else {
            super._transfer(sender, recipient, amount);
        }
    }
}
