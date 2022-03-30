import fs from "fs";
import {deploy, deployAmbBridge, deployEthBridge} from "../deployBridges";
import {ambSigner, ethSigner, relayAddress, vsContractAddress} from "./cfg";

const YAML = require("js-yaml");


const timeframe = 1;  // each second
const lockTime = 1;
const minSafetyBlocks = 10;


async function main() {


  const ambErc20 = await deploy("MockERC20", ambSigner, []);
  const ethErc20 = await deploy("MockERC20", ethSigner, []);
  fs.writeFileSync("../mockTokensAddresses.json", JSON.stringify({
    ambErc20Address: ambErc20.address,
    ethErc20Address: ethErc20.address,
  }));


  const ambBridge = await deployAmbBridge(ambSigner, {
    sideBridgeAddress: null,
    relayAddress,
    tokens: {[ambErc20.address]: ethErc20.address},
    fee: 10,
    feeRecipient: relayAddress,
    timeframe, lockTime, minSafetyBlocks
  });


  const ethBridge = await deployEthBridge(ethSigner, {
    sideBridgeAddress: ambBridge.address,
    relayAddress,
    tokens: {[ethErc20.address]: ambErc20.address},
    fee: 1000,
    feeRecipient: relayAddress,
    timeframe, lockTime, minSafetyBlocks
  }, vsContractAddress, ambSigner);


  await ambBridge.setSideBridge(ethBridge.address, {gasLimit: 30_000});  // auto gas exceed network gasLimit

  // todo use real erc20 contract
  // todo set bridge addresses to erc20 contracts

  setContractAddressesYml(ambBridge.address, ethBridge.address);
}


function setContractAddressesYml(ambAddress: String, ethAddress: String) {
  const path = "../relay/configs/dev.yml";
  const raw = fs.readFileSync(path);
  const data = YAML.load(raw);

  data.network.amb["contract-addr"] = ambAddress;
  data.network.eth["contract-addr"] = ethAddress;

  const yaml = YAML.dump(data);
  fs.writeFileSync(path, yaml);
}

main().then(() => {
}).catch((error) => {
  console.error(error);
  process.exit(1);
});
