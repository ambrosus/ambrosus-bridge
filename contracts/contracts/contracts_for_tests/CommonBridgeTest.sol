// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonBridge.sol";
import "../checks/CheckReceiptsProof.sol";

contract CommonBridgeTest is CommonBridge {

    // normal constructor can't have calldata args
    function constructor_(CommonStructs.ConstructorArgs calldata args) public {
        __CommonBridge_init(args);

        // used for signature check
        _setupRole(RELAY_ROLE, address(0x295C2707319ad4BecA6b5bb4086617fD6F240CfE));
    }

    function lockTransfersTest(CommonStructs.Transfer[] calldata events, uint eventId) public {
        checkEventId(eventId);  // now its more like submitTransferTest
        lockTransfers(events, eventId);
    }

    function addElementToQueue() public {
        queue.push(CommonStructs.Transfer(address(0), address(0), 100));
    }

    // checkReceiptsProof

    function calcTransferReceiptsHashTest(CommonStructs.TransferProof calldata p, address eventContractAddress) public pure returns (bytes32) {
        return calcTransferReceiptsHash(p, eventContractAddress);
    }

    function checkSignatureTest(bytes32 hash, bytes memory signature) public view returns(address) {
        return ecdsaRecover(hash, signature);
    }


    function FeeCheckTest(address token, bytes calldata signature, uint fee1, uint fee2) public payable {
        feeCheck(token, signature, fee1, fee2, msg.value);
    }

    function getSignatureFeeCheckNumber() public view returns (uint) {
        return signatureFeeCheckNumber;
    }
}
