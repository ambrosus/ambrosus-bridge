import {deployments, ethers, getNamedAccounts, network} from "hardhat";
import type {Contract, ContractReceipt, ContractTransaction, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

const ADMIN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const RELAY_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("RELAY_ROLE"));
const BRIDGE_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("BRIDGE_ROLE"));

const [token1, token2, token3, token4] = [
  "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002",
  "0x0000000000000000000000000000000000000003", "0x0000000000000000000000000000000000000004"];


describe("Common tests", () => {
  let ownerS: Signer;
  let relayS: Signer;
  let userS: Signer;
  let owner: string;
  let relay: string;
  let user: string;


  let commonBridge: Contract;
  let ethBridge: Contract;
  let ambBridge: Contract;
  let mockERC20: Contract;
  let wAmb: Contract;


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, relay, user} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);
    userS = await ethers.getSigner(user);

    commonBridge = await ethers.getContract("CommonBridge", ownerS);
    ambBridge = await ethers.getContract("AmbBridgeTest", ownerS);
    ethBridge = await ethers.getContract("EthBridgeTest", ownerS);
    mockERC20 = await ethers.getContract("BridgeERC20Test", ownerS);
    wAmb = await ethers.getContract("wAMB", ownerS);

  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state

    for (let bridge of [commonBridge, ethBridge, ambBridge]) {
      await bridge.grantRole(ADMIN_ROLE, owner);
      await bridge.grantRole(RELAY_ROLE, relay);
      await bridge.tokensAddBatch([wAmb.address, mockERC20.address], [token1, token2]);
      await mockERC20.grantRole(BRIDGE_ROLE, bridge.address);
    }
    await ambBridge.setAmbWrapper(wAmb.address);

    await mockERC20.mint(owner, 10000);
    await mockERC20.increaseAllowance(commonBridge.address, 5000);
  });


  describe("Test Withdraw", async () => {
    it("token balance changed", async () => {
      await expect(() => commonBridge.withdraw(mockERC20.address, user, 1, {value: 1000}))
        .to.changeTokenBalance(mockERC20, ownerS, -1);
    });

    it("withdraw event_id increased", async () => {
      await commonBridge.withdraw(mockERC20.address, user, 2, {value: 1000});
      await nextTimeframe();
      let tx1Amb: ContractTransaction = await commonBridge.withdraw(mockERC20.address, user, 1337, {value: 1000});
      await commonBridge.withdraw(mockERC20.address, user, 3, {value: 1000});
      await commonBridge.withdraw(mockERC20.address, user, 4, {value: 1000});
      await nextTimeframe();

      // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
      let tx2Amb: ContractTransaction = await commonBridge.withdraw(mockERC20.address, user, 1337, {value: 1000});
      await commonBridge.withdraw(mockERC20.address, user, 5, {value: 1000});

      let receipt1Amb: ContractReceipt = await tx1Amb.wait();
      let receipt2Amb: ContractReceipt = await tx2Amb.wait();

      // todo use truffle helpers for catch events

      let events1Amb: any = await getEvents(receipt1Amb);
      let events2Amb: any = await getEvents(receipt2Amb);

      // Checking that event_id increased
      expect(events2Amb[0].args.event_id).eq(events1Amb[0].args.event_id.add("1"));
    });
  });


  describe('Test wrap_withdraw in AmbBridge', async () => {

    it('Test wrap part', async () => {
      const fee = +await ambBridge.fee();

      await ambBridge.wrap_withdraw(user, {value: fee + 50});

      await expect(() => ambBridge.wrap_withdraw(user, {value: fee + 50}))
        .to.changeTokenBalance(wAmb, ambBridge, 50);
    });

    it('Test withdraw part', async () => {
      const fee = +await ambBridge.fee();

      await ambBridge.wrap_withdraw(user, {value: fee + 1});
      await ambBridge.wrap_withdraw(user, {value: fee + 1});
      await nextTimeframe();

      // will catch previous txs (because nextTimeframe happened)
      let tx1Amb: ContractTransaction = await ambBridge.wrap_withdraw(user, {value: fee + 1});
      await ambBridge.wrap_withdraw(user, {value: fee + 1});
      await ambBridge.wrap_withdraw(user, {value: fee + 1});
      await nextTimeframe();

      // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
      let tx2Amb: ContractTransaction = await ambBridge.wrap_withdraw(user, {value: fee + 1});
      await ambBridge.wrap_withdraw(user, {value: fee + 1});

      let receipt1Amb: ContractReceipt = await tx1Amb.wait();
      let receipt2Amb: ContractReceipt = await tx2Amb.wait();

      let events1Amb: any = await getEvents(receipt1Amb);
      let events2Amb: any = await getEvents(receipt2Amb);

      // Checking that event_id increased
      expect(events2Amb[0].args.event_id).eq(events1Amb[0].args.event_id.add("1"));
    });
  });


  describe("Token addresses", () => {
    it("add tokens", async () => {
      await commonBridge.tokensAdd(token1, token2);

      expect(await commonBridge.tokenAddresses(token1)).eq(token2);
    });

    it("add tokens batch", async () => {
      await commonBridge.tokensAddBatch([token1, token2], [token3, token4]);

      expect(await commonBridge.tokenAddresses(token1)).eq(token3);
      expect(await commonBridge.tokenAddresses(token2)).eq(token4);
    });

    it("remove tokens", async () => {
      await commonBridge.tokensAdd(token1, token2);
      await commonBridge.tokensRemove(token1);

      expect(await commonBridge.tokenAddresses(token1)).eq(ethers.constants.AddressZero);
    });

    it("remove tokens batch", async () => {
      await commonBridge.tokensAddBatch([token1, token2], [token3, token4]);
      await commonBridge.tokensRemoveBatch([token1, token2]);

      expect(await commonBridge.tokenAddresses(token1)).eq(ethers.constants.AddressZero);
      expect(await commonBridge.tokenAddresses(token2)).eq(ethers.constants.AddressZero);
    });
  });

  describe("Test change methods", () => {
    it("Test changeFeeRecipient", async () => {
      await commonBridge.changeFeeRecipient(user);
      await expect(() => commonBridge.withdraw(mockERC20.address, owner, 5, {value: 1000}))
        .to.changeEtherBalance(userS, 1000);
    });

    it("Test changeMinSafetyBlocks", async () => {
      await commonBridge.changeMinSafetyBlocks(20);
      expect(await commonBridge.minSafetyBlocks()).eq(20);
    });

    it("Test changeTimeframeSeconds", async () => {
      await commonBridge.changeTimeframeSeconds(20000);
      expect(await commonBridge.timeframeSeconds()).eq(20000);
    });

    it("Test changeLockTime", async () => {
      await commonBridge.changeLockTime(2000);
      expect(await commonBridge.lockTime()).eq(2000);
    });

    it("Test setSideBridge [AMB]", async () => {
      await ambBridge.setSideBridge(mockERC20.address);
      expect(await ambBridge.sideBridgeAddress()).eq(mockERC20.address);
    });
  });


  describe("Test Transfer lock/unlock/remove", () => {
    beforeEach(async () => {
      const data0 = [[mockERC20.address, user, 10]];
      const data1 = [[mockERC20.address, user, 20], [mockERC20.address, user, 30]];

      await ambBridge.lockTransfersTest(data0, 0);
      await ambBridge.lockTransfersTest(data1, 1);
    });

    it("locked correctly", async () => {
      const d0 = await ambBridge.getLockedTransferTest(0);
      const d1 = await ambBridge.getLockedTransferTest(1);

      expect(d0.transfers[0].amount).eq(10)
      expect(d1.transfers[0].amount).eq(20)
      expect(d1.transfers[1].amount).eq(30)
    });

    it("unlock", async () => {
      await nextTimeframe();

      await expect(() => ambBridge.unlockTransfers(0))
        .to.changeTokenBalance(mockERC20, userS, 10);
    });

    it("unlock before endTime passed", async () => {
      await expect(ambBridge.unlockTransfers(0))
        .to.be.revertedWith("lockTime has not yet passed")
    });

    it("unlock not oldest", async () => {
      await nextTimeframe();

      await expect(ambBridge.unlockTransfers(1))
        .to.be.revertedWith("can unlock only oldest event")
    });

    it("unlock batch", async () => {
      await nextTimeframe();

      await expect(() => ambBridge.unlockTransfersBatch())
        .to.changeTokenBalance(mockERC20, userS, 10 + 20 + 30);
    });

    it("remove transfers from 0", async () => {
      await ambBridge.pause();
      await ambBridge.removeLockedTransfers(0);

      expect((await ambBridge.getLockedTransferTest(0)).transfers).to.be.empty;
      expect((await ambBridge.getLockedTransferTest(1)).transfers).to.be.empty;
    });

    it("remove transfers from 1", async () => {
      await ambBridge.pause();
      await ambBridge.removeLockedTransfers(1);

      expect((await ambBridge.getLockedTransferTest(0)).transfers).to.not.be.empty;
      expect((await ambBridge.getLockedTransferTest(1)).transfers).to.be.empty;
    });

    it("remove unlocked", async () => {
      await nextTimeframe();
      await ambBridge.unlockTransfers(0);

      await ambBridge.pause();
      await expect(ambBridge.removeLockedTransfers(0)).to.be.revertedWith("event_id must be >= oldestLockedEventId");
    });
  });


  it('Test CalcTransferReceiptsHash', async () => {
    const receiptProof = require("./data-pow/receipt-proof-checkpow.json");
    const transferProof = [
      receiptProof, 1,
      [["0xc4b907fc242097D47eFd47f36eaee5Da2C239aDd", "0x8FC84c829d9cB1982f2121F135624E25aac679A9", 10]]
    ];
    const sideBridgeAddress = "0xd34baced0bf45ad4752783ad610450d0167ef6c7";

    expect(await ambBridge.CalcTransferReceiptsHash(transferProof, sideBridgeAddress))
      .to.eq("0x3cd6a7c9c4b79bd7231f9c85f7c6ef783b012faaadf908e54fb75c0b28ee2f88");
  });


});


let currentTimeframe = Math.floor(Date.now() / 14400);
const nextTimeframe = async (amount = 1) => {
  currentTimeframe += amount;
  const timestamp = currentTimeframe * 14400 + amount * 14400;
  await network.provider.send("evm_setNextBlockTimestamp", [timestamp]);
};

const getEvents = async (receipt: any) =>
  receipt.events?.filter((x: any) => x.event == "Transfer");
