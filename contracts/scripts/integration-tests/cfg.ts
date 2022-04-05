import {ethers} from "hardhat";
import fs from "fs";

import config from "../../../relay/configs/integr.json";

//


// websocket need for events subscribe
const ambNet = new ethers.providers.WebSocketProvider(config.network.eth.wsUrl);
const ethNet = new ethers.providers.WebSocketProvider(config.network.eth.wsUrl);


// accounts with money; details in /tests/README.md
const ambPk = "80f702eb861f36fe8fbbe1a7ccceb04ef7ddef714604010501a5f67c8065d446";
// todo add amb validators pk
const ethPk = "51d098d8aee092622149d8f3a79cc7b1ce36ff97fadaa2fbd623c65badeefadc";
// must be equal to pk set in relay env var
const relayPk = "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5";


export const ambSigner = new ethers.Wallet(ambPk, ambNet);
export const ethSigner = new ethers.Wallet(ethPk, ethNet);

export const relayAddress = new ethers.Wallet(relayPk).address;

export const vsContractAddress = config.network.amb.vsContractAddr;

export const options = {gasLimit: 800000};  // amb & eth gas estimator broken

//

export function setContractAddressesYml(ambAddress: string, ethAddress: string) {
  config.network.amb.contractAddr = ambAddress;
  config.network.eth.contractAddr = ethAddress;
  fs.writeFileSync("../../../relay/configs/integr.json", JSON.stringify(config));
}

// wait for transaction to be mined
export async function w(call: Promise<any>): Promise<any> {
  const tx = await (await call).wait();
  console.log('Transaction mined');
  return tx;
}

