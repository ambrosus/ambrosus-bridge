// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "hardhat/console.sol";


contract EthBridge {


    struct Block {
        bytes header;
        uint timeOffset;
        uint timeLength;
    }


    // todo maybe not always
    uint constant BloomOffset = 180 + 12;
    uint constant BloomLength = 256;

    bytes1 constant b1 = bytes1(0x01);


    constructor() {}



    function BlockHashTest(bytes memory block) public view returns (bytes32) {
        return keccak256(block);
    }


//    function TimestampTest(Block memory block) public view returns (uint) {
//
//        uint block_time;
//        assembly {
//            block_time := mload(add(block.header, block.timeOffset))
//        }
//        block_time = block_time >> (32 - block.timeLength);
//
//
//        return block_time;
//    }


    function BloomTest(bytes memory bloom, bytes memory topicHash) public view returns (bool) {
        bytes32 hashbuf = keccak256(topicHash);

        // todo asm
        bytes1 v1 = b1 << uint8(hashbuf[1] & 0x07);
        bytes1 v2 = b1 << uint8(hashbuf[3] & 0x07);
        bytes1 v3 = b1 << uint8(hashbuf[5] & 0x07);

        uint i1 = BloomLength - uint((Uint16(hashbuf[0], hashbuf[1]) & 0x7ff) >> 3) - 1;
        uint i2 = BloomLength - uint((Uint16(hashbuf[2], hashbuf[3]) & 0x7ff) >> 3) - 1;
        uint i3 = BloomLength - uint((Uint16(hashbuf[4], hashbuf[5]) & 0x7ff) >> 3) - 1;

        return
        v1 == v1 & bloom[i1] &&
        v2 == v2 & bloom[i2] &&
        v3 == v3 & bloom[i3];
    }


    function Uint16(bytes2 a, bytes1 b) internal view returns (uint16) {
        return uint16(a) + uint16(uint8(b));
    }


}
