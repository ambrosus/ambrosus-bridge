// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonStructs.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";


contract CheckUntrastless {
    uint public confirmationsThreshold;
    mapping(bytes32 => Confirmations) public confirmations;  // [hash of (eventId, transfers)] => confirmations

    struct Confirmations {
        mapping(address => bool) addresses;
        uint count;
    }

    event Confirmation(address indexed sender, uint indexed eventId, bytes32 hash);

    function checkUntrastless_(uint eventId, CommonStructs.Transfer[] calldata transfers) internal view returns (bool) {
        bytes32 hash = transfersHash(eventId, transfers);
        Confirmations storage conf = confirmations[hash];

        require(conf.addresses[msg.sender] == false, "You have already confirmed this transfer");

        conf.count++;
        emit Confirmation(msg.sender, eventId, hash);

        if (conf.count >= confirmationsThreshold) {
            confirmations[hash];
            return true;
        }
        return false;
    }


    // note: it can return false when transfer already processed and deleted
    function isConfirmedByRelay(address relay, uint eventId, CommonStructs.Transfer[] calldata transfers) returns (bool) {
        bytes32 hash = transfersHash(eventId, transfers);
        return confirmations[hash][msg.sender];
    }

    function transfersHash(uint eventId, CommonStructs.Transfer[] calldata transfers) internal pure returns (bytes32) {
        return keccak256(abi.encodePacked(eventId, transfers));
    }

}
