// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../common/CheckPoW.sol";

contract AmbBridge is CommonBridge, CheckPoW {
    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, uint timeframeSeconds_, uint lockTime_, uint minSafetyBlocks_
    )
    CommonBridge(_sideBridgeAddress, relayAddress,
        tokenThisAddresses, tokenSideAddresses,
        fee_, timeframeSeconds_, lockTime_, minSafetyBlocks_
    ) {
        emitTestEvent(address(this), msg.sender, 10);
    }

    function submitTransfer(PoWProof memory powProof) public onlyRole(RELAY_ROLE) {

        require(powProof.transfer.event_id == inputEventId + 1);
        inputEventId++;

        CheckPoW_(powProof);

//        lockTransfers(events, event_id);
    }
}
