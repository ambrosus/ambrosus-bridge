import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {readConfig} from "./utils/prod_addresses";
import {parseNet} from "./utils/utils";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const isMainNet = parseNet(hre.network).stage === 'main'

  const {owner} = await hre.getNamedAccounts();


  const relays = isMainNet ?
    Object.entries(readConfig())
      .filter(([name, _]) => name.endsWith("AmbBridge")) // only relays from amb side
      .map(([_, addresses]) => addresses.masterRelayAddress) :
    [owner];

  console.log("Deploying faucet, relays:", relays);

  await hre.deployments.deploy("Faucet", {
    contract: "Faucet",
    from: owner,
    args: [relays],
    log: true,
  });
};

export default func;
func.tags = ["faucet"];
func.skip = async (hre: HardhatRuntimeEnvironment) => !hre.network.tags["amb"]; // only amb
