// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckReceiptsProof.sol";

contract CommonBridgeTest is CommonBridge {
    constructor(
        CommonStructs.ConstructorArgs memory args
    ) {
        __CommonBridge_init(args);
    }

    function getLockedTransferTest(uint eventId) public view returns (CommonStructs.LockedTransfers memory) {
        return lockedTransfers[eventId];
    }

    function lockTransfersTest(CommonStructs.Transfer[] memory events, uint eventId) public {
        lockTransfers(events, eventId);
    }

    function getOutputEventId() public view returns(uint) {
        return outputEventId;
    }

    function addElementToQueue() public {
        queue.push(CommonStructs.Transfer(address(0), address(0), 100));
    }

    // checkReceiptsProof

    function calcTransferReceiptsHashTest(CommonStructs.TransferProof memory p, address eventContractAddress) public pure returns (bytes32) {
        return calcTransferReceiptsHash(p, eventContractAddress);
    }
}
