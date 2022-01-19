// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";


contract CommonBridge is AccessControl {
    // OWNER_ROLE must be DEFAULT_ADMIN_ROLE because by default only this role able to grant or revoke other roles
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant RELAY_ROLE = keccak256("RELAY_ROLE");

    event Test(uint indexed event_id, uint unimportant);
    event TransferEvent(uint indexed event_id, Transfer[] queue);


    struct Transfer {
        address tokenAddress;
        address toAddress;
        uint amount;
    }


    // this network to side network
    mapping(address => address) tokenAddresses;

    uint fee;

    uint lastTimeframe;

    Transfer[] queue;

    uint outputEventId;
    uint public inputEventId;

    address public sideBridgeAddress;


    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_)
    {
        sideBridgeAddress = _sideBridgeAddress;

        // initialise tokenAddresses with start values
        tokensAddBatch(tokenThisAddresses, tokenSideAddresses);

        _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);

        fee = fee_;
    }


    function withdraw(address tokenAmbAddress, address toAddress, uint amount) payable public {
        require(msg.value == fee, "Sent value and fee must be same");

        uint nowTimeframe = block.timestamp / 4 hours;

        if (nowTimeframe != lastTimeframe) {
            emit TransferEvent(outputEventId, queue);
            outputEventId += 1;
            delete queue;
            lastTimeframe = nowTimeframe;
        }

        queue.push(Transfer(tokenAmbAddress, toAddress, amount));
    }


    function acceptBlock() public onlyRole(RELAY_ROLE) {

    }


    // todo



    function eventTest(uint event_id) public {
        emit Test(event_id, 1337);
    }


    function Transfer_(Transfer memory transfer) public {
        require(IERC20(transfer.tokenAddress).transferFrom(msg.sender, transfer.toAddress, transfer.amount), "Fail transfer coins");
    }






    // token addressed mapping

    // todo only admin

    function tokensAdd(address tokenThisAddress, address tokenSideAddress) public onlyRole(ADMIN_ROLE) {
        tokenAddresses[tokenThisAddress] = tokenSideAddress;
    }

    function tokensRemove(address tokenThisAddress) public onlyRole(ADMIN_ROLE) {
        delete tokenAddresses[tokenThisAddress];
    }

    function tokensAddBatch(address[] memory tokenThisAddresses, address[] memory tokenSideAddresses) public onlyRole(ADMIN_ROLE) {
        require(tokenThisAddresses.length == tokenSideAddresses.length, "sizes of tokenThisAddresses and tokenSideAddresses must be same");
        uint arrayLength = tokenThisAddresses.length;
        for (uint i = 0; i < arrayLength; i++) {
            tokenAddresses[tokenThisAddresses[i]] = tokenSideAddresses[i];
        }
    }

    function tokensRemoveBatch(address[] memory tokenThisAddresses) public onlyRole(ADMIN_ROLE) {
        uint arrayLength = tokenThisAddresses.length;
        for (uint i = 0; i < arrayLength; i++) {
            delete tokenAddresses[tokenThisAddresses[i]];
        }
    }


    function changeFee(uint fee_) public onlyRole(ADMIN_ROLE) {
        fee = fee_;
    }
}
