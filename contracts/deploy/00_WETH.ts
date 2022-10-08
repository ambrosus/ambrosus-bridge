import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {parseNet, readConfig_} from "./utils/utils";
import {isAddress} from "ethers/lib/utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["eth"]) return;
  const isMainNet = parseNet(hre.network).stage === 'main'
  if (isMainNet) {
    console.log("No need to deploy WETH on mainnet");
    return;
  }

  const {owner} = await hre.getNamedAccounts();

  let configFile = readConfig_(hre.network);

  const weth = configFile.tokens.WETH;

  if (weth === undefined) {
    console.log("wETH is not in config")
    return
  }
  if (isAddress(weth.addresses.eth)) {
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
