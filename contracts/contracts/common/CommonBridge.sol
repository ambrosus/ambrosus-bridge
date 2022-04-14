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


    // this network to side network token addresses mapping
    mapping(address => address) public tokenAddresses;
    address public wrapperAddress;


    uint public fee;
    address payable feeRecipient;

    address public sideBridgeAddress;
    uint public minSafetyBlocks;
    uint public timeframeSeconds;
    uint public lockTime;

    uint public inputEventId;
    uint outputEventId;
    uint public oldestLockedEventId;

    uint lastTimeframe;

    event Withdraw(address indexed from, uint event_id, uint feeAmount);
    event Transfer(uint indexed event_id, CommonStructs.Transfer[] queue);
    event TransferSubmit(uint indexed event_id);
    event TransferFinish(uint indexed event_id);


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
        address tokenExternalAddress = tokenAddresses[wrapperAddress];
        require(tokenExternalAddress != address(0), "Unknown token address");

        require(msg.value > fee, "msg.value can't be lesser than fee");
        feeRecipient.transfer(fee);

        uint restOfValue = msg.value - fee;
        IWrapper(wrapperAddress).deposit{value: restOfValue}();

        //
        queue.push(CommonStructs.Transfer(tokenExternalAddress, toAddress, restOfValue));
        emit Withdraw(msg.sender, outputEventId, fee);

        withdraw_finish();
    }

    function withdraw(address tokenAmbAddress, address toAddress, uint amount) payable public {
        address tokenExternalAddress = tokenAddresses[tokenAmbAddress];
        require(tokenExternalAddress != address(0), "Unknown token address");

        require(msg.value == fee, "Sent value != fee");
        feeRecipient.transfer(msg.value);

        require(IERC20(tokenAmbAddress).transferFrom(msg.sender, address(this), amount), "Fail transfer coins");

        queue.push(CommonStructs.Transfer(tokenAmbAddress, toAddress, amount));
        emit Withdraw(msg.sender, outputEventId, fee);

        withdraw_finish();
    }

    function withdraw_finish() internal {
        uint nowTimeframe = block.timestamp / timeframeSeconds;
        if (nowTimeframe != lastTimeframe) {
            emit Transfer(++outputEventId, queue);
            delete queue;

            lastTimeframe = nowTimeframe;
        }
    }


    // locked transfers from another network

    // submitted transfers save here for `lockTime` period
    function lockTransfers(CommonStructs.Transfer[] memory events, uint event_id) internal {
        lockedTransfers[event_id].endTimestamp = block.timestamp + lockTime;
        for (uint i = 0; i < events.length; i++)
            lockedTransfers[event_id].transfers.push(events[i]);
    }

    // after `lockTime` period, transfers can  be unlocked
    function unlockTransfers(uint event_id) public whenNotPaused {
        require(event_id == oldestLockedEventId, "can unlock only oldest event");

        CommonStructs.LockedTransfers memory transfersLocked = lockedTransfers[event_id];
        require(transfersLocked.endTimestamp > 0, "no locked transfers with this id");
        require(transfersLocked.endTimestamp < block.timestamp, "lockTime has not yet passed");

        CommonStructs.Transfer[] memory transfers = transfersLocked.transfers;
        for (uint i = 0; i < transfers.length; i++)
            require(IERC20(transfers[i].tokenAddress).transfer(transfers[i].toAddress, transfers[i].amount), "Fail transfer coins");

        delete lockedTransfers[event_id];
        emit TransferFinish(event_id);

        oldestLockedEventId = event_id+1;
    }

    // optimized version of unlockTransfers that unlock all transfer that can be unlocked in one call
    function unlockTransfersBatch() public whenNotPaused {
        uint event_id = oldestLockedEventId;
        for (;; event_id++) {
            CommonStructs.LockedTransfers memory transfersLocked = lockedTransfers[event_id];
            if (transfersLocked.endTimestamp == 0 || transfersLocked.endTimestamp > block.timestamp) break;

            CommonStructs.Transfer[] memory transfers = transfersLocked.transfers;
            for (uint i = 0; i < transfers.length; i++)
                require(IERC20(transfers[i].tokenAddress).transfer(transfers[i].toAddress, transfers[i].amount), "Fail transfer coins");

            delete lockedTransfers[event_id];
            emit TransferFinish(event_id);
        }
        oldestLockedEventId = event_id;
    }

    // delete transfers with passed event_id and all after it
    function removeLockedTransfers(uint event_id) public onlyRole(ADMIN_ROLE) whenPaused {
        require(event_id >= oldestLockedEventId, "event_id must be >= oldestLockedEventId");
        for ( ;lockedTransfers[event_id].endTimestamp != 0; event_id++)
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

    function checkEventId(uint event_id) internal {
        require(event_id == ++inputEventId, "EventId out of order");
    }
}
