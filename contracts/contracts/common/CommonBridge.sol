// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/access/AccessControl.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "./CommonStructs.sol";
import "../tokens/IWrapper.sol";



contract CommonBridge is AccessControl, Pausable {
    // OWNER_ROLE must be DEFAULT_ADMIN_ROLE because by default only this role able to grant or revoke other roles
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant RELAY_ROLE = keccak256("RELAY_ROLE");


    // queue to be pushed in another network
    CommonStructs.Transfer[] queue;

    // locked transfers from another network
    mapping(uint => CommonStructs.LockedTransfers) public lockedTransfers;
    uint public oldestLockedEventId = 1;  // head index of lockedTransfers 'queue' mapping


    // this network to side network token addresses mapping
    mapping(address => address) public tokenAddresses;
    address public wrapperAddress;

    uint public fee;
    address payable feeRecipient;

    address public sideBridgeAddress;
    uint public minSafetyBlocks;
    uint public timeframeSeconds;
    uint public lockTime;

    uint public inputEventId; // last processed event from side network
    uint outputEventId = 1;  // last created event in this network. start from 1 coz 0 consider already processed

    uint lastTimeframe;

    event Withdraw(address indexed from, uint eventId, uint feeAmount);
    event Transfer(uint indexed eventId, CommonStructs.Transfer[] queue);
    event TransferSubmit(uint indexed eventId);
    event TransferFinish(uint indexed eventId);


    constructor(CommonStructs.ConstructorArgs memory args)
    {
        _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _setupRole(RELAY_ROLE, args.relayAddress);
        _setupRole(ADMIN_ROLE, args.adminAddress);

        // initialise tokenAddresses with start values
        _tokensAddBatch(args.tokenThisAddresses, args.tokenSideAddresses);
        wrapperAddress = args.wrappingTokenAddress;

        sideBridgeAddress = args.sideBridgeAddress;
        fee = args.fee;
        feeRecipient = args.feeRecipient;
        minSafetyBlocks = args.minSafetyBlocks;
        timeframeSeconds = args.timeframeSeconds;
        lockTime = args.lockTime;
    }


    function wrap_withdraw(address toAddress) public payable {
        address tokenSideAddress = tokenAddresses[wrapperAddress];
        require(tokenSideAddress != address(0), "Unknown token address");

        require(msg.value > fee, "msg.value can't be lesser than fee");
        feeRecipient.transfer(fee);

        uint restOfValue = msg.value - fee;
        IWrapper(wrapperAddress).deposit{value: restOfValue}();

        //
        queue.push(CommonStructs.Transfer(tokenSideAddress, toAddress, restOfValue));
        emit Withdraw(msg.sender, outputEventId, fee);

        withdraw_finish();
    }

    function withdraw(address tokenThisAddress, address toAddress, uint amount) payable public {
        address tokenSideAddress = tokenAddresses[tokenThisAddress];
        require(tokenSideAddress != address(0), "Unknown token address");

        require(msg.value == fee, "Sent value != fee");
        feeRecipient.transfer(msg.value);

        require(IERC20(tokenThisAddress).transferFrom(msg.sender, address(this), amount), "Fail transfer coins");

        queue.push(CommonStructs.Transfer(tokenSideAddress, toAddress, amount));
        emit Withdraw(msg.sender, outputEventId, fee);

        withdraw_finish();
    }

    function withdraw_finish() internal {
        uint nowTimeframe = block.timestamp / timeframeSeconds;
        if (nowTimeframe != lastTimeframe) {
            emit Transfer(outputEventId++, queue);
            delete queue;

            lastTimeframe = nowTimeframe;
        }
    }


    // locked transfers from another network

    function getLockedTransfers(uint eventId) public view returns (CommonStructs.LockedTransfers memory) {
        return lockedTransfers[eventId];
    }

    // submitted transfers save here for `lockTime` period
    function lockTransfers(CommonStructs.Transfer[] memory events, uint eventId) internal {
        lockedTransfers[eventId].endTimestamp = block.timestamp + lockTime;
        for (uint i = 0; i < events.length; i++)
            lockedTransfers[eventId].transfers.push(events[i]);
    }

    // after `lockTime` period, transfers can  be unlocked
    function unlockTransfers(uint eventId) public whenNotPaused {
        require(eventId == oldestLockedEventId, "can unlock only oldest event");

        CommonStructs.LockedTransfers memory transfersLocked = lockedTransfers[eventId];
        require(transfersLocked.endTimestamp > 0, "no locked transfers with this id");
        require(transfersLocked.endTimestamp < block.timestamp, "lockTime has not yet passed");

        CommonStructs.Transfer[] memory transfers = transfersLocked.transfers;
        for (uint i = 0; i < transfers.length; i++)
            require(IERC20(transfers[i].tokenAddress).transfer(transfers[i].toAddress, transfers[i].amount), "Fail transfer coins");

        delete lockedTransfers[eventId];
        emit TransferFinish(eventId);

        oldestLockedEventId = eventId+1;
    }

    // optimized version of unlockTransfers that unlock all transfer that can be unlocked in one call
    function unlockTransfersBatch() public whenNotPaused {
        uint eventId = oldestLockedEventId;
        for (;; eventId++) {
            CommonStructs.LockedTransfers memory transfersLocked = lockedTransfers[eventId];
            if (transfersLocked.endTimestamp == 0 || transfersLocked.endTimestamp > block.timestamp) break;

            CommonStructs.Transfer[] memory transfers = transfersLocked.transfers;
            for (uint i = 0; i < transfers.length; i++)
                require(IERC20(transfers[i].tokenAddress).transfer(transfers[i].toAddress, transfers[i].amount), "Fail transfer coins");

            delete lockedTransfers[eventId];
            emit TransferFinish(eventId);
        }
        oldestLockedEventId = eventId;
    }

    // delete transfers with passed eventId and all after it
    function removeLockedTransfers(uint eventId) public onlyRole(ADMIN_ROLE) whenPaused {
        require(eventId >= oldestLockedEventId, "eventId must be >= oldestLockedEventId");
        for ( ;lockedTransfers[eventId].endTimestamp != 0; eventId++)
            delete lockedTransfers[eventId];
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
        _tokensAddBatch(tokenThisAddresses, tokenSideAddresses);
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

    // pause

    function pause() public onlyRole(ADMIN_ROLE) {
        _pause();
    }

    function unpause() public onlyRole(ADMIN_ROLE) {
        _unpause();
    }

    // internal

    function checkEventId(uint eventId) internal {
        require(eventId == ++inputEventId, "EventId out of order");
    }
}
