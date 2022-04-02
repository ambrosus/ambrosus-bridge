import { HardhatRuntimeEnvironment } from "hardhat/types";
import { DeployFunction } from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name !== "hardhat" && hre.network.name !== "amb") return;
  const { owner } = await hre.getNamedAccounts();
  await hre.deployments.deploy("CommonBridge", {
    from: owner,
    args: [
      {
        sideBridgeAddress: "0x295c2707319ad4beca6b5bb4086617fd6f240cfe",
        relayAddress: "0x295c2707319ad4beca6b5bb4086617fd6f240cfe",
        tokenThisAddresses: [
          "0x195c2707319ad4beca6b5bb4086617fd6f240cfe",
          "0x295c2707319ad4beca6b5bb4086617fd6f240cfe",
          "0x395c2707319ad4beca6b5bb4086617fd6f240cfe",
        ],
        tokenSideAddresses: [
          "0x495c2707319ad4beca6b5bb4086617fd6f240cfe",
          "0x595c2707319ad4beca6b5bb4086617fd6f240cfe",
          "0x695c2707319ad4beca6b5bb4086617fd6f240cfe",
        ],
        fee: 1000,
        feeRecipient: "0x295c2707319ad4beca6b5bb4086617fd6f240cfe",
        timeframeSeconds: 14400,
        lockTime: 1000,
        minSafetyBlocks: 10,
      }
    ],
    log: true,
  });
};

export default func;
func.tags = ["commonbridge"];
