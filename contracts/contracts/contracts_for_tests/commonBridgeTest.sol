// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckReceiptsProof.sol";

contract CommonBridgeTest is CommonBridge, CheckReceiptsProof {
    constructor(
        CommonStructs.ConstructorArgs memory args
    ) {
        __CommonBridge_init(args);

        _setupRole(RELAY_ROLE, address(0x295C2707319ad4BecA6b5bb4086617fD6F240CfE));
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

    // checkReceiptsProof

    function calcTransferReceiptsHashTest(CommonStructs.TransferProof memory p, address eventContractAddress) public pure returns (bytes32) {
        return calcTransferReceiptsHash(p, eventContractAddress);
    }


    function FeeCheckTest(address token, bytes calldata signature, uint fee1, uint fee2) public {
        feeCheck(token, signature, fee1, fee2);
    }

    function getSignatureFeeCheckNumber() public view returns(uint) {
        return signatureFeeCheckNumber;
    }
}
