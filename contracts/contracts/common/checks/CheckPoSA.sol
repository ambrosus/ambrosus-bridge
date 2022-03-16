// SPDX-License-Identifier: MIT
pragma solidity >=0.7.3 <0.9.0;

import "./CheckReceiptsProof.sol";

contract CheckPoSA is CheckReceiptsProof {
    struct BlockPoSA {

    }

    struct PoSAProof {
        BlockPoSA[] blocks;
        CommonStructs.TransferProof transfer;
    }

    address[] private validatorSet;

    constructor(address[] memory _initialValidators) {
        require(_initialValidators.length > 0, "Length of _initialValidators must be bigger than 0");
        validatorSet = _initialValidators;
    }

    function CheckPoSA_(PoSAProof memory posaProof) public {
        for (uint i = 1; i < posaProof.blocks.length; i++) {
            BlockPoSA memory block_ = posaProof.blocks[i];
        }
    }

    function CheckBlock(BlockPoSA memory block) {}
}
