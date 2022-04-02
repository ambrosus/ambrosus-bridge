import {HardhatRuntimeEnvironment} from "hardhat/types";
import {DeployFunction} from "hardhat-deploy/types";

const func: DeployFunction = async function (hre: HardhatRuntimeEnvironment) {
    if (hre.network.name !== 'hardhat' && hre.network.name != 'rinkeby') return
    const {owner} = await hre.getNamedAccounts();

    await hre.deployments.deploy("MockERC20", {
        contract: "MockERC20",
        from: owner,
        args: [
            ["0xa74393612236d7927edc3947161e7efc2eca7769"] // bridgeAddresses
        ],
        log: true,
    });

};



export default func;
func.tags = ["mocktoken"];
