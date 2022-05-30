import {ethers} from "hardhat"
import hardhat from "hardhat";


// ETH_AmbBridge, ETH_EthBridge
const bridgeName = "ETH_AmbBridge";
const transactionNum = -1;

async function main() {
    const {proxyAdmin} = await hardhat.getNamedAccounts();
    const proxyAdminS = await ethers.getSigner(proxyAdmin);

    const contract = await ethers.getContract(bridgeName);
    const Factory = await ethers.getContractFactory("MultiSigWallet");
    const proxy = await Factory.attach(contract.address);

    await (await proxy.connect(proxyAdminS).confirmTransaction(transactionNum)).wait();

    console.log("Confirmed!");
}

main().catch(reason => {
    console.log(reason);
    process.exitCode = -1;
})
