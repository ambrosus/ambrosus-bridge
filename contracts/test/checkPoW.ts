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
  let ambBridgeTest: Contract;


  before(async () => {
    await deployments.fixture(["ambbridge", "ambbridgetest"]);
    ({owner, relay} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);

    ambBridge = await ethers.getContract("AmbBridge", ownerS);
    ambBridgeTest = await ethers.getContract("AmbBridgeTest", ownerS);

    await ambBridge.grantRole(ADMIN_ROLE, owner);
    await ambBridgeTest.grantRole(ADMIN_ROLE, owner);

    await ambBridge.grantRole(RELAY_ROLE, relay);
    await ambBridgeTest.grantRole(RELAY_ROLE, relay);
  });

  beforeEach(async () => {
    await deployments.fixture(["ambbridge", "ambbridgetest"]); // reset contracts state
  });

  it("Test Ethash PoW", async () => {
    const blockPoW = require("../../relay/cmd/dump-test-data/BlockPoW-14257704.json");
    const epoch = require("../../relay/cmd/dump-test-data/epoch-475.json");

    await submitEpochData(ambBridge, epoch);
    expect(await ambBridge.isEpochDataSet(epoch.Epoch)).to.be.true;

    await ambBridge.verifyEthash(blockPoW);
  });

  it("TEST epochdata new", async () => {
    const blockDelete = require("../../relay/cmd/dump-test-data/BlockPoW-228.json");
    const epoch = require("../../relay/assets/testdata/epoch-128.json");
    await ambBridge.setEpochDataTest(epoch.Epoch, epoch.FullSizeIn128Resolution,
      epoch.BranchDepth, epoch.MerkleNodes, {gasLimit: 30000000});

    await ambBridge.verifyEthash(blockDelete);
  })

  it("Test setEpochData deleting old epochs", async () => {
    const epoch1 = require("../../relay/assets/testdata/epoch-475.json");
    const epoch2 = require("../../relay/assets/testdata/epoch-476.json");
    const epoch3 = require("../../relay/assets/testdata/epoch-477.json");

    await submitEpochData(ambBridge, epoch1);
    await submitEpochData(ambBridge, epoch2);
    expect(await ambBridge.isEpochDataSet(epoch1.Epoch)).to.be.true;
    await submitEpochData(ambBridge, epoch3);
    expect(await ambBridge.isEpochDataSet(epoch1.Epoch)).to.be.false;
    expect(await ambBridge.isEpochDataSet(epoch2.Epoch)).to.be.true;
    expect(await ambBridge.isEpochDataSet(epoch3.Epoch)).to.be.true;
  });

  it("Test blockHash", async () => {
    const blockPoW = require("../../relay/cmd/dump-test-data/BlockPoW-14257704.json");
    const expectedBlockHash = "0xc4ca0efd5d528d67691abd9e10e9d4ca570f16235779e1f314b036caa5b455a1";

    const realBlockHash = await ambBridgeTest.blockHashTest(blockPoW);
    expect(realBlockHash).eq(expectedBlockHash);
  });


  const submitEpochData = async (ethashContractInstance: Contract, epoch: any) => {
    let start = 0;
    let nodes: any = [];
    let mnlen = 0;
    let index = 0;

    for (let mn of epoch.MerkleNodes) {
      nodes.push(mn);

      if (nodes.length === 40 || index === epoch.MerkleNodes.length - 1) {
        mnlen = nodes.length;
        if (index < 440 && epoch.Number === 128) {
          start = start + mnlen;
          nodes = [];
          return;
        }
        await ethashContractInstance.setEpochData(
          epoch.Epoch, epoch.FullSizeIn128Resolution, epoch.BranchDepth,
          nodes, start, mnlen
        );

        start = start + mnlen;
        nodes = [];
      }
      index++;
    }
  };

});
