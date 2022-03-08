// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./CommonStructs.sol";
import "hardhat/console.sol";


contract CommonBridge is AccessControl {
    // OWNER_ROLE must be DEFAULT_ADMIN_ROLE because by default only this role able to grant or revoke other roles
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant RELAY_ROLE = keccak256("RELAY_ROLE");


    // queue to be pushed in another network
    CommonStructs.Transfer[] queue;
    // locked transfers from another network
    mapping(uint => CommonStructs.LockedTransfers) public lockedTransfers;


    // this network to side network token addresses mapping
    mapping(address => address) public tokenAddresses;

    uint public fee;
    address payable feeRecipient;

    address public sideBridgeAddress;
    uint public minSafetyBlocks;
    uint public timeframeSeconds;
    uint public lockTime;

    uint public inputEventId;
    uint outputEventId;

    uint lastTimeframe;

    event Withdraw(address indexed from, uint event_id, uint feeAmount);
    event Transfer(uint indexed event_id, CommonStructs.Transfer[] queue);


    constructor(
        address _sideBridgeAddress, address relayAddress,
        address[] memory tokenThisAddresses, address[] memory tokenSideAddresses,
        uint fee_, address payable feeRecipient_, uint timeframeSeconds_, uint lockTime_, uint minSafetyBLocks_)
    {
        _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _setupRole(RELAY_ROLE, relayAddress);

        // initialise tokenAddresses with start values
        _tokensAddBatch(tokenThisAddresses, tokenSideAddresses);

        sideBridgeAddress = _sideBridgeAddress;
        fee = fee_;
        feeRecipient = feeRecipient_;
        minSafetyBlocks = minSafetyBLocks_;
        timeframeSeconds = timeframeSeconds_;
        lockTime = lockTime_;
    }


    // todo remove
    event Test(uint indexed a, address indexed b, string c, uint d);
    function emitTestEvent(address tokenAmbAddress, address toAddress, uint amount, bool transferEvent) public {
        emit Test(1, address(this), "asd", 123);

        queue.push(CommonStructs.Transfer(tokenAmbAddress, toAddress, amount));
        emit Withdraw(msg.sender, outputEventId, fee);

        if (transferEvent) {
            emit Transfer(outputEventId++, queue);
            delete queue;
        }

        emit Test(2, address(msg.sender), "dfg", 456);
    }




    function withdraw(address tokenAmbAddress, address toAddress, uint amount) payable public {
        require(msg.value == fee, "Sent value and fee must be same");
        feeRecipient.send(msg.value);

        queue.push(CommonStructs.Transfer(tokenAmbAddress, toAddress, amount));
        emit Withdraw(msg.sender, outputEventId, fee);

        uint nowTimeframe = block.timestamp / timeframeSeconds;
        if (nowTimeframe != lastTimeframe) {
            emit Transfer(outputEventId++, queue);
            delete queue;

            lastTimeframe = nowTimeframe;
        }
    }






    function lockTransfers(CommonStructs.Transfer[] memory events, uint event_id) internal {
        lockedTransfers[event_id].endTimestamp = block.timestamp + lockTime;
        for (uint i = 0; i < events.length; i++) {
            lockedTransfers[event_id].transfers.push(events[i]);
        }
    }

    function unlockTransfers(uint event_id) public onlyRole(RELAY_ROLE) {
        CommonStructs.LockedTransfers memory transfersLocked = lockedTransfers[event_id];
        require(transfersLocked.endTimestamp < block.timestamp, "lockTime has not yet passed");

        CommonStructs.Transfer[] memory transfers = transfersLocked.transfers;

        for (uint i = 0; i < transfers.length; i++) {
            require(IERC20(transfers[i].tokenAddress).transfer(transfers[i].toAddress, transfers[i].amount), "Fail transfer coins");
        }

        delete lockedTransfers[event_id];
    }


    // admin setters

    function changeMinSafetyBlocks(uint minSafetyBlocks_) public onlyRole(ADMIN_ROLE) {
        minSafetyBlocks = minSafetyBlocks_;
    }

    function changeFee(uint fee_) public onlyRole(ADMIN_ROLE) {
        fee = fee_;
    }

    function changeFeeRecipient(address payable feeRecipient_) public onlyRole(ADMIN_ROLE) {
        feeRecipient = feeRecipient_;
    }

    function changeTimeframeSeconds(uint timeframeSeconds_) public onlyRole(ADMIN_ROLE) {
        timeframeSeconds = timeframeSeconds_;
    }

    function changeLockTime(uint lockTime_) public onlyRole(ADMIN_ROLE) {
        lockTime = lockTime_;
    }


    // token addressed mapping

    function tokensAdd(address tokenThisAddress, address tokenSideAddress) public onlyRole(ADMIN_ROLE) {
        tokenAddresses[tokenThisAddress] = tokenSideAddress;
    }

    function tokensRemove(address tokenThisAddress) public onlyRole(ADMIN_ROLE) {
        delete tokenAddresses[tokenThisAddress];
    }

    function tokensAddBatch(address[] memory tokenThisAddresses, address[] memory tokenSideAddresses) public onlyRole(ADMIN_ROLE) {
        tokensAddBatch(tokenThisAddresses, tokenSideAddresses);
    }

    function _tokensAddBatch(address[] memory tokenThisAddresses, address[] memory tokenSideAddresses) private {
        require(tokenThisAddresses.length == tokenSideAddresses.length, "sizes of tokenThisAddresses and tokenSideAddresses must be same");
        uint arrayLength = tokenThisAddresses.length;
        for (uint i = 0; i < arrayLength; i++)
            tokenAddresses[tokenThisAddresses[i]] = tokenSideAddresses[i];
    }

    function tokensRemoveBatch(address[] memory tokenThisAddresses) public onlyRole(ADMIN_ROLE) {
        uint arrayLength = tokenThisAddresses.length;
        for (uint i = 0; i < arrayLength; i++)
            delete tokenAddresses[tokenThisAddresses[i]];
    }

}
