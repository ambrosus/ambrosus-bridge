pragma solidity ^0.8.6;

library CommonStructs {
    struct Transfer {
        address tokenAddress;
        address toAddress;
        uint amount;
    }

    struct TransferProof {
        bytes[] receipt_proof;
        uint event_id;
        Transfer[] transfers;
    }

    struct LockedTransfers {
        Transfer[] transfers;
        uint endTimestamp;
    }


}
