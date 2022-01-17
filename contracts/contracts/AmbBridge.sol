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

    mapping(address => address) fromAmb;

    Withdraw[] queue;

    uint lastTimeframeWithActions;

    uint eventWithdrawId;


    constructor(
        address[] memory ambAddress,
        address[] memory ethAddress) {
        require(ambAddress.length == ethAddress.length, "sizes of ambAddress and ethAddress must be same");

        uint arrayLength = ambAddress.length;
        for (uint i = 0; i < arrayLength; i++) {
            fromAmb[ambAddress[i]] = ethAddress[i];
        }
    }


    function getTimeframe(uint timestamp_) private pure returns (uint) {
        return timestamp_ / uint(4);
    }

    function withdraw(address tokenAmbAddress, address toAddress, uint amount) public {
        if (lastTimeframeWithActions != getTimeframe(block.timestamp)) {
            emit newWithdraw(eventWithdrawId, queue);
            eventWithdrawId += 1;
            delete queue;
        }

        queue.push(Withdraw(tokenAmbAddress, msg.sender, toAddress, amount));
        lastTimeframeWithActions = getTimeframe(block.timestamp);
    }


    function eventTest() public {
//        emit Test(keccak256(abi.encode(queue)), queue);
        delete queue;
    }
}
