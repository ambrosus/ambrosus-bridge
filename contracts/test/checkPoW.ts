import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;


describe("Check PoW", () => {
  let ownerS: Signer;
  let owner: string;

  let ambBridge: Contract;


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);

    ambBridge = await ethers.getContract("CheckPoWTest", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state
  });


  it("Test CheckPoW", async function () {
    this.timeout(2 * 60 * 1000); // lol too long test

    const powProof = require("./fixtures/powProof.json");
    const epoch = require("./fixtures/epoch-406.json");

    await ambBridge.setEpochData(epoch.Epoch, epoch.FullSizeIn128Resolution, epoch.BranchDepth, epoch.MerkleNodes);
    await ambBridge.checkPoWTest(powProof, "0xf9427deDdAa899d388db70c0Fb4dA84A06976C85", {gasLimit: 40000000});
  });


  it("Test Ethash PoW", async () => {
    const blockPoW = require("./fixtures/BlockPoW-14257704.json");
    const epoch = require("./fixtures/epoch-475.json");

    await ambBridge.setEpochData(epoch.Epoch, epoch.FullSizeIn128Resolution, epoch.BranchDepth, epoch.MerkleNodes);
    expect(await ambBridge.isEpochDataSet(epoch.Epoch)).to.be.true;

    await ambBridge.verifyEthashTest(blockPoW);
  });

  // epoch-128 has 512 MerkleNodes
  it("Test submit epoch-128", async () => {
    const block = require("./fixtures/BlockPoW-3840001.json");
    const epoch = require("./fixtures/epoch-128.json");
    await ambBridge.setEpochData(epoch.Epoch, epoch.FullSizeIn128Resolution,
      epoch.BranchDepth, epoch.MerkleNodes, {gasLimit: 30000000});

    await ambBridge.verifyEthashTest(block);
  })

  it("Test setEpochData deleting old epochs", async () => {
    const epoch1 = require("./fixtures/epoch-475.json");
    const epoch2 = require("./fixtures/epoch-476.json");
    const epoch3 = require("./fixtures/epoch-477.json");
    await ambBridge.setEpochData(epoch1.Epoch, epoch1.FullSizeIn128Resolution,
        epoch1.BranchDepth, epoch1.MerkleNodes);
    await ambBridge.setEpochData(epoch2.Epoch, epoch2.FullSizeIn128Resolution,
        epoch2.BranchDepth, epoch2.MerkleNodes);
    expect(await ambBridge.isEpochDataSet(epoch1.Epoch)).to.be.true;
    await ambBridge.setEpochData(epoch3.Epoch, epoch3.FullSizeIn128Resolution,
        epoch3.BranchDepth, epoch3.MerkleNodes);
    expect(await ambBridge.isEpochDataSet(epoch1.Epoch)).to.be.false;
    expect(await ambBridge.isEpochDataSet(epoch2.Epoch)).to.be.true;
    expect(await ambBridge.isEpochDataSet(epoch3.Epoch)).to.be.true;
  });

  it("Test blockHash", async () => {
    const blockPoW = require("./fixtures/BlockPoW-14257704.json");
    const expectedBlockHash = "0xc4ca0efd5d528d67691abd9e10e9d4ca570f16235779e1f314b036caa5b455a1";

    const realBlockHash = await ambBridge.blockHashTest(blockPoW);
    expect(realBlockHash).eq(expectedBlockHash);
  });
});
