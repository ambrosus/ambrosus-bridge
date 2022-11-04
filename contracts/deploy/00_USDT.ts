import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {getSideNetDecimalsOrTokenDenomination, parseNet, readConfig_} from "./utils/utils";
import {isAddress} from "ethers/lib/utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && hre.network.tags["amb"]) return;
  const isMainNet = parseNet(hre.network).stage === 'main'
  if (isMainNet) {
    console.log("No need to deploy USDT on mainnet");
    return;
  }

  const netName = parseNet(hre.network).name;
  const {owner} = await hre.getNamedAccounts();

  let configFile = readConfig_(hre.network);

  const usdt = configFile.tokens.USDT;

  if (usdt === undefined || !usdt.isActive) {
    console.log("USDT is not in config")
    return
  }
  if (usdt.addresses[netName] != "DEPLOY") return;  // already deployed or shouldn't be deployed
  if (isAddress(usdt.addresses[netName])) {
    console.log("USDT already deployed on", netName);
    return;
  }

  const deployResult = await hre.deployments.deploy(usdt.symbol, {
    contract: "MintableERC20",
    from: owner,
    args: [usdt.name, usdt.symbol, getSideNetDecimalsOrTokenDenomination(usdt, netName)],
    log: true,
  });

  usdt.addresses[netName] = deployResult.address;
  configFile.save();


};

export default func;
func.tags = ["USDT", "tokens"];
