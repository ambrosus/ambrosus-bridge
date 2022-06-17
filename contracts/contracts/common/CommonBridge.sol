// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts-upgradeable/access/AccessControlUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/security/PausableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "./CommonStructs.sol";
import "../tokens/IWrapper.sol";
import "../checks/SignatureCheck.sol";


contract CommonBridge is Initializable, AccessControlUpgradeable, PausableUpgradeable {
    // OWNER_ROLE must be DEFAULT_ADMIN_ROLE because by default only this role able to grant or revoke other roles
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant RELAY_ROLE = keccak256("RELAY_ROLE");

    uint private constant SIGNATURE_FEE_TIMESTAMP = 1800;  // 30 min

    // queue to be pushed in another network
    CommonStructs.Transfer[] queue;

    // locked transfers from another network
    mapping(uint => CommonStructs.LockedTransfers) public lockedTransfers;
    uint public oldestLockedEventId;  // head index of lockedTransfers 'queue' mapping


    // this network to side network token addresses mapping
    mapping(address => address) public tokenAddresses;
    address public wrapperAddress;

    address payable transferFeeRecipient;
    address payable bridgeFeeRecipient;

    address public sideBridgeAddress;
    uint public minSafetyBlocks;
    uint public timeframeSeconds;
    uint public lockTime;

    uint public inputEventId; // last processed event from side network
    uint outputEventId;  // last created event in this network. start from 1 coz 0 consider already processed

    uint public lastTimeframe; // timestamp / timeframeSeconds of latest withdraw

    uint internal signatureFeeCheckNumber;

    event Withdraw(address indexed from, uint eventId, address tokenFrom, address tokenTo, uint amount,
                   uint transferFeeAmount, uint bridgeFeeAmount);
    event Transfer(uint indexed eventId, CommonStructs.Transfer[] queue);
    event TransferSubmit(uint indexed eventId);
    event TransferFinish(uint indexed eventId);

    function __CommonBridge_init(CommonStructs.ConstructorArgs memory args) internal initializer {
        _setupRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _setupRole(RELAY_ROLE, args.relayAddress);
        _setupRole(ADMIN_ROLE, args.adminAddress);

        // initialise tokenAddresses with start values
        _tokensAddBatch(args.tokenThisAddresses, args.tokenSideAddresses);
        wrapperAddress = args.wrappingTokenAddress;

        sideBridgeAddress = args.sideBridgeAddress;
        transferFeeRecipient = args.transferFeeRecipient;
        bridgeFeeRecipient = args.bridgeFeeRecipient;
        minSafetyBlocks = args.minSafetyBlocks;
        timeframeSeconds = args.timeframeSeconds;
        lockTime = args.lockTime;

        oldestLockedEventId = 1;
        outputEventId = 1;

        signatureFeeCheckNumber = 3;

        lastTimeframe = block.timestamp / timeframeSeconds;
    }

    function wrapWithdraw(address toAddress, bytes calldata signature, uint transferFee, uint bridgeFee) public payable {
        address tokenSideAddress = tokenAddresses[wrapperAddress];
        require(tokenSideAddress != address(0), "Unknown token address");

        require(msg.value > transferFee + bridgeFee, "Sent value <= fee");

        uint amount = msg.value - transferFee - bridgeFee;
        feeCheck(wrapperAddress, signature, transferFee, bridgeFee, amount);
        transferFeeRecipient.transfer(transferFee);
        bridgeFeeRecipient.transfer(bridgeFee);

        IWrapper(wrapperAddress).deposit{value : amount}();

        //
        queue.push(CommonStructs.Transfer(tokenSideAddress, toAddress, amount));
        emit Withdraw(msg.sender, outputEventId, address(0), tokenSideAddress, amount, transferFee, bridgeFee);

        withdrawFinish();
    }

    function withdraw(
        address tokenThisAddress,
        address toAddress,
        uint amount,
        bool unwrapSide,
        bytes calldata signature,
        uint transferFee,
        uint bridgeFee
    ) payable public {
        address tokenSideAddress;
        if (unwrapSide) {
            require(tokenAddresses[address(0)] == tokenThisAddress, "Token not point to native token");
            // tokenSideAddress will be 0x0000000000000000000000000000000000000000 - for native token
        } else {
            tokenSideAddress = tokenAddresses[tokenThisAddress];
            require(tokenSideAddress != address(0), "Unknown token address");
        }

        require(msg.value == transferFee + bridgeFee, "Sent value != fee");

        require(amount > 0, "Cannot withdraw 0");

        feeCheck(tokenThisAddress, signature, transferFee, bridgeFee, amount);
        transferFeeRecipient.transfer(transferFee);
        bridgeFeeRecipient.transfer(bridgeFee);

        require(IERC20(tokenThisAddress).transferFrom(msg.sender, address(this), amount), "Fail transfer coins");

        queue.push(CommonStructs.Transfer(tokenSideAddress, toAddress, amount));
        emit Withdraw(msg.sender, outputEventId, tokenThisAddress, tokenSideAddress, amount, transferFee, bridgeFee);

        withdrawFinish();
    }

    function triggerTransfers() public {
        require(queue.length != 0, "Queue is empty");

        emit Transfer(outputEventId++, queue);
        delete queue;
    }

    function withdrawFinish() internal {
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


    function proceedTransfers(CommonStructs.Transfer[] memory transfers) internal {
        for (uint i = 0; i < transfers.length; i++) {

            if (transfers[i].tokenAddress == address(0)) {// native token
                IWrapper(wrapperAddress).withdraw(transfers[i].amount);
                payable(transfers[i].toAddress).transfer(transfers[i].amount);
            } else {// ERC20 token
                require(
                    IERC20(transfers[i].tokenAddress).transfer(transfers[i].toAddress, transfers[i].amount),
                    "Fail transfer coins");
            }

        }
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

        proceedTransfers(transfersLocked.transfers);

        delete lockedTransfers[eventId];
        emit TransferFinish(eventId);

        oldestLockedEventId = eventId + 1;
    }

    // optimized version of unlockTransfers that unlock all transfer that can be unlocked in one call
    function unlockTransfersBatch() public whenNotPaused {
        uint eventId = oldestLockedEventId;
        for (;; eventId++) {
            CommonStructs.LockedTransfers memory transfersLocked = lockedTransfers[eventId];
            if (transfersLocked.endTimestamp == 0 || transfersLocked.endTimestamp > block.timestamp) break;

            proceedTransfers(transfersLocked.transfers);

            delete lockedTransfers[eventId];
            emit TransferFinish(eventId);
        }
        oldestLockedEventId = eventId;
    }

    // delete transfers with passed eventId and all after it
    function removeLockedTransfers(uint eventId) public onlyRole(ADMIN_ROLE) whenPaused {
        require(eventId >= oldestLockedEventId, "eventId must be >= oldestLockedEventId");
        for (; lockedTransfers[eventId].endTimestamp != 0; eventId++)
            delete lockedTransfers[eventId];
        inputEventId = eventId-1; // pretend like we don't receive that event
    }

    function isQueueEmpty() public view returns (bool) {
        return queue.length == 0;
    }


    // admin setters

    function changeMinSafetyBlocks(uint minSafetyBlocks_) public onlyRole(ADMIN_ROLE) {
        minSafetyBlocks = minSafetyBlocks_;
    }

    function changeTransferFeeRecipient(address payable feeRecipient_) public onlyRole(ADMIN_ROLE) {
        transferFeeRecipient = feeRecipient_;
    }

    function changeBridgeFeeRecipient(address payable feeRecipient_) public onlyRole(ADMIN_ROLE) {
        bridgeFeeRecipient = feeRecipient_;
    }

    function changeTimeframeSeconds(uint timeframeSeconds_) public onlyRole(ADMIN_ROLE) {
        lastTimeframe = (lastTimeframe * timeframeSeconds) / timeframeSeconds_;
        timeframeSeconds = timeframeSeconds_;
    }

    function changeLockTime(uint lockTime_) public onlyRole(ADMIN_ROLE) {
        lockTime = lockTime_;
    }

    function changeSignatureFeeCheckNumber(uint signatureFeeCheckNumber_) public onlyRole(ADMIN_ROLE) {
        signatureFeeCheckNumber = signatureFeeCheckNumber_;
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

    function feeCheck(address token, bytes calldata signature, uint transferFee, uint bridgeFee, uint amount) internal {
        bytes32 messageHash;
        address signer;
        uint timestampEpoch = block.timestamp / SIGNATURE_FEE_TIMESTAMP;

        for (uint i = 0; i < signatureFeeCheckNumber; i++) {
            messageHash = keccak256(abi.encodePacked(
                    "\x19Ethereum Signed Message:\n32",
                    keccak256(abi.encodePacked(token, timestampEpoch, transferFee, bridgeFee, amount))
                ));

            signer = ecdsaRecover(messageHash, signature);
            if (hasRole(RELAY_ROLE, signer))
                return;
            timestampEpoch--;
        }
        revert("Signature check failed");
    }

    function checkEventId(uint eventId) internal {
        require(eventId == ++inputEventId, "EventId out of order");
    }

    receive() external payable {}  // need to receive native token from wrapper contract

    uint256[15] private __gap;
}
