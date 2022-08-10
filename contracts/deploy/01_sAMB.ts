import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {parseNet, readConfig_} from "./utils/utils";
import {isAddress} from "ethers/lib/utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  let configFile = readConfig_(hre.network);

  const samb = configFile.tokens.SAMB;

  if (isAddress(samb.addresses.amb)) {
      console.log("sAMB already deployed");
      return;
  }

  const {owner} = await hre.getNamedAccounts();

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
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["amb"]; // only amb
