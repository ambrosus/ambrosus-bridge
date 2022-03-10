// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../common/checks/CheckPoW.sol";

contract AmbBridge is CommonBridge, CheckPoW {
    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, address payable feeRecipient_,
        uint timeframeSeconds_, uint lockTime_, uint minSafetyBlocks_
    )
    CommonBridge(_sideBridgeAddress, relayAddress,
                 tokenThisAddresses, tokenSideAddresses,
                 fee_, feeRecipient_,
                 timeframeSeconds_, lockTime_, minSafetyBlocks_)
    {

        // relay uses this event to know from what moment to synchronize the validator set;
        // side bridge contract must be deployed with validator set actual at the time this event was emitted.
        emit Transfer(0, queue);


        emitTestEvent(address(this), msg.sender, 10, true);

    }

    function submitTransfer(PoWProof memory powProof) public onlyRole(ADMIN_ROLE) {

        require(powProof.transfer.event_id == inputEventId + 1);
        inputEventId++;

        CheckPoW_(powProof, sideBridgeAddress);

        //        lockTransfers(events, event_id);
    }

    function setSideBridge(address _sideBridgeAddress) public onlyRole(ADMIN_ROLE) {
        sideBridgeAddress = _sideBridgeAddress;
    }


}
