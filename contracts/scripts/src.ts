import { ethers, getNamedAccounts, getUnnamedAccounts } from "hardhat";
import { BigNumber, Signer, Wallet } from "ethers";
// import {getNamedSigner, getNamedSigners} from "hardhat-deploy-ethers/internal/helpers";


async function main() {
    // const ownerS = await ethers.getNamedSigner("owner");

    // const b = await ethers.getContract("WETH");
    // await (await b.deposit({value: BigNumber.from("20000000000000000")})).wait();

    const b = await ethers.getContract("BSC_BscBridge");
    // const b = await ethers.getContract("ETH_AmbBridge");
    // const b = await ethers.getContract("BSC_AmbBridge");
    // const token = await ethers.getContract("testSAMB");
    // const b = await ethers.getContract("ETH_EthBridge");
    // console.log(await b.lastProcessedBlock())
    // return
    // console.log("here");
    // for (let i = 0; i < 20; i++) {
    //     console.log(i)
    //     await (await b.wrapWithdraw("0xa74393612236d7927edc3947161E7eFc2EcA7769", "0xf2fc3596318b36f30b38222929d0425db687c88a04a80a4c0fa328435a813fb45e8b39533daa60a9cac0966a9aeb9a73ec63fc2f45ab08f4e456a258b188cc9600", "0x32b3a7486d104eda7", "0x26a052e4477eb70530", {value: BigNumber.from("781983984226631414487")})).wait();
    // }
    // await (await b.withdraw("0x562188D39352ad4CB5B456212C11FABC3a2294c2", "0x295C2707319ad4BecA6b5bb4086617fD6F240CfE", BigNumber.from("1000000000000000"), false, "0x58296f57e5b009596200d396fd7fad5b9dfd2dcd05be9191813201b36d68c94176deae6f2c7f7477eaa809b7af2a1dacf31dcf711fbfb0da426b091a11d9170400", "0x2386f26fc10000", "0x103b2ff6589182", {value: BigNumber.from("14568676809871746")})).wait();
    // await (await b.wrapWithdraw("0x295C2707319ad4BecA6b5bb4086617fD6F240CfE", "0x80e27d64d7c58a999561b680125367e57611d2ddffa6af401c14659dfe1ed00d33a4c92af9c78419cb730201da86c9dba15605de6f62fe3bb9da7e30c9aef19701", "0xde0b6b3a7640000", "0x2d5314d64c5737cc3b", {value: BigNumber.from("837096128744981187643")})).wait();
    // await (await b.wrapWithdraw("0x295C2707319ad4BecA6b5bb4086617fD6F240CfE", "0x803d7b46bc6c8882c3f46ae460f18fc3b2879b91be1410d16c19035f15edd2d026dc6d5a5e86384640b5dd0579faa5f0e118d63c976b29f199e3f1de6874214501", "0x2386f26fc10000", "0xeaef399aa4192", {value: BigNumber.from("20133010952307090")})).wait();
    // await (await b.changeTransferFeeRecipient("адреса")).wait()
    // await (await token.deposit({value: BigNumber.from("255000000000000000000")})).wait();
    // await (await b.withdraw("0xd24c6D48726C0e0b15F55Bb1E1D89D875287694a", "0xD693a3cc5686e74Ca2e72e8120A2F2013B8eE66E", "0x1bc16d674ec80000", false, "0xa1ea7b2dc79e4ac751aae6ab26e72369ad3100e5fbdf8007c81aec96a2bf38d763eb4002d1a6555911ba1d06de7885f65392ce39497e47d094fc8097501a633400", "0xaa84ad107", "0x3f956ef8f72f7f", {value: BigNumber.from("17897273163810054")})).wait();
    await (await b.setFeeProviderAndWatchdogs("0xD693a3cc5686e74Ca2e72e8120A2F2013B8eE66E", ["0xD693a3cc5686e74Ca2e72e8120A2F2013B8eE66E"])).wait()

    // console.log(await b.timeframeSeconds())
    // console.log(await b.lockTime())
    // await (await b.changeTimeframeSeconds(30)).wait();
    // await (await b.changeLockTime(1)).wait();
    // const tx = await b.populateTransaction.setRelaysAndConfirmations([], [], 3)
    // console.log(tx.data)




    // const b = await ethers.getContract("BSC_BscBridge");
    // const relayRole = await b.RELAY_ROLE();
    // console.log(await b.wrapperAddress())
    // console.log(await b.currentEpoch())
    // console.log(await b.minSafetyBlocksValidators())
}

main().then(() => process.exit(0)).catch(error => {
    console.error(error);
    process.exit(1);
});
