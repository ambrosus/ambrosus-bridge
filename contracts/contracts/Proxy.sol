pragma solidity ^0.8.6;

import "hardhat-deploy/solc_0.8/openzeppelin/proxy/transparent/TransparentUpgradeableProxy.sol";
import "hardhat-deploy/solc_0.8/openzeppelin/proxy/transparent/ProxyAdmin.sol";

contract proxyTransparent is TransparentUpgradeableProxy {
    constructor(
        address _logic,
        address admin_,
        bytes memory _data
    ) TransparentUpgradeableProxy(_logic, admin_, _data) {}
}

contract proxyAdmin is ProxyAdmin {
    constructor(address initialOwner) ProxyAdmin(initialOwner) {}
}
