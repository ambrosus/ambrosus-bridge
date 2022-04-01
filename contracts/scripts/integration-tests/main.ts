import chai from "chai";
import {ethers} from "hardhat";
import 'mocha';
import {ambSigner, ethSigner, loadYaml, relayAddress, w, options} from "./cfg";
import mockTokens from "./mockTokensAddresses.json";
import {Contract, providers} from "ethers";



chai.should();
const expect = chai.expect;


describe("Integration tests", function () {
  this.timeout(5 * 60 * 1000); // 5 minutes

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

    // setup relay
    console.log("Setup relay");
    const relayRole = await ethBridge.ADMIN_ROLE();
    await w(ethBridge.grantRole(relayRole, relayAddress, options));
    await w(ambBridge.grantRole(relayRole, relayAddress, options));

    await w(ethSigner.sendTransaction({to: relayAddress, value: ethers.utils.parseEther("0.1"), ...options}));
    await w(ambSigner.sendTransaction({to: relayAddress, value: ethers.utils.parseEther("1"), ...options}));


    // mint tokens
    console.log("Mint tokens");
    await w(ambToken.mint(ambSigner.address, 10, options));
    await w(ethToken.mint(ethSigner.address, 10, options));

    await w(ethToken.increaseAllowance(ethBridge.address, 10, options));
    await w(ambToken.increaseAllowance(ambBridge.address, 10, options));

  });

  // it('Amb -> Eth', async () => {
  // });


  it('Eth -> Amb', async () => {
    // check mint working
    const ethBefore = await ethToken.balanceOf(ethSigner.address);
    const ambBefore = await ambToken.balanceOf(ambSigner.address);

    expect(ethBefore).gte(10);
    expect(ambBefore).gte(10);


    // withdraw
    console.log("Withdraw");
    const fee = await ethBridge.fee();
    await w(ethBridge.withdraw(ethToken.address, ambSigner.address, 5, {value: fee, ...options}));
    await sleep(2000); // waiting for 100% new timeframe (timeframe is 1 second)
    await w(ethBridge.withdraw(ethToken.address, ambSigner.address, 1, {value: fee, ...options}));  // must emit event, todo check

    // wait for minSafetyBlocks confirmations
    console.log(`Wait for confirmations`);
    const minSafetyBlocks = await ambBridge.minSafetyBlocks();
    await waitConfirmations(minSafetyBlocks, ethSigner.provider);

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



async function waitConfirmations(minSafetyBlocks: number, provider: providers.Provider) {
  const currentBlock = (await ethSigner.provider.getBlock('latest')).number;
  await new Promise<void>((resolve) => {
    provider.on('block', async (block) => {
      if (block - currentBlock > minSafetyBlocks) {
        await provider.removeAllListeners('block');
        return resolve();
      }
      console.log(`  Confirmations: ${block - currentBlock}/${minSafetyBlocks}`)
    });
  });

}




function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function getContractAddressYml() {
  const yaml = loadYaml()
  return {
    ambAddress: yaml.network.amb["contract-addr"],
    ethAddress: yaml.network.eth["contract-addr"]
  };
}
