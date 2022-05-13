import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {readConfig} from "./utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;

  const {owner} = await hre.getNamedAccounts();

  let configFile = readConfig(hre.network);

  const samb = configFile.tokens.SAMB;
  if (!!samb.addresses.amb) {
    console.log("sAMB already deployed");
    return;
  }

  const deployResult = await hre.deployments.deploy(samb.symbol, {
    contract: "sAMB",
    from: owner,
    args: [samb.name, samb.symbol],
    log: true,
  });

  samb.addresses.amb = deployResult.address;
  configFile.save();

};

export default func;
func.tags = ["wAMB", "tokens"];
