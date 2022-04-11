// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../common/CommonStructs.sol";
import "../checks/CheckPoW.sol";
import "../tokens/IwAMB.sol";


contract AmbBridge is CommonBridge, CheckPoW {
    address ambWrapperAddress;

    constructor(
        CommonStructs.ConstructorArgs memory args,
        address ambWrapper_
    )
    CommonBridge(args)
    {
        ambWrapperAddress = ambWrapper_;
    }

    function wrap_withdraw(address toAddress) public payable {
        address tokenExternalAddress = tokenAddresses[ambWrapperAddress];
        require(tokenExternalAddress != address(0), "Unknown token address");

        require(msg.value > fee, "msg.value can't be lesser than fee");
        feeRecipient.transfer(fee);

        uint restOfValue = msg.value - fee;
        IwAMB(ambWrapperAddress).wrap{value: restOfValue}();

        //
        queue.push(CommonStructs.Transfer(tokenExternalAddress, toAddress, restOfValue));
        emit Withdraw(msg.sender, outputEventId, fee);

        withdraw_finish();
    }

    function submitTransfer(PoWProof memory powProof) public onlyRole(RELAY_ROLE) whenNotPaused {
        emit TransferSubmit(powProof.transfer.event_id);

        checkEventId(powProof.transfer.event_id);

        CheckPoW_(powProof, sideBridgeAddress);

        lockTransfers(powProof.transfer.transfers, powProof.transfer.event_id);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(DEFAULT_ADMIN_ROLE) {
        require(sideBridgeAddress == address(0), "sideBridgeAddress already set");
        sideBridgeAddress = _sideBridgeAddress;
    }

    function setAmbWrapper(address wrapper) public onlyRole(DEFAULT_ADMIN_ROLE) {
        ambWrapperAddress = wrapper;
    }
}
