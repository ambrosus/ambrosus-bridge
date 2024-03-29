import chai from "chai";
import {ethers} from "hardhat";
import 'mocha';
import {Contract, providers} from "ethers";
import config_c from "../configs/integr.json";
import config_r from "../../relay/configs/integr.json";

// websocket need for events subscribe
const ambNet = new ethers.providers.WebSocketProvider(config_r.network.amb.wsUrl);
const ethNet = new ethers.providers.WebSocketProvider(config_r.network.eth.wsUrl);

// accounts with money; details in /tests/README.md; also set in hardhat config
const ambSigner = new ethers.Wallet("80f702eb861f36fe8fbbe1a7ccceb04ef7ddef714604010501a5f67c8065d446", ambNet);
const ethSigner = new ethers.Wallet("51d098d8aee092622149d8f3a79cc7b1ce36ff97fadaa2fbd623c65badeefadc", ethNet);

// must be equal to pk set in relay env var
const relayAddress = new ethers.Wallet("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5").address;
const vsContractAddress = config_r.network.amb.vsContractAddr;

const validators = [
  "0x4c9785451bb2CA3E91B350C06bcB5f974cA33F79",
  "0x90B2Ce3741188bCFCe25822113e93983ecacfcA0",
  "0xAccdb7a2268BC4Af0a1898e725138888ba1Ca6Fc"
];

const options = {gasLimit: 8_000_000};  // amb & eth gas estimator broken


chai.should();
const expect = chai.expect;


describe("Integration tests", function () {
  this.timeout(5 * 60 * 1000); // 5 minutes

  let ambBridge: Contract;
  let ethBridge: Contract;
  let ambToken: Contract;
  let ethToken: Contract;
  let vs: Contract;

  before(async () => {
    ambBridge = await ethers.getContractAt("AmbBridge", config_c.bridges.eth.amb, ambSigner);
    ethBridge = await ethers.getContractAt("EthBridge", config_c.bridges.eth.side, ethSigner);
    ambToken = await ethers.getContractAt("sAMB", config_c.tokens.SAMB.addresses.amb, ambSigner);
    ethToken = await ethers.getContractAt("BridgeERC20Test", config_c.tokens.SAMB.addresses.eth, ethSigner);
    vs = await ethers.getContractAt("ModifiedValidatorSet", vsContractAddress, ambSigner);

    console.log("Setup relay")
    // set relay role
    const relayRole = await ethBridge.RELAY_ROLE();
    await ws(
      ethBridge.grantRole(relayRole, relayAddress, options),
      ambBridge.grantRole(relayRole, relayAddress, options),

      // send money to relay
      ethSigner.sendTransaction({to: relayAddress, value: ethers.utils.parseEther("0.1"), ...options}),
      ambSigner.sendTransaction({to: relayAddress, value: ethers.utils.parseEther("1"), ...options})
    );

  });


  it('Eth -> Amb', async () => {
    console.log("ETH");
    const ambBefore = await ambToken.balanceOf(ambSigner.address);

    // mint tokens
    console.log("Mint tokens");
    await ws(
      ethToken.mint(ethSigner.address, 10, options),
      ethToken.increaseAllowance(ethBridge.address, 10, options));

    // withdraw
    console.log("Withdraw");
    const fee = await ethBridge.fee();
    await w(ethBridge.withdraw(ethToken.address, ambSigner.address, 5, false, {value: fee, ...options}));  // must emit event, todo check

    // wait for minSafetyBlocks confirmations
    console.log(`Wait for confirmations`);
    const minSafetyBlocks = await ambBridge.minSafetyBlocks();
    await waitConfirmations(minSafetyBlocks, ethSigner.provider);

    await waitForFinish(ethBridge)
    // todo check event_id in events


    // check user balance
    expect(await ambToken.balanceOf(ambSigner.address)).eq(ambBefore + 5);

  });

  it('Amb -> Eth', async () => {
    console.log("AMB");
    const ethBefore = await ethToken.balanceOf(ethSigner.address);

    // validator set changes
    console.log("Validator set changes");

    // return all validators to set if they were removed
    const currentSet = await vs.getValidators();
    await ws(...
      validators
        .filter(v => !currentSet.includes(v))
        .map(v => vs.addValidator(v, options))
    )

    // todo this doesn't work, seems we have to use a contract for this, or add validators without nodes
    // in one block
    // await ws(
    //   vs.removeValidator(validators[1]),
    //   vs.addValidator(validators[1]))

    // one at time
    await w(vs.removeValidator(validators[2]))
    await w(vs.addValidator(validators[2]))


    // withdraw
    console.log("Withdraw");
    const fee = await ambBridge.fee();
    await w(ambBridge.wrapWithdraw(ambSigner.address, {value: +fee + 5, ...options}));


    // wait for minSafetyBlocks confirmations
    console.log(`Wait for confirmations`);
    const minSafetyBlocks = await ethBridge.minSafetyBlocks();
    await waitConfirmations(minSafetyBlocks, ambSigner.provider);

    await waitForFinish(ambBridge)

    // todo check event_id in events


    // check user balances
    expect(await ethToken.balanceOf(ethSigner.address)).eq(ethBefore + 5);
  });


});


async function waitConfirmations(minSafetyBlocks: number, provider: providers.Provider) {
  const currentBlock = (await provider.getBlock('latest')).number;
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

async function waitForFinish(bridge: Contract) {
  // wait for transfer submit event (tx from relay)
  console.log("Waiting for transfer submit event");
  await new Promise(resolve => bridge.once("TransferSubmit", resolve));

  // wait for transfer unlock event
  console.log("Waiting for transfer unlock event");
  await new Promise(resolve => bridge.once("TransferUnlock", resolve));
}


const getEvents = (receipt: any, eventName: string) => {
  return receipt.events?.filter((x: any) => x.event === eventName);
};

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// wait for transaction to be mined
async function w(call: Promise<any>): Promise<any> {
  const tx = await (await call).wait();
  console.log('Transaction mined');
  return tx;
}

// wait for transactions to be mined
async function ws(...calls: Promise<any>[]): Promise<any> {
  return Promise.all(calls.map(c => w(c)));
}
