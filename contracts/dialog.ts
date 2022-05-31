import Dialog from "./dialog_model";
import hardhat from "hardhat";
import config from "./hardhat.config";
import fs from "fs";
import {execSync} from "child_process";
import {ethers} from "hardhat";

// todo add confirmation

async function main() {
    const getNetworkName = (type: string, side: string) => `${type}/${side}`;

    const actions = ["redeploy", "upgrade", "confirmTransaction", "exit"];
    let action;

    while (true) {
        action = await Dialog.askToChooseFromArray(actions, "Choose action:");
        if (action === "exit") return;

        const confirm = await Dialog.confirmation();
        if (confirm) break;
    }


    let networks: any = {};
    for (let net in config.networks) {
        if (net.startsWith("hardhat")) continue;

        const [type, side] = net.split("/");
        if (networks[type] === undefined) {
            networks[type] = [side];
        } else {
            networks[type].push(side);
        }
    }

    const networkType = await Dialog.askToChooseFromArray(Object.keys(networks), "Choose network type");
    const networkSide = await Dialog.askToChooseFromArray(networks[networkType], "Choose network side");

    const path = `./deployments/${networkType}/${networkSide}`;
    let bridgeName = "";

    let bridgeDeploymentExists = false;
    let deploymentsFiles: any;
    if (fs.existsSync(path)) {
        deploymentsFiles = fs.readdirSync(path);
        const pattern = "(.+)_" +
                        `${networkSide.charAt(0).toUpperCase() + networkSide.slice(1)}` +
                        "Bridge\.json";

        for (let i in deploymentsFiles) {
            const res = new RegExp(pattern).exec(deploymentsFiles[i]);
            if (res !== null) {
                bridgeName = res[0].slice(0, -5);  // .json remove
                bridgeDeploymentExists = true;
            }
        }
    }

    const bridgeTypes = ["eth", "bsc"];
    const bridgeType = await Dialog.askToChooseFromArray(bridgeTypes, "Choose bridge type");

    if (action === "redeploy") {
        // todo
        if (bridgeDeploymentExists) {
            const configPath = `./configs/${getConfigName(networkType)}`;
            let obj = require(configPath);
            obj.bridges[bridgeType]["amb"] = "";
            obj.bridges[bridgeType]["side"] = "";
            fs.writeFileSync(configPath, JSON.stringify(obj, null, 2));

            if (networkType === "main") {
                const pattern = "(.+)_" +
                    `${networkSide.charAt(0).toUpperCase() + networkSide.slice(1)}` +
                    "Bridge";

                for (let i in deploymentsFiles) {
                    const res = new RegExp(pattern).exec(deploymentsFiles[i]);
                    if (res !== null) {
                        execSync(`rm ${path}/${res.input}`);
                    }
                }
            } else {
                execSync(`rm -r ${path}`);
            }
        }
        execSync(`./deploy.sh ${bridgeType} ${networkType}`);

        Dialog.output("Successfully redeployed!");
    }

    // const bridgeSides = ["amb", "side"];
    // const bridgeSide = await Dialog.askToChooseBridgeSide(bridgeSides);

    if (action === "upgrade") {
        execSync(`yarn hardhat deploy --network ${getNetworkName(networkType, networkSide)} --tags bridges_${bridgeType}`);

        const {proxyAdmin} = await hardhat.getNamedAccounts();
        const proxyAdminS = await ethers.getSigner(proxyAdmin);

        const contract = await ethers.getContract(bridgeName);
        const Factory = await ethers.getContractFactory("MultiSigWallet");
        const proxy = await Factory.attach(contract.address);

        const lastTransaction = await proxy.getTransactionCount(true, true) - 1;
        await (await proxy.connect(proxyAdminS).confirmTransaction(lastTransaction)).wait();

        Dialog.output("Successfully upgraded!");
    }

    if (action === "confirmTransaction") {
        // todo
    }
}
// if (bridgeDeploymentExists) {
//     const configPath = `./configs/${getConfigName(networkType)}`;
//     let obj = require(configPath);
//     obj.bridges[bridgeType][bridgeSide] = "";
//     fs.writeFileSync(configPath, JSON.stringify(obj, null, 2));
//
//     const pattern = "(.+)_" +
//         `${networkSide.charAt(0).toUpperCase() + networkSide.slice(1)}` +
//         "Bridge";
//
//     for (let i in deploymentsFiles) {
//         const res = new RegExp(pattern).exec(deploymentsFiles[i]);
//         if (res !== null) {
//             execSync(`rm ${path}/${res.input}`);
//         }
//     }
// }
// execSync(`yarn hardhat deploy --network ${getNetworkName(networkType, networkSide)} --tags bridges_${bridgeType}`);

const getConfigName = (networkType: string) => {
    if (networkType === "integr")
        return "config-integr.json"
    else
        return `config-${networkType}net.json`
}

main().catch(reason => {
    console.log(reason);
    process.exitCode = -1;
});