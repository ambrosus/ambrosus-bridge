// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckReceiptsProof.sol";

contract CommonBridgeTest is CommonBridge, CheckReceiptsProof {
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

    // checkReceiptsProof

    function calcTransferReceiptsHashTest(CommonStructs.TransferProof memory p, address eventContractAddress) public pure returns (bytes32) {
        return calcTransferReceiptsHash(p, eventContractAddress);
    }


}
