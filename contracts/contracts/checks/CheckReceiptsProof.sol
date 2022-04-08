// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../common/CommonStructs.sol";


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


    function CalcTransferReceiptsHash(CommonStructs.TransferProof memory p, address eventContractAddress) public pure returns (bytes32) {
        bytes32 el = keccak256(abi.encodePacked(
                p.receipt_proof[0],
                eventContractAddress,
                p.receipt_proof[1],
                toBinary(p.event_id),
                p.receipt_proof[2],
                abi.encode(p.transfers),
                p.receipt_proof[3]
            ));
        return CalcReceiptsHash(p.receipt_proof, el, 4);
        // start from proof[4]
    }


    function toBinary(uint _x) private pure returns (bytes memory) {
        bytes memory b = new bytes(32);
        assembly {
            mstore(add(b, 32), _x)
        }
        uint i;
        for (i = 0; i < 32; i++) {
            if (b[i] != 0) {
                break;
            }
        }
        bytes memory res = new bytes(32 - i);
        for (uint j = 0; j < res.length; j++) {
            res[j] = b[i++];
        }
        return res;
    }
}
