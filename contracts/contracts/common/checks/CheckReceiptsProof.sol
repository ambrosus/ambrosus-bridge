// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../CommonStructs.sol";


contract CheckReceiptsProof {
    // check readme for focs
    function CalcReceiptsHash(bytes[] memory proof, bytes32 el, uint proofStart) public pure returns (bytes32) {
        bytes memory s;

        for (uint i = proofStart; i < proof.length; i += 2) {
            s = abi.encodePacked(proof[i], el, proof[i + 1]);
            el = (s.length > 32) ? keccak256(s) : bytes32(s);
        }

        return el;
    }


    function CalcTransferReceiptsHash(CommonStructs.TransferProof memory p, address eventContractAddress) public view returns (bytes32) {
        bytes32 el = keccak256(abi.encodePacked(
                p.receipt_proof[0],
                eventContractAddress,
                p.receipt_proof[1],
                p.event_id,
                p.receipt_proof[2],
                abi.encode(p.transfers),
                p.receipt_proof[3]
            ));
        return CalcReceiptsHash(p.receipt_proof, el, 4);
        // start from proof[4]
    }

}
