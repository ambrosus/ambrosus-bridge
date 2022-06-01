import Dialog from "./dialog_model";
import hardhat from "hardhat";
import config from "./hardhat.config";
import fs from "fs";
import {execSync} from "child_process";
import {ethers} from "hardhat";


async function main() {
    const getNetworkName = (type: string, side: string) => `${type}/${side}`;

    const actions = ["redeploy", "upgrade", "confirmTransaction", "exit"];
    let action;
    let fullRedeploy = false;

    while (true) {
        action = await Dialog.askToChooseFromArray(actions, "Choose action:");
        if (action === "exit") return;

        if (action === "redeploy") {
            const msg = "Do you want a full redeploy? (Tokens will be redeployed):"
            fullRedeploy = await Dialog.confirmation(msg);
        }

        const msg = `[${!fullRedeploy ? action : 'full_redeploy'}] Are you sure?`;
        const confirm = await Dialog.confirmation(msg);
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

    const networkType = await Dialog.askToChooseFromArray(Object.keys(networks), "Choose network type:");

    const bridgeTypes = ["eth", "bsc"];
    const bridgeType = await Dialog.askToChooseFromArray(bridgeTypes, "Choose bridge type:");

    if (action === "redeploy") {
        const configPath = `./configs/${getConfigName(networkType)}`;
        let obj = require(configPath);

        obj.bridges[bridgeType]["amb"] = "";
        obj.bridges[bridgeType]["side"] = "";

        Dialog.output(`Beginning of the deploy process: ${new Date().toLocaleTimeString()}`);
        Dialog.output(`If deploy is stuck at ${networkType}/${bridgeType} stage - continue with manual deploy\n`);

        if (fullRedeploy) {
            obj.tokens["SAMB"]["addresses"]["amb"]= "";
            obj.tokens["SAMB"]["addresses"][bridgeType]= "";
            obj.tokens["WETH"]["addresses"]["amb"]= "";
            obj.tokens["WETH"]["addresses"][bridgeType]= "";
            fs.writeFileSync(configPath, JSON.stringify(obj, null, 2));

            for (let i in networks[networkType]) {
                execSync(`rm -r ./deployments/${networkType}/${networks[networkType][i]}`);
            }

            execSync(`./deploy.sh ${bridgeType} ${networkType}`, {stdio: 'inherit'});

        } else {
            fs.writeFileSync(configPath, JSON.stringify(obj, null, 2));
            const pattern = `${bridgeType.toUpperCase()}` + "_(.+)Bridge";

            for (let i in networks[networkType]) {
                const path = `./deployments/${networkType}/${networks[networkType][i]}`;
                const files = fs.readdirSync(path);

                for (let f in files) {
                    const res = new RegExp(pattern).exec(files[f]);
                    if (res !== null)
                        execSync(`rm ${path}/${res.input}`);
                }
            }

            execSync(`yarn hardhat deploy --network ${networkType}/amb --tags bridges_${bridgeType}`, {stdio: 'inherit'});
            execSync(`yarn hardhat deploy --network ${networkType}/${bridgeType} --tags bridges_${bridgeType}`, {stdio: 'inherit'});

            execSync(`yarn hardhat deploy --network ${networkType}/amb --tags bridges_${bridgeType}`, {stdio: 'inherit'});
            Dialog.output(`sideBridgeAddress was set`);

            execSync(`yarn hardhat deploy --network ${networkType}/amb --tags tokens_add_bridges`, {stdio: 'inherit'});
            execSync(`yarn hardhat deploy --network ${networkType}/${bridgeType} --tags tokens_add_bridges`, {stdio: 'inherit'});
            Dialog.output(`Tokens are successfully added`);

        }

        Dialog.output("Successfully redeployed!");
        return;
    }

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

    if (action === "upgrade") {
        execSync(`yarn hardhat deploy --network ${getNetworkName(networkType, networkSide)} --tags bridges_${bridgeType}`);

        const {proxyAdmin} = await hardhat.getNamedAccounts();
        const proxyAdminS = await ethers.getSigner(proxyAdmin);

        const contract = await ethers.getContract(bridgeName);
        const Factory = await ethers.getContractFactory("MultiSigWallet");
        const proxy = await Factory.attach(contract.address);

        const lastTransactionNum = await proxy.getTransactionCount(true, true) - 1;
        await (await proxy.connect(proxyAdminS).confirmTransaction(lastTransactionNum)).wait();

        Dialog.output("Successfully upgraded!");
    }

    if (action === "confirmTransaction") {
        // todo
    }
}

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