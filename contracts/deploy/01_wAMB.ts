import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;

  const {owner} = await hre.getNamedAccounts();

  await hre.deployments.deploy("wAMB", {
    contract: "wAMB",
    from: owner,
    args: ["wAMB", "wAmb"],
    log: true,
  });

};

export default func;
func.tags = ["wAMB"];
