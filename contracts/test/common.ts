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
  let proxyAdminS: Signer;
  let owner: string;
  let relay: string;
  let user: string;
  let proxyAdmin: string;


  let commonBridge: Contract;
  let ethBridge: Contract;
  let ambBridge: Contract;
  let mockERC20: Contract;
  let sAmb: Contract;


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, relay, user, proxyAdmin} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);
    userS = await ethers.getSigner(user);
    proxyAdminS = await ethers.getSigner(proxyAdmin);

    commonBridge = await ethers.getContract("CommonBridgeTest", ownerS);
    ambBridge = await ethers.getContract("AmbBridgeTest", ownerS);
    ethBridge = await ethers.getContract("EthBridgeTest", ownerS);
    mockERC20 = await ethers.getContract("BridgeERC20Test", ownerS);
    sAmb = await ethers.getContract("sAMB", ownerS);

  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state

    const tokens = {
      [sAmb.address]: token1,
      [mockERC20.address]: token2,
      [ethers.constants.AddressZero]: mockERC20.address,  // mean that mockERC20 point to native token in side network
    }

    for (let bridge of [commonBridge, ethBridge, ambBridge]) {
      await bridge.grantRole(ADMIN_ROLE, owner);
      await bridge.grantRole(RELAY_ROLE, relay);
      await bridge.tokensAddBatch(Object.keys(tokens), Object.values(tokens));
      await mockERC20.grantRole(BRIDGE_ROLE, bridge.address);
    }

    await mockERC20.mint(owner, 10000);
    await mockERC20.increaseAllowance(commonBridge.address, 5000);
  });

  describe("Test Proxy", async () => {
    it("ChangeAdmin check",async () => {
      await ambBridge.connect(proxyAdminS).changeAdmin(user);
      expect(await ambBridge.connect(userS).callStatic.admin()).eq(user);
    })
  });

  describe("Test Withdraw", async () => {
    it("token balance changed", async () => {
      await expect(() => commonBridge.withdraw(mockERC20.address, user, 1, false, {value: 1000}))
        .to.changeTokenBalance(mockERC20, ownerS, -1);
    });

    it("withdraw eventId increased", async () => {
      await commonBridge.withdraw(mockERC20.address, user, 2, false, {value: 1000});
      await nextTimeframe();
      let tx1Amb: ContractTransaction = await commonBridge.withdraw(mockERC20.address, user, 1337, false, {value: 1000});
      await commonBridge.withdraw(mockERC20.address, user, 3, false, {value: 1000});
      await commonBridge.withdraw(mockERC20.address, user, 4, false, {value: 1000});
      await nextTimeframe();

      // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
      let tx2Amb: ContractTransaction = await commonBridge.withdraw(mockERC20.address, user, 1337, false, {value: 1000});
      await commonBridge.withdraw(mockERC20.address, user, 5, false, {value: 1000});

      let receipt1Amb: ContractReceipt = await tx1Amb.wait();
      let receipt2Amb: ContractReceipt = await tx2Amb.wait();

      // todo use truffle helpers for catch events

      let events1Amb: any = await getEvents(receipt1Amb);
      let events2Amb: any = await getEvents(receipt2Amb);

      // Checking that eventId increased
      expect(events2Amb[0].args.eventId).eq(events1Amb[0].args.eventId.add("1"));
    });

    it("unwrapSide == true", async () => {
      const tx = await commonBridge.withdraw(mockERC20.address, user, 1, true, {value: 1000})
      // todo check address in event
    });

    it("unwrapSide == true, but wrong token", async () => {
      await expect(commonBridge.withdraw(token1, user, 1, true, {value: 1000}))
        .to.be.revertedWith("Token not point to native token")
    });

  });


  describe('Test wrapWithdraw', async () => {

    it('Test wrap part', async () => {
      const fee = +await commonBridge.fee();

      await commonBridge.wrapWithdraw(user, {value: fee + 50});

      await expect(() => commonBridge.wrapWithdraw(user, {value: fee + 50}))
        .to.changeTokenBalance(sAmb, commonBridge, 50);
    });

    it('Test withdraw part', async () => {
      const fee = +await commonBridge.fee();

      await commonBridge.wrapWithdraw(user, {value: fee + 1});
      await commonBridge.wrapWithdraw(user, {value: fee + 1});
      await nextTimeframe();

      // will catch previous txs (because nextTimeframe happened)
      let tx1Amb: ContractTransaction = await commonBridge.wrapWithdraw(user, {value: fee + 1});
      await commonBridge.wrapWithdraw(user, {value: fee + 1});
      await commonBridge.wrapWithdraw(user, {value: fee + 1});
      await nextTimeframe();

      // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
      let tx2Amb: ContractTransaction = await commonBridge.wrapWithdraw(user, {value: fee + 1});
      await commonBridge.wrapWithdraw(user, {value: fee + 1});

      let receipt1Amb: ContractReceipt = await tx1Amb.wait();
      let receipt2Amb: ContractReceipt = await tx2Amb.wait();

      let events1Amb: any = await getEvents(receipt1Amb);
      let events2Amb: any = await getEvents(receipt2Amb);

      // Checking that eventId increased
      expect(events2Amb[0].args.eventId).eq(events1Amb[0].args.eventId.add("1"));
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
      await expect(() => commonBridge.withdraw(mockERC20.address, owner, 5, false, {value: 1000}))
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
      const data1 = [[mockERC20.address, user, 10]];
      const data2 = [[mockERC20.address, user, 20], [mockERC20.address, user, 30]];

      await ambBridge.lockTransfersTest(data1, 1);
      await ambBridge.lockTransfersTest(data2, 2);
    });

    it("locked correctly", async () => {
      const d1 = await ambBridge.getLockedTransferTest(1);
      const d2 = await ambBridge.getLockedTransferTest(2);

      expect(d1.transfers[0].amount).eq(10)
      expect(d2.transfers[0].amount).eq(20)
      expect(d2.transfers[1].amount).eq(30)
    });

    it("unlock", async () => {
      await nextTimeframe();

      await expect(() => ambBridge.unlockTransfers(1))
        .to.changeTokenBalance(mockERC20, userS, 10);
    });

    it("unlock before endTime passed", async () => {
      await expect(ambBridge.unlockTransfers(1))
        .to.be.revertedWith("lockTime has not yet passed")
    });

    it("unlock not oldest", async () => {
      await nextTimeframe();

      await expect(ambBridge.unlockTransfers(2))
        .to.be.revertedWith("can unlock only oldest event")
    });

    it("unlock batch", async () => {
      await nextTimeframe();

      await expect(() => ambBridge.unlockTransfersBatch())
        .to.changeTokenBalance(mockERC20, userS, 10 + 20 + 30);
    });

    it("remove transfers from 1", async () => {
      await ambBridge.pause();
      await ambBridge.removeLockedTransfers(1);

      expect((await ambBridge.getLockedTransferTest(1)).transfers).to.be.empty;
      expect((await ambBridge.getLockedTransferTest(2)).transfers).to.be.empty;
    });

    it("remove transfers from 2", async () => {
      await ambBridge.pause();
      await ambBridge.removeLockedTransfers(2);

      expect((await ambBridge.getLockedTransferTest(1)).transfers).to.not.be.empty;
      expect((await ambBridge.getLockedTransferTest(2)).transfers).to.be.empty;
    });

    it("remove unlocked", async () => {
      await nextTimeframe();
      await ambBridge.unlockTransfers(1);

      await ambBridge.pause();
      await expect(ambBridge.removeLockedTransfers(1)).to.be.revertedWith("eventId must be >= oldestLockedEventId");
    });

    it("unlock native coins", async () => {
      await ambBridge.wrapWithdraw(user, {value: +await commonBridge.fee() + 50});  // lock some SAMB tokens on bridge
      await ambBridge.lockTransfersTest([[ethers.constants.AddressZero, user, 25]], 1);
      await nextTimeframe();

      await expect(() => ambBridge.unlockTransfers(1))
        .to.changeEtherBalance(userS, 25);
    });
  });


  it('Test calcTransferReceiptsHash', async () => {
    const receiptProof = require("./fixtures/transfer-event-proof.json");
    const transferProof = [
      receiptProof, 1,
      [["0xc4b907fc242097D47eFd47f36eaee5Da2C239aDd", "0x8FC84c829d9cB1982f2121F135624E25aac679A9", 10]]
    ];
    const sideBridgeAddress = "0xd34baced0bf45ad4752783ad610450d0167ef6c7";

    expect(await ambBridge.calcTransferReceiptsHashTest(transferProof, sideBridgeAddress))
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
