pragma solidity ^0.8.6;

import "@openzeppelin/contracts/proxy/Proxy.sol";
import "@openzeppelin/contracts/utils/StorageSlot.sol";
import "@openzeppelin/contracts/utils/Address.sol";
import "./MultiSigWallet.sol";


contract proxyMultiSig is Proxy, MultiSigWallet {
    bytes32 internal constant _IMPLEMENTATION_SLOT = 0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc;
    bytes32 private constant _ROLLBACK_SLOT = 0x4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd9143;

    bytes32 private constant ADMIN_STORAGE_LOCATION = 0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103;

    bytes4 private constant UPGRADE_TO_AND_CALL__SIGNATURE = bytes4(keccak256("upgradeToAndCall_(address,bytes)"));

    event Upgraded(address indexed implementation);


    constructor(
        address _logic,
        bytes memory _data,
        address[] memory owners,
        uint _required

    ) MultiSigWallet(owners, _required) {
        _upgradeToAndCall(_logic, _data, false);

        StorageSlot.getAddressSlot(ADMIN_STORAGE_LOCATION).value = msg.sender;  // trick the hardhat-deploy
    }


    function implementation() external returns (address implementation_) {
        implementation_ = _implementation();
    }

    function upgradeTo(address newImplementation) external payable ownerExists(msg.sender) {
        submitTransaction(
            address(this),
            msg.value,
            abi.encodeWithSelector(UPGRADE_TO_AND_CALL__SIGNATURE, newImplementation, bytes(""))
        );
    }

    function upgradeToAndCall(address newImplementation, bytes calldata data) external payable ownerExists(msg.sender) {
        submitTransaction(
            address(this),
            msg.value,
            abi.encodeWithSelector(UPGRADE_TO_AND_CALL__SIGNATURE, newImplementation, data)
        );
    }

    function upgradeToAndCall_(address newImplementation, bytes calldata data) external onlyWallet payable {
        _upgradeToAndCall(newImplementation, data, false);
    }

    function _upgradeToAndCall(
        address newImplementation,
        bytes memory data,
        bool forceCall
    ) internal {
        _upgradeTo(newImplementation);
        if (data.length > 0 || forceCall) {
            Address.functionDelegateCall(newImplementation, data);
        }
    }

    function _upgradeTo(address newImplementation) internal {
        _setImplementation(newImplementation);
        emit Upgraded(newImplementation);
    }

    function _setImplementation(address newImplementation) private {
        require(Address.isContract(newImplementation), "ERC1967: new implementation is not a contract");
        StorageSlot.getAddressSlot(_IMPLEMENTATION_SLOT).value = newImplementation;
    }

    function _implementation() internal view override returns (address) {
        return StorageSlot.getAddressSlot(_IMPLEMENTATION_SLOT).value;
    }

    receive() external payable override {
        if (msg.value > 0)
            emit Deposit(msg.sender, msg.value);
        super._fallback();
    }
}
