import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name !== 'hardhat' && hre.network.name != 'rinkeby') return
  const {owner} = await hre.getNamedAccounts();

  await hre.deployments.deploy("EthBridge", {
    contract: "EthBridge",
    from: owner,
    args: [owner],
    log: true,
  });

};



export default func;
func.tags = ["ethbridge"];
