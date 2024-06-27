import { deployments, ethers, getNamedAccounts } from "hardhat";
import type { Contract, Signer } from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;


describe("Check Bridge Networks Contracts", () => {
  let ownerS: Signer;
  let owner: string;
  let relayS: Signer;
  let relay: string;

  let bridgeAmbBsc: Contract;
  let bridgeBscAmb: Contract;
  let bridgeAmbEth: Contract;
  let bridgeEthAmb: Contract;

  const transfer = [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002",
    3,
  ];


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({ owner, relay } = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);

    bridgeAmbBsc = await ethers.getContract("BSC_AmbBridge", ownerS);
    bridgeBscAmb = await ethers.getContract("BSC_BscBridge", ownerS);
    bridgeAmbEth = await ethers.getContract("ETH_AmbBridge", ownerS);
    bridgeEthAmb = await ethers.getContract("ETH_EthBridge", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state
  });

  it("bsc->amb", async function () {
    await checkSubmitUntrustless(bridgeBscAmb);
  });
  it("amb->bsc", async function () {
    await checkSubmitUntrustless(bridgeAmbBsc);
    await checkSetSideBridge(bridgeAmbBsc)
  });

  it("eth->amb", async function () {
    await checkSubmitUntrustless(bridgeEthAmb);
  });
  it("amb->eth", async function () {
    await checkSubmitUntrustless(bridgeAmbEth);
    await checkSetSideBridge(bridgeAmbEth)
  });

  async function checkSubmitUntrustless(contract: Contract) {
    await expect(contract.submitTransferUntrustless(1, [transfer])).to.be.reverted;

    await contract.connect(relayS).submitTransferUntrustless(1, [transfer]);
    expect(await contract.inputEventId()).to.eq(1);
  }


  async function checkSetSideBridge(contract: Contract) {
    await contract.setSideBridge(bridgeAmbBsc.address);
    expect(await contract.sideBridgeAddress()).to.eq(bridgeAmbBsc.address);
    await expect(contract.setSideBridge(bridgeAmbBsc.address)).to.be.reverted;
  }


});

