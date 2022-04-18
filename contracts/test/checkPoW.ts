import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

const ADMIN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const RELAY_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("RELAY_ROLE"));

describe("Check PoW", () => {
  let ownerS: Signer;
  let relayS: Signer;
  let owner: string;
  let relay: string;

  let ambBridge: Contract;


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, relay} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);

    ambBridge = await ethers.getContract("AmbBridgeTest", ownerS);

    await ambBridge.grantRole(ADMIN_ROLE, owner);
    await ambBridge.grantRole(RELAY_ROLE, relay);
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state
  });

  it("Test Ethash PoW", async () => {
    const blockPoW = require("./fixtures/BlockPoW-14257704.json");
    const epoch = require("./fixtures/epoch-475.json");

    await ambBridge.setEpochData(epoch.Epoch, epoch.FullSizeIn128Resolution,
        epoch.BranchDepth, epoch.MerkleNodes);
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
