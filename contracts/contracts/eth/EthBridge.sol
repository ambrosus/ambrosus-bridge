// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";
import "../common/CommonBridge.sol";
import "../common/checks/CheckAura.sol";


contract EthBridge is CommonBridge, CheckAura {

    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, address payable feeRecipient_,
        uint timeframeSeconds_, uint lockTime_, uint minSafetyBlocks_,
        address[] memory initialValidators
    )
    CommonBridge(
        _sideBridgeAddress, relayAddress,
        tokenThisAddresses, tokenSideAddresses,
        fee_, feeRecipient_,
        timeframeSeconds_, lockTime_, minSafetyBlocks_
    )
    CheckAura(initialValidators)
    {
        emitTestEvent(address(this), msg.sender, 10, true);
    }

    function submitTransfer(AuraProof memory auraProof, address sideBridgeAddress) public onlyRole(ADMIN_ROLE) {

        require(auraProof.transfer.event_id == inputEventId + 1);
        inputEventId++;

        CheckAura_(auraProof, minSafetyBlocks, sideBridgeAddress);

        //        lockTransfers(events, event_id);
    }
}
