import {ethers} from "hardhat";

const ambNetUrl = "http://127.0.0.1:8545";
const ethNetUrl = "http://127.0.0.1:8502";

const ambNet = new ethers.providers.JsonRpcProvider(ambNetUrl);
const ethNet = new ethers.providers.JsonRpcProvider(ethNetUrl);

export const ambSigner = new ethers.Wallet("80f702eb861f36fe8fbbe1a7ccceb04ef7ddef714604010501a5f67c8065d446", ambNet);
export const ethSigner = new ethers.Wallet("51d098d8aee092622149d8f3a79cc7b1ce36ff97fadaa2fbd623c65badeefadc", ethNet);

export const relayAddress = "0x0000000000000000000000000000000000000000"; // todo
export const vsContractAddress = "0x0000000000000000000000000000000000000F00";
