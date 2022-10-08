import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {parseNet, readConfig_} from "./utils/utils";
import {isAddress} from "ethers/lib/utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["bsc"]) return;
  const isMainNet = parseNet(hre.network).stage === 'main'
  if (isMainNet) {
    console.log("No need to deploy WBNB on mainnet");
    return;
  }

  const {owner} = await hre.getNamedAccounts();

  let configFile = readConfig_(hre.network);

  const wbnb = configFile.tokens.WBNB;

  if (wbnb === undefined) {
    console.log("wBNB is not in config")
    return
  }
  if (isAddress(wbnb.addresses.bsc)) {
    console.log("wBNB already deployed");
    return;
  }

  const deployResult = await hre.deployments.deploy(wbnb.symbol, {
    contract: "sAMB",
    from: owner,
    args: [wbnb.name, wbnb.symbol],
    log: true,
  });

  wbnb.addresses.bsc = deployResult.address;
  configFile.save();

};

export default func;
func.tags = ["wBNB", "tokens"];
