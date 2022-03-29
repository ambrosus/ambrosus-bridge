import chai from "chai";
import {ethers} from "hardhat";
import 'mocha';
import fs from "fs";
import {ambSigner, ethSigner, relayAddress} from "./cfg";

const YAML = require("js-yaml");


chai.should();


describe("Contract", () => {
  let ambNet;
  let ethNet;

  // it('Amb -> Eth', async () => {
  // });


  it('Eth -> Amb', async () => {
    const bridges = getContractAddressYml();
    const ambBridge = await ethers.getContractAt("AmbBridge", bridges.ambAddress, ambSigner);
    const ethBridge = await ethers.getContractAt("EthBridge", bridges.ethAddress, ethSigner);

    const mockTokens = require("../../mockTokensAddresses.json");
    const ambMockToken = await ethers.getContractAt("MockERC20", mockTokens.ambErc20Address, ambSigner);
    const ethMockToken = await ethers.getContractAt("MockERC20", mockTokens.ethErc20Address, ethSigner);


    // send funds

    // todo .connect(ethSigner) is redundant when contract deployed or got with ethSigned arg

    const hashAdmin = await ethBridge.callStatic.ADMIN_ROLE();

    await ethBridge.connect(ethSigner).grantRole(hashAdmin, ethSigner.address, {gasLimit: 100000});

    await ethSigner.sendTransaction({to: relayAddress, value: ethers.utils.parseEther("0.1")});


    await ethBridge.setSideBridge(ambBridge.address, {gasLimit: 100000});

    await ethMockToken.connect(ethSigner).mint(ethSigner.address, 30, {gasLimit: 100000});


    await ethMockToken.connect(ethSigner).increaseAllowance(ethBridge.address, 1000, {gasLimit: 8000000});
    await ethMockToken.increaseAllowance(ethSigner.address, 1000, {gasLimit: 8000000});
    await ethBridge.withdraw(ethMockToken.address, ambSigner.address, 10, {value: 1000, gasPrice: 1});

    console.log("wait for next timeframe...")
    await sleep(2000); // waiting for 100% new timeframe
    console.log("waiting is over")

    await ethBridge.callStatic.withdraw(ethMockToken.address, ambSigner.address, 3, {
      value: 1000,
      gasLimit: 8000000,
      gasPrice: 1
    });

    // event here
    console.log(await ambMockToken.connect(ambSigner).balanceOf(ambSigner.address));

  });


});


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
