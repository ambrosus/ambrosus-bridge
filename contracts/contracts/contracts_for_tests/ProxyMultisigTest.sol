pragma solidity ^0.8.0;

contract ProxyMultisigTest {
    bytes32 public lastProcessedBlock;

    event lastProcessedBlockHasChanged();

    function updateLastProcessedBlock(bytes32 lastProcessedBlock_) public {
        require(msg.sender == address(this), "Only this contract can call this function");
        lastProcessedBlock = lastProcessedBlock_;
        emit lastProcessedBlockHasChanged();
    }
}
