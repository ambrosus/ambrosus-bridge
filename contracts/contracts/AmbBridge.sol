// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

contract AmbBridge {
    event Test(uint indexed event_id, Withdraw[] unimportant);
    event newWithdraw(uint indexed event_id, Withdraw[] queue);


    struct Block {
        bytes p1;
        bytes32 prevHashOrReceiptRoot;
        bytes p2;
        bytes difficulty;
        bytes p3;
    }

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


    function TestAll(
        Block[] memory blocks,
        Withdraw[] memory events,
        bytes[] memory proof) public {
        TestReceiptsProof(proof, abi.encode(events), blocks[0].prevHashOrReceiptRoot);

        bytes32 hash = calcReceiptsRoot(proof, abi.encode(events));

        for (uint i = 0; i < blocks.length; i++) {
            require(blocks[i].prevHashOrReceiptRoot == hash, "prevHash or receiptRoot wrong");
            hash = keccak256(abi.encodePacked(blocks[i].p1, blocks[i].prevHashOrReceiptRoot, blocks[i].p2, blocks[i].difficulty, blocks[i].p3));

            TestPoW(hash, blocks[i].difficulty);
        }


        //        require(!TestBloom(bloom, abi.encode(events_hash)), "Failed to verify bloom");

        // TestTransfer(MyTokenAddress, recipient, amount);
        for (uint i = 0; i < events.length; i++) {
            emit DepositEvent(events[i].fromAddress, events[i].toAddress, events[i].amount);
        }
    }

    function TestPoW(bytes32 hash, bytes memory difficulty) internal view {
        require(uint(hash) < bytesToUint(difficulty), "hash must be less than difficulty");
    }

    function eventTest(uint event_id) public {
        emit Test(event_id, [1, 3, 3, 7]);
    }

    function calcReceiptsRoot(bytes[] memory proof, bytes memory eventToSearch) public view returns (bytes32){
        bytes32 el = keccak256(abi.encodePacked(proof[0], eventToSearch, proof[1]));
        bytes memory s;

        for (uint i = 2; i < proof.length - 1; i += 2) {
            s = abi.encodePacked(proof[i], el, proof[i + 1]);
            el = (s.length > 32) ? keccak256(s) : bytes32(s);
        }

        return el;
    }

    function TestReceiptsProof(bytes[] memory proof, bytes memory eventToSearch, bytes32 receiptsRoot) public {
        require(calcReceiptsRoot(proof, eventToSearch) == receiptsRoot, "Failed to verify receipts proof");
    }

    function TestTransfer(address TokenAddress, address recipient, uint256 amount) {
        // require(IERC20(sender.contractAddress).transferFrom(sender.user, receiver, actualPrice), "Fail transfer coins");
        (bool success, bytes memory data) = TokenAddress.call(
            abi.encodeWithSignature('transferFrom(address,address,uint256)',
            address(this),
            recipient,
            amount));

        require(success, "Transfer failed");
    }

    function bytesToUint(bytes memory b) public view returns (uint){
        return uint(bytes32(b)) >> (256 - b.length * 8);
    }
}
