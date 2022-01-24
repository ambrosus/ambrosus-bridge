import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name !== 'hardhat' && hre.network.name != 'rinkeby') return
  const {owner} = await hre.getNamedAccounts();

  await hre.deployments.deploy("AmbBridge", {
    contract: "AmbBridge",
    from: owner,
    args: [
      "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", // _sideBridgeAddress
      "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", // relayAddress
      ["0x195c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x395c2707319ad4beca6b5bb4086617fd6f240cfe"], // tokenThisAddresses
      ["0x495c2707319ad4beca6b5bb4086617fd6f240cfe", "0x595c2707319ad4beca6b5bb4086617fd6f240cfe", "0x695c2707319ad4beca6b5bb4086617fd6f240cfe"], // tokenSideAddresses
      1000, // fee
      14400, // timeframeSeconds
      1000, // lockTime
      10 // minSafetyBlocks
    ],
    log: true,
  });

};



export default func;
func.tags = ["ambbridge"];
