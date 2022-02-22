// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";
import "../common/CommonBridge.sol";
import "../common/checks/CheckAura.sol";


contract EthBridge is CommonBridge, CheckAura {

    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timeframeSeconds_, uint lockTime_, uint minSafetyBlocks_,
        address[] memory initialValidators
    )
    CommonBridge(
        _sideBridgeAddress, relayAddress,
        tokenThisAddresses, tokenSideAddresses,
        fee_, timeframeSeconds_, lockTime_, minSafetyBlocks_
    )
    CheckAura(
        initialValidators
    ) {
        emitTestEvent(address(this), msg.sender, 10, true);
    }

    function submitTransfer(AuraProof memory auraProof) public onlyRole(RELAY_ROLE) {

        require(auraProof.transfer.event_id == inputEventId + 1);
        inputEventId++;

        CheckAura_(auraProof, minSafetyBlocks);

        //        lockTransfers(events, event_id);
    }
}
