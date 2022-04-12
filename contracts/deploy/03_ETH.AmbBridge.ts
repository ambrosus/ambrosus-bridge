import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {ethers} from "hardhat";
import {configPath, getTokensPair, readConfig, writeConfig} from "./utils";

const relayAddress = "0x295c2707319ad4beca6b5bb4086617fd6f240cfe" // todo get from something?

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.live && !hre.network.tags["amb"]) return;
  const path = configPath(hre.network);
  let configFile = readConfig(path);

  const {owner} = await hre.getNamedAccounts();
  const [tokensThis, tokensSide] = getTokensPair("amb", "eth", hre.network)

  const deployResult = await hre.deployments.deploy("AmbBridge", {
    from: owner,
    args: [
      {
        sideBridgeAddress: ethers.constants.AddressZero, // amb deployed before eth
        relayAddress: relayAddress,
        tokenThisAddresses: tokensThis,
        tokenSideAddresses: tokensSide,
        fee: 1000,  // todo
        feeRecipient: owner,   // todo
        timeframeSeconds: hre.network.live ? 14400 : 1,
        lockTime: hre.network.live ? 1000 : 1,
        minSafetyBlocks: 10,
      },
      configFile.tokens.wAMB.addresses["amb"],
    ],
    log: true,
  });


  configFile.bridges.eth.amb = deployResult.address;
  writeConfig(path, configFile);

};

export default func;
func.tags = ["bridges"];
