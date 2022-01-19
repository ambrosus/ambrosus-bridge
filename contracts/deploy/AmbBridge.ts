import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name !== 'hardhat' && hre.network.name != 'rinkeby') return
  const {owner} = await hre.getNamedAccounts();

  await hre.deployments.deploy("AmbBridge", {
    contract: "AmbBridge",
    from: owner,
    args: ["0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe",
      ["0x195c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x395c2707319ad4beca6b5bb4086617fd6f240cfe"],
      ["0x495c2707319ad4beca6b5bb4086617fd6f240cfe", "0x595c2707319ad4beca6b5bb4086617fd6f240cfe", "0x695c2707319ad4beca6b5bb4086617fd6f240cfe"],
      1000, 14400],
    log: true,
  });

};



export default func;
func.tags = ["ambbridge"];
