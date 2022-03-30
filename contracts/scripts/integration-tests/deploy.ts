import fs from "fs";
import {deploy, deployAmbBridge, deployEthBridge} from "../deployBridges";
import {ambSigner, ethSigner, relayAddress, setContractAddressesYml, vsContractAddress} from "./cfg";
import path from "path";



const timeframe = 1;  // each second
const lockTime = 1;
const minSafetyBlocks = 10;


async function main() {
  const ambErc20 = await deploy("MockERC20", ambSigner, []);
  const ethErc20 = await deploy("MockERC20", ethSigner, []);
  fs.writeFileSync(path.resolve(__dirname, "./mockTokensAddresses.json"), JSON.stringify({
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

main().then(() => {
  process.exit(0);  // "gracefully" close websockets :)
}).catch((error) => {
  console.error(error);
  process.exit(1);
});
