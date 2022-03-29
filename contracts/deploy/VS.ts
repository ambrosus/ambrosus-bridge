import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
  if (hre.network.name != 'rinkeby') return
  const {owner} = await hre.getNamedAccounts();

  const validators = ["0x4c9785451bb2CA3E91B350C06bcB5f974cA33F79",
    "0x90B2Ce3741188bCFCe25822113e93983ecacfcA0",
    "0xAccdb7a2268BC4Af0a1898e725138888ba1Ca6Fc"];

  await hre.deployments.deploy("ModifiedValidatorSet", {
    contract: "ModifiedValidatorSet",
    from: owner,
    args: [
      "0x4c9785451bb2CA3E91B350C06bcB5f974cA33F79",
      validators,
      '0xfffffffffffffffffffffffffffffffffffffffe'
    ],
    log: true,
  });

};


export default func;
func.tags = ["vs"];
