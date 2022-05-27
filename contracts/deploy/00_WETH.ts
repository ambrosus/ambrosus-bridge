import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {networkType, readConfig} from "./utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["eth"]) return;
  if (networkType(hre.network) === 'mainnet') {
    console.log("No need to deploy WETH on mainnet");
    return;
  }

  const {owner} = await hre.getNamedAccounts();

  let configFile = readConfig(hre.network);

  const weth = configFile.tokens.WETH;
  if (!!weth.addresses.eth) {
    console.log("wETH already deployed");
    return;
  }

  const deployResult = await hre.deployments.deploy(weth.symbol, {
    contract: "sAMB",
    from: owner,
    args: [weth.name, weth.symbol],
    log: true,
  });

  weth.addresses.eth = deployResult.address;
  configFile.save();

};

export default func;
func.tags = ["wETH", "tokens"];
