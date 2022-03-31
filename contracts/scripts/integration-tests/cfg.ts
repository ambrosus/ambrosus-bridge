import {ethers} from "hardhat";
import fs from "fs";

const yaml = require("js-yaml");

const yamlPath = "../relay/configs/integr.yml";
const config = loadYaml();

//


// websocket need for events subscribe
const ambNet = new ethers.providers.WebSocketProvider(config.network.amb["ws-url"]);
const ethNet = new ethers.providers.WebSocketProvider(config.network.eth["ws-url"]);


// accounts with money; details in /tests/README.md
const ambPk = "80f702eb861f36fe8fbbe1a7ccceb04ef7ddef714604010501a5f67c8065d446";
// todo add amb validators pk
const ethPk = "51d098d8aee092622149d8f3a79cc7b1ce36ff97fadaa2fbd623c65badeefadc";
// must be equal to pk set in relay env var
const relayPk = "34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5";


export const ambSigner = new ethers.Wallet(ambPk, ambNet);
export const ethSigner = new ethers.Wallet(ethPk, ethNet);

export const relayAddress = new ethers.Wallet(relayPk).address;

export const vsContractAddress = config.network.amb["vs-contract-addr"];


//

export function loadYaml() {
  return yaml.load(fs.readFileSync(yamlPath));
}

export function setContractAddressesYml(ambAddress: String, ethAddress: String) {
  const data = loadYaml()
  data.network.amb["contract-addr"] = ambAddress;
  data.network.eth["contract-addr"] = ethAddress;
  fs.writeFileSync(yamlPath, yaml.dump(data));
}
