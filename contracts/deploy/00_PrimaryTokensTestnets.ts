import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";
import {parseNet, readConfig_} from "./utils/utils";
import {isTokenPrimary, isTokenIsNativeCoinWrapper} from "./utils/config";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  const netName = parseNet(hre.network).name;
  const { owner } = await hre.getNamedAccounts();

  let configFile = readConfig_(hre.network);

  for (const token of Object.values(configFile.tokens)) {
    if (!token.isActive) {
      console.log(token.symbol, "is not active");
      continue;
    }
    if (token.networks[netName]?.address != "DEPLOY") {
      console.log(token.symbol, "already or should not be deployed.");
      continue;
    }
    if (!isTokenPrimary(token, netName)) {
      console.log(token.symbol, "is synthetic token, skip it.");
      continue;
    }

    let address: string;
    if (isTokenIsNativeCoinWrapper(token, netName)) {
      console.log("Deploying test wrappable token", token.symbol, "on", netName);
      ({ address } = await hre.deployments.deploy(token.symbol, {
        contract: "sAMB",
        from: owner,
        args: [token.name, token.symbol],
        log: true,
      }));
    } else {
      console.log("Deploying test mintable token", token.symbol, "on", netName);
      ({ address } = await hre.deployments.deploy(token.symbol, {
        contract: "MintableERC20",
        from: owner,
        args: [token.name, token.symbol, token.networks[netName].denomination],
        log: true,
      }));
    }

    token.networks[netName].address = address;
    configFile.save();
  }
};

export default func;
func.tags = ["primary_tokens_testnets"];
func.skip = async (hre: HardhatRuntimeEnvironment) => parseNet(hre.network).stage == "main"; // only testnets
