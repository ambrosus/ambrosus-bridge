// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract AmbBridge {
    // event Test(bytes32 indexed withdraws_hash, Withdraw[] withdraws);
    event newWithdraw(uint event_id, Withdraw[] queue);


    struct Withdraw {
        address tokenExtAddress;
        address fromAddress;
        address toAddress;
        uint amount;
    }

    Withdraw[] queue;

    uint lastTimeframeWithActions;

    uint eventId;

    constructor() {}


    function getTimeframe(uint timestamp_) private view returns (uint) {
        return timestamp_ / uint(4);
    }

    function withdraw(address tokenAmbAddress, address toAddress, uint amount) public {
        if (block.timestamp != getTimeframe(lastTimeframeWithActions)) {
            emit newWithdraw(eventId, queue);
            delete queue;
        }

        queue.push(Withdraw(tokenAmbAddress, msg.sender, toAddress, amount));
        lastTimeframeWithActions = getTimeframe(block.timestamp);

        eventId += 1;
    }


    function eventTest() public {
//        emit Test(keccak256(abi.encode(queue)), queue);
        delete queue;
    }
}
