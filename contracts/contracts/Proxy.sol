pragma solidity ^0.8.6;

import "@openzeppelin/contracts/proxy/Proxy.sol";
import "@openzeppelin/contracts/utils/StorageSlot.sol";
import "./MultiSigWallet.sol";

// todo rename contract
// todo + create utils dir for multisig and proxy
contract proxyTransparent is Proxy, MultiSigWallet {
    bytes32 internal constant _IMPLEMENTATION_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;

    constructor(
        address _logic,
        bytes memory _data,
        address[] memory owners,
        uint _required

    ) MultiSigWallet(owners, _required) {}

    function _implementation() internal view virtual override returns (address impl) {
        return StorageSlot.getAddressSlot(_IMPLEMENTATION_SLOT).value;
    }

    receive() external payable override {
        if (msg.value > 0)
            emit Deposit(msg.sender, msg.value);
        super._fallback();
    }
}
