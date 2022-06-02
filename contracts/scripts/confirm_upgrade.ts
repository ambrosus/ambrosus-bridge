import hardhat from "hardhat";
import fs from "fs";

// Do not use this script directly.
// It can be used only by dialog.ts
async function main() {
    const {proxyAdmin} = await hardhat.getNamedAccounts();
    const proxyAdminS = await hardhat.ethers.getSigner(proxyAdmin);

    const bridgeName = fs.readFileSync("./scripts/tmpBridgeName", "utf-8");
    const contract = await hardhat.ethers.getContract(bridgeName);
    const Factory = await hardhat.ethers.getContractFactory("MultiSigWallet");
    const proxy = await Factory.attach(contract.address);

    let transactionCount = await proxy.getTransactionCount(true, true);
    if (transactionCount === 0) {
        console.log("There is no transactions to confirm");
    } else {
        await (await proxy.connect(proxyAdminS).confirmTransaction(transactionCount.toNumber() - 1)).wait();
    }
}

main().catch((reason) => {
    console.log(reason);
    process.exitCode = -1;
});