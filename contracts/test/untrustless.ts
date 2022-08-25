import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;


describe("Check Untrustless", () => {
  let ownerS: Signer;
  let owner: string;
  let relayS: Signer;
  let relay: string;

  let bridge: Contract;

  const transfer = [
    "0x0000000000000000000000000000000000000001",
    "0x0000000000000000000000000000000000000002",
    3,
  ];


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, relay} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);

    bridge = await ethers.getContract("CheckUntrustlessTest", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state
  });

  describe("setRelaysAndConfirmationsTest", () => {

    it("relays", async function () {
      expect(await bridge.isRelay(owner)).to.be.false;
      expect(await bridge.isRelay(relay)).to.be.false;

      await bridge.setRelaysAndConfirmationsTest([], [owner], 2);
      expect(await bridge.isRelay(owner)).to.be.true;

      await bridge.setRelaysAndConfirmationsTest([owner], [relay], 2);
      expect(await bridge.isRelay(owner)).to.be.false;
      expect(await bridge.isRelay(relay)).to.be.true;
    });

    it("confirmations", async function () {
      await bridge.setRelaysAndConfirmationsTest([], [], 2);
      expect(await bridge.confirmationsThreshold()).to.eq(2);

      await bridge.setRelaysAndConfirmationsTest([], [], 3);
      expect(await bridge.confirmationsThreshold()).to.eq(3);
    });

    it("RelayAdd event", async function () {
      await expect(bridge.setRelaysAndConfirmationsTest([], [owner], 2))
        .to.emit(bridge, "RelayAdd").withArgs(owner);
    });

    it("RelayAdd already relay", async function () {
      bridge.setRelaysAndConfirmationsTest([], [owner], 1);

      await expect(bridge.setRelaysAndConfirmationsTest([], [owner], 1))
        .to.be.revertedWith("Already relay");
    });


    it("RelayRemove event", async function () {
      bridge.setRelaysAndConfirmationsTest([], [owner], 1);

      await expect(bridge.setRelaysAndConfirmationsTest([owner], [], 1))
        .to.emit(bridge, "RelayRemove").withArgs(owner);
    });

    it("RelayRemove not relay", async function () {
      bridge.setRelaysAndConfirmationsTest([], [relay], 1);

      await expect(bridge.setRelaysAndConfirmationsTest([owner], [], 1))
        .to.be.revertedWith("Not a relay");
    });

    it("Change threshold event", async function () {
      bridge.setRelaysAndConfirmationsTest([], [], 1);

      await expect(bridge.setRelaysAndConfirmationsTest([], [], 2))
        .to.emit(bridge, "ThresholdChange").withArgs(2);
    });

    it("Change threshold no event", async function () {
      bridge.setRelaysAndConfirmationsTest([], [], 1);

      await expect(bridge.setRelaysAndConfirmationsTest([], [], 1))
        .to.not.emit(bridge, "ThresholdChange");
    });

  });

  describe("checkUntrustless", () => {

    it("Not in whitelist", async function () {
      await expect(bridge.checkUntrustlessTest(1, [transfer]))
        .to.be.revertedWith("You not in relay whitelist");
    });


    it("Double confirm", async function () {
      await bridge.setRelaysAndConfirmationsTest([], [owner], 2);

      await bridge.checkUntrustlessTest(1, [transfer]);

      expect(await bridge.isConfirmedByRelay(owner, 1, [transfer])).to.be.true;

      await expect(bridge.checkUntrustlessTest(1, [transfer]))
        .to.be.revertedWith("You have already confirmed this");
    });


    it("Confirm confirmed", async function () {
      await bridge.setRelaysAndConfirmationsTest([], [owner, relay], 1);

      expect(await getResult(bridge.checkUntrustlessTest(1, [transfer]))).to.be.true; // confirmed!

      await expect(bridge.connect(relayS).checkUntrustlessTest(1, [transfer]))
        .to.be.revertedWith("Already confirmed");
    });


    it("Different hash", async function () {
      await bridge.setRelaysAndConfirmationsTest([], [owner], 2);

      expect(await getResult(bridge.checkUntrustlessTest(1, [transfer]))).to.be.false; // not enough confirmations

      expect(await getResult(bridge.checkUntrustlessTest(2, [transfer]))).to.be.false; // different hash
      expect(await getResult(bridge.checkUntrustlessTest(1, [transfer, transfer]))).to.be.false; // different hash again
    });

    it("Successfully confirm", async function () {
      await bridge.setRelaysAndConfirmationsTest([], [owner, relay], 2);

      expect(await getResult(bridge.checkUntrustlessTest(1, [transfer]))).to.be.false; // not enough confirmations
      expect(await getResult(bridge.connect(relayS).checkUntrustlessTest(1, [transfer]))).to.be.true; // confirmed!
    });

    // get result from event, emitted by test wrapper (coz can't get return value from non view function)
    async function getResult(func: any) {
      const res = await (await func).wait();
      return res.events[1].args.result;
    }
  });


});
