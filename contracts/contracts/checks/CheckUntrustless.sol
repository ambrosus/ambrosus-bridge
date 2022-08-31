// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonStructs.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";


contract CheckUntrustless {
    uint public confirmationsThreshold;

    // [hash of (eventId, transfers)] => [relayAddress] => isConfirmed
    mapping(bytes32 => mapping(address => bool)) public confirmations;

    // using this instead of RELAY_ROLE coz it simpler
    address[] public relays;

    event RelayAdd(address indexed relay);
    event RelayRemove(address indexed relay);
    event ThresholdChange(uint newThreshold);
    event RelayConfirmation(address indexed sender, uint indexed eventId, bytes32 hash);

    // return true if current call reach confirmationsThreshold
    function checkUntrustless_(uint eventId, CommonStructs.Transfer[] calldata transfers) internal returns (bool) {
        require(isRelay(msg.sender), "You not in relay whitelist");

        bytes32 hash = transfersHash(eventId, transfers);
        require(confirmations[hash][msg.sender] == false, "You have already confirmed this");

        uint confirmCount = confirmedCount(hash);
        require(confirmCount < confirmationsThreshold, "Already confirmed");

        confirmations[hash][msg.sender] = true;
        emit RelayConfirmation(msg.sender, eventId, hash);

        // +1 coz current relay confirmed just now
        return confirmCount + 1 >= confirmationsThreshold;
    }


    function isConfirmedByRelay(address relay, uint eventId, CommonStructs.Transfer[] calldata transfers) public view returns (bool) {
        bytes32 hash = transfersHash(eventId, transfers);
        return confirmations[hash][relay];
    }


    function confirmedCount(bytes32 hash) public view returns (uint) {
        uint res;
        for (uint i = 0; i < relays.length; i++)
            if (confirmations[hash][relays[i]])
                res++;
        return res;
    }


    function isRelay(address relay) public view returns (bool){
        for (uint i = 0; i < relays.length; i++)
            if (relays[i] == relay)
                return true;
        return false;
    }

    function getRelays() public view returns (address[] memory) {
        return relays;
    }

    function transfersHash(uint eventId, CommonStructs.Transfer[] calldata transfers) public pure returns (bytes32) {
        bytes memory payload = abi.encodePacked(eventId);

        // i guess we can hash transfers only in this way
        for (uint i = 0; i < transfers.length; i++)
            payload = abi.encodePacked(payload, transfers[i].amount, transfers[i].toAddress, transfers[i].tokenAddress);

        return keccak256(payload);
    }


    function _setRelaysAndConfirmations(address[] memory toRemove, address[] memory toAdd, uint _confirmations) internal {
        for (uint i = 0; i < toRemove.length; i++)
            _removeRelay(toRemove[i]);
        for (uint i = 0; i < toAdd.length; i++)
            _addRelay(toAdd[i]);

        if (confirmationsThreshold != _confirmations) {
            confirmationsThreshold = _confirmations;
            emit ThresholdChange(_confirmations);
        }
    }

    function _removeRelay(address relay) internal {
        for (uint i = 0; i < relays.length; i++) {
            if (relays[i] == relay) {
                relays[i] = relays[relays.length - 1];
                relays.pop();
                emit RelayRemove(relay);
                return;
            }
        }
        revert("Not a relay");
    }

    function _addRelay(address relay) internal {
        require(!isRelay(relay), "Already relay");
        relays.push(relay);
        emit RelayAdd(relay);
    }

    uint256[15] private ___gap;

}
