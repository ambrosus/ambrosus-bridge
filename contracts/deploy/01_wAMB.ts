import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {configPath, readConfig, writeConfig} from "./utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;

  const {owner} = await hre.getNamedAccounts();

  const wAmb = await hre.deployments.deploy("wAMB", {
    contract: "wAMB",
    from: owner,
    args: ["wAMB", "wAmb"],
    log: true,
  });

  const path = configPath(hre.network);
  let configFile = readConfig(path);
  configFile.tokens.wAMB.addresses.amb = wAmb.address;
  writeConfig(path, configFile);

};

export default func;
func.tags = ["wAMB", "tokens"];
