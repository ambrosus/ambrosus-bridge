// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract CommonBridge {
    struct Transfer {
        address tokenAddress;
        address toAddress;
        uint amount;
    }

    // this network to side network
    mapping(address => address) tokenAddresses;

    uint outputEventId;
    uint public inputEventId;

    address public sideBridgeAddress;


    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses)
    {

        sideBridgeAddress = _sideBridgeAddress;

        // initialise tokenAddresses with start values
        tokensAddBatch(tokenThisAddresses, tokenSideAddresses);

        // todo set roles

    }


    function withdraw(address tokenAmbAddress, address toAddress, uint amount) public {
        uint nowTimeframe = block.timestamp / 4 hours;

        if (nowTimeframe != lastTimeframe) {
            emit newWithdraw(outputEventId, queue);
            outputEventId += 1;
            delete queue;
            lastTimeframe = nowTimeframe;
        }

        queue.push(Transfer(tokenAmbAddress, msg.sender, toAddress, amount));
    }


    // todo for relay role
    function acceptBlock() public {

    }


    // todo



    function eventTest(uint event_id) public {
        emit Test(event_id, [1, 3, 3, 7]);
    }


    function Transfer(Transfer memory transfer) {
        require(IERC20(transfer.tokenAddress).transferFrom(sender.user, receiver, actualPrice), "Fail transfer coins");
    }






    // token addressed mapping

    // todo only admin

    function tokensAdd(address tokenThisAddress, address tokenSideAddress) public {
        tokenAddresses[tokenThisAddress] = tokenSideAddress;
    }

    function tokensRemove(address tokenThisAddress) public {
        delete tokenAddresses[tokenThisAddress];
    }

    function tokensAddBatch(address[] memory tokenThisAddresses, address[] memory tokenSideAddresses) public {
        require(tokenThisAddresses.length == tokenSideAddresses.length, "sizes of tokenThisAddresses and tokenSideAddresses must be same");
        uint arrayLength = tokenThisAddresses.length;
        for (uint i = 0; i < arrayLength; i++) {
            tokenAddresses[tokenThisAddresses[i]] = tokenSideAddresses[i];
        }
    }

    function tokensRemoveBatch(address[] memory tokenThisAddresses) public {
        uint arrayLength = tokenThisAddresses.length;
        for (uint i = 0; i < arrayLength; i++) {
            delete tokenAddresses[tokenThisAddresses[i]];
        }
    }




    // utils

    function bytesToUint(bytes memory b) public view returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }


}
