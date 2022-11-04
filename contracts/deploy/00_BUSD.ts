import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {getSideNetDecimalsOrTokenDenomination, parseNet, readConfig_} from "./utils/utils";
import {isAddress} from "ethers/lib/utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && hre.network.tags["amb"]) return;
  const isMainNet = parseNet(hre.network).stage === 'main'
  if (isMainNet) {
    console.log("No need to deploy BUSD on mainnet");
    return;
  }

  const netName = parseNet(hre.network).name;
  const {owner} = await hre.getNamedAccounts();

  let configFile = readConfig_(hre.network);

  const busd = configFile.tokens.BUSD;

  if (busd === undefined) {
    console.log("BUSD is not in config")
    return
  }
  if (busd.addresses[netName] != "DEPLOY") return;  // already deployed or shouldn't be deployed
  if (isAddress(busd.addresses[netName])) {
    console.log("BUSD already deployed on", netName);
    return;
  }

  const deployResult = await hre.deployments.deploy(busd.symbol, {
    contract: "MintableERC20",
    from: owner,
    args: [busd.name, busd.symbol, getSideNetDecimalsOrTokenDenomination(busd, netName)],
    log: true,
  });

  busd.addresses[netName] = deployResult.address;
  configFile.save();


};

export default func;
func.tags = ["BUSD", "tokens"];
