import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {parseNet, readConfig_} from "./utils/utils";
import {isAddress} from "ethers/lib/utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && hre.network.tags["amb"]) return;
  const isMainNet = parseNet(hre.network).stage === 'main'
  if (isMainNet) {
    console.log("No need to deploy USDC on mainnet");
    return;
  }

  const netName = parseNet(hre.network).name;
  const {owner} = await hre.getNamedAccounts();

  let configFile = readConfig_(hre.network);

  const usdc = configFile.tokens.USDC;

  if (usdc === undefined) {
    console.log("USDC is not in config")
    return
  }
  if (isAddress(usdc.addresses[netName])) {
    console.log("USDC already deployed on", netName);
    return;
  }

  const deployResult = await hre.deployments.deploy(usdc.symbol, {
    contract: "MintableERC20",
    from: owner,
    args: [usdc.name, usdc.symbol, usdc.denomination],
    log: true,
  });

  usdc.addresses[netName] = deployResult.address;
  configFile.save();


};

export default func;
func.tags = ["USDC", "tokens"];
