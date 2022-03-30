import chai from "chai";
import {ethers} from "hardhat";
import 'mocha';
import fs from "fs";
import {ambSigner, ethSigner, relayAddress} from "./cfg";
import mockTokens from "./mockTokensAddresses.json";
import {Contract} from "ethers";


const YAML = require("js-yaml");


chai.should();
export const expect = chai.expect;


describe("Contract", function () {
  this.timeout(5 * 60 * 1000); // 5 minutes

  const ambOptions = {gasLimit: 800000};  // amb gas estimator broken

  let ambBridge: Contract;
  let ethBridge: Contract;
  let ambToken: Contract;
  let ethToken: Contract;

  before(async () => {
    const bridges = getContractAddressYml();

    ambBridge = await ethers.getContractAt("AmbBridge", bridges.ambAddress, ambSigner);
    ethBridge = await ethers.getContractAt("EthBridge", bridges.ethAddress, ethSigner);
    ambToken = await ethers.getContractAt("MockERC20", mockTokens.ambErc20Address, ambSigner);
    ethToken = await ethers.getContractAt("MockERC20", mockTokens.ethErc20Address, ethSigner);
  });

  // it('Amb -> Eth', async () => {
  // });


  it('Eth -> Amb', async () => {
    // setup relay
    console.log("Setup relay");
    const relayRole = await ethBridge.callStatic.ADMIN_ROLE();
    await w(ethBridge.grantRole(relayRole, relayAddress));
    await w(ambBridge.grantRole(relayRole, relayAddress, ambOptions));

    await w(ethSigner.sendTransaction({to: relayAddress, value: ethers.utils.parseEther("0.1")}));
    await w(ambSigner.sendTransaction({to: relayAddress, value: ethers.utils.parseEther("1"), ...ambOptions}));


    // mint tokens
    console.log("Mint tokens");
    await w(ambToken.mint(ambSigner.address, 10, ambOptions));
    await w(ethToken.mint(ethSigner.address, 10));

    await w(ethToken.increaseAllowance(ethBridge.address, 10));
    await w(ambToken.increaseAllowance(ambBridge.address, 10, ambOptions));

    const ethBefore = await ethToken.balanceOf(ethSigner.address);
    const ambBefore = await ambToken.balanceOf(ambSigner.address);

    // check mint working
    expect(ethBefore).gte(10);
    expect(ambBefore).gte(10);


    // withdraw
    console.log("Withdraw");
    const fee = await ethBridge.fee();
    await w(ethBridge.withdraw(ethToken.address, ambSigner.address, 5, {value: fee}));
    await sleep(2000); // waiting for 100% new timeframe (timeframe is 1 second)
    await w(ethBridge.withdraw(ethToken.address, ambSigner.address, 1, {value: fee}));  // must emit event, todo check

    // wait for minSafetyBlocks confirmations
    console.log(`Wait for confirmations`);
    const minSafetyBlocks = await ambBridge.minSafetyBlocks();
    const currentBlock = (await ethSigner.provider.getBlock('latest')).number;
    await new Promise<void>((resolve) => {
      ethSigner.provider.on('block', async (block) => {
        if (block - currentBlock > minSafetyBlocks) {
          await ethSigner.provider.removeAllListeners('block');
          return resolve();
        }
        console.log(`  Confirmations: ${block - currentBlock}/${minSafetyBlocks}`)
      });
    });


    // wait for transfer submit event (tx from relay)
    console.log("Waiting for transfer submit event");
    await new Promise(resolve => ambBridge.once("TransferSubmit", resolve));

    // wait for transfer unlock event
    console.log("Waiting for transfer unlock event");
    await new Promise(resolve => ambBridge.once("TransferUnlock", resolve));

    // todo check event_id in events


    // check user balances
    const ethAfter = await ethToken.balanceOf(ethSigner.address);
    const ambAfter = await ambToken.balanceOf(ambSigner.address);

    expect(ethAfter).eq(ethBefore - 5);
    expect(ambAfter).eq(ambBefore + 5);

  });


});

// wait for transaction to be mined
async function w(call: Promise<any>): Promise<any> {
  return await (await call).wait();
}

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function getContractAddressYml() {
  const path = "../relay/configs/dev.yml";
  const raw = fs.readFileSync(path);
  const data = YAML.load(raw);

  return {
    ambAddress: data.network.amb["contract-addr"],
    ethAddress: data.network.eth["contract-addr"]
  };
}
