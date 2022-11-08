import {deployments, ethers, getNamedAccounts, network} from "hardhat";
import type {Contract, ContractReceipt, ContractTransaction, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

const ADMIN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const RELAY_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("RELAY_ROLE"));
const WATCHDOG_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("WATCHDOG_ROLE"));
const FEE_PROVIDER_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("FEE_PROVIDER_ROLE"));


const [token1, token2, token3, token4] = [
  "0x0000000000000000000000000000000000000001", "0x0000000000000000000000000000000000000002",
  "0x0000000000000000000000000000000000000003", "0x0000000000000000000000000000000000000004"];

const START_TIMESTAMP = 1652466039;

const transferFee = 111;
const bridgeFee = 222;

describe("Common tests", () => {
  let ownerS: Signer;
  let relayS: Signer;
  let userS: Signer;
  let adminS: Signer;
  let owner: string;
  let relay: string;
  let user: string;


  let commonBridge: Contract;
  let mockERC20: Contract;
  let sAmb: Contract;


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, relay, user} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);
    userS = await ethers.getSigner(user);

    commonBridge = await ethers.getContract("CommonBridgeTest", ownerS);
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

    await commonBridge.grantRole(ADMIN_ROLE, owner);
    await commonBridge.grantRole(WATCHDOG_ROLE, owner);
    await commonBridge.grantRole(RELAY_ROLE, relay);
    await commonBridge.grantRole(FEE_PROVIDER_ROLE, relay);
    await commonBridge.tokensAddBatch(Object.keys(tokens), Object.values(tokens));
    await mockERC20.setBridgeAddress(commonBridge.address);

    await mockERC20.mint(owner, 10000);
    await mockERC20.increaseAllowance(commonBridge.address, 5000);
  });

  // todo move to another test file?
  // describe("Test Proxy", async () => {
  //   it("ChangeAdmin check",async () => {
  //     await ambBridge.connect(proxyAdminS).changeAdmin(user);
  //     expect(await ambBridge.connect(userS).callStatic.admin()).eq(user);
  //   })
  // });

  describe("Test Withdraw", async () => {
    it("token balance changed", async () => {
      const signature = await getSignature(relayS, mockERC20.address, START_TIMESTAMP, 1);

      await expect(() =>
        commonBridge.withdraw(
          mockERC20.address,
          user, 1, false,
          signature, transferFee, bridgeFee,
          {value: transferFee + bridgeFee}))
        .to.changeTokenBalance(mockERC20, ownerS, -1);
    });

    it("withdraw eventId increased", async () => {
      let signature;
      let changedTimestamp;

      signature = await getSignature(relayS, mockERC20.address, START_TIMESTAMP, 1);
      await commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature));
      changedTimestamp = await nextTimeframe();

      signature = await getSignature(relayS, mockERC20.address, changedTimestamp, 1);
      let tx1Amb: ContractTransaction = await commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature));
      await commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature));
      await commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature));
      changedTimestamp = await nextTimeframe();

      // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
      signature = await getSignature(relayS, mockERC20.address, changedTimestamp, 1);
      let tx2Amb: ContractTransaction = await commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature));
      await commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature));

      let receipt1Amb: ContractReceipt = await tx1Amb.wait();
      let receipt2Amb: ContractReceipt = await tx2Amb.wait();

      // todo use truffle helpers for catch events

      let events1Amb: any = await getEvents(receipt1Amb);
      let events2Amb: any = await getEvents(receipt2Amb);

      // Checking that eventId increased
      expect(events2Amb[0].args.eventId).eq(events1Amb[0].args.eventId.add("1"));
    });

    it("unwrapSide == true (tokenTo should be 0x0)", async () => {
      const signature = await getSignature(relayS, mockERC20.address, START_TIMESTAMP, 1);
      await expect(commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature, true)))
        .to.emit(commonBridge, "Withdraw")
        .withNamedArgs({tokenTo: ethers.constants.AddressZero})
    });

    it("unwrapSide == true, but wrong token", async () => {
      const signature = await getSignature(relayS, token1, START_TIMESTAMP, 1);
      await expect(commonBridge.withdraw(...withdrawArgs(token1, user, signature, true)))
        .to.be.revertedWith("Token not point to native token")
    });

    it("wrong signature", async () => {
      const signature = await getSignature(relayS, mockERC20.address, START_TIMESTAMP, 1);
      await network.provider.send("evm_setNextBlockTimestamp", [START_TIMESTAMP + 4400]);
      await expect(commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature)))
          .to.be.revertedWith("Signature check failed");
    });

    it("withdraw feeCheck with delay", async () => {
      const signature = await getSignature(relayS, mockERC20.address, START_TIMESTAMP, 1);
      await network.provider.send("evm_setNextBlockTimestamp", [START_TIMESTAMP + 2400]);
      await commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature));
    });

    it("withdraw msg.value != transferFee + bridgeFee", async () => {
      const feeAddition = 60;
      const signature = await getSignature(relayS, mockERC20.address, START_TIMESTAMP, feeAddition);
      await expect(commonBridge.withdraw(...withdrawArgs(mockERC20.address, user, signature, false, feeAddition)))
          .to.be.revertedWith("Sent value != fee");
    });
  });


  describe('Test wrapWithdraw', async () => {

    it('Test wrap part', async () => {
      const fee = transferFee + bridgeFee;

      const wrapperAddress = await commonBridge.wrapperAddress();
      const feeAddition = 50;
      const signature = await getSignature(relayS, wrapperAddress, START_TIMESTAMP, feeAddition);
      await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + feeAddition});

      await expect(() => commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + feeAddition}))
        .to.changeTokenBalance(sAmb, commonBridge, 50);
    });

    it('Test withdraw part', async () => {
      let signature;
      let changedTimestamp = START_TIMESTAMP;
      const wrapperAddress = await commonBridge.wrapperAddress();

      const fee = transferFee + bridgeFee;

      signature = await getSignature(relayS, wrapperAddress, changedTimestamp, 1);
      await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + 1});
      await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + 1});
      changedTimestamp = await nextTimeframe();

      // will catch previous txs (because nextTimeframe happened)
      signature = await getSignature(relayS, wrapperAddress, changedTimestamp, 1);
      let tx1Amb: ContractTransaction = await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + 1});
      await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + 1});
      await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + 1});
      changedTimestamp = await nextTimeframe();

      // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
      signature = await getSignature(relayS, wrapperAddress, changedTimestamp, 1);
      let tx2Amb: ContractTransaction = await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + 1});
      await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee + 1});

      let receipt1Amb: ContractReceipt = await tx1Amb.wait();
      let receipt2Amb: ContractReceipt = await tx2Amb.wait();

      let events1Amb: any = await getEvents(receipt1Amb);
      let events2Amb: any = await getEvents(receipt2Amb);

      // Checking that eventId increased
      expect(events2Amb[0].args.eventId).eq(events1Amb[0].args.eventId.add("1"));
    });

    it('Check msg.value', async () => {
      const fee = transferFee + bridgeFee;
      const wrapperAddress = await commonBridge.wrapperAddress();
      const signature = await getSignature(relayS, wrapperAddress, START_TIMESTAMP, 0);

      await expect(commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: fee}))
          .to.be.revertedWith("Sent value <= fee");
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
      await commonBridge.changeTransferFeeRecipient(user);
      const signature = await getSignature(relayS, mockERC20.address, START_TIMESTAMP, 1);
      await expect(() => commonBridge.withdraw(...withdrawArgs(mockERC20.address, owner, signature)))
        .to.changeEtherBalance(userS, transferFee);

      await commonBridge.changeBridgeFeeRecipient(relay);
      await expect(() => commonBridge.withdraw(...withdrawArgs(mockERC20.address, owner, signature)))
          .to.changeEtherBalance(relayS, bridgeFee);
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

    it("Test ChangeSignatureFeeCheckNumber", async () => {
      await commonBridge.changeSignatureFeeCheckNumber(5);

      expect(await commonBridge.getSignatureFeeCheckNumber()).eq(5);
    });

    // todo move to another test file?
    // it("Test setSideBridge [AMB]", async () => {
    //   await ambBridge.setSideBridge(mockERC20.address);
    //   expect(await ambBridge.sideBridgeAddress()).eq(mockERC20.address);
    // });
  });


  describe("Test Transfer lock/unlock/remove/trigger", () => {
    beforeEach(async () => {
      const data1 = [[mockERC20.address, user, 10]];
      const data2 = [[mockERC20.address, user, 20], [mockERC20.address, user, 30]];

      await commonBridge.lockTransfersTest(data1, 1);
      await commonBridge.lockTransfersTest(data2, 2);
    });

    it("locked correctly", async () => {
      const d1 = await commonBridge.getLockedTransfers(1);
      const d2 = await commonBridge.getLockedTransfers(2);

      expect(d1.transfers[0].amount).eq(10)
      expect(d2.transfers[0].amount).eq(20)
      expect(d2.transfers[1].amount).eq(30)
    });

    it("unlock", async () => {
      await nextTimeframe();

      await expect(() => commonBridge.unlockTransfers(1))
        .to.changeTokenBalance(mockERC20, userS, 10);
    });

    it("unlock before endTime passed", async () => {
      await expect(commonBridge.unlockTransfers(1))
        .to.be.revertedWith("lockTime has not yet passed")
    });

    it("unlock not oldest", async () => {
      await nextTimeframe();

      await expect(commonBridge.unlockTransfers(2))
        .to.be.revertedWith("can unlock only oldest event")
    });

    it("unlock batch", async () => {
      await nextTimeframe();

      await expect(() => commonBridge.unlockTransfersBatch())
        .to.changeTokenBalance(mockERC20, userS, 10 + 20 + 30);
    });

    it("remove transfers from 1", async () => {
      await commonBridge.pause();
      await commonBridge.removeLockedTransfers(1);

      expect((await commonBridge.getLockedTransfers(1)).transfers).to.be.empty;
      expect((await commonBridge.getLockedTransfers(2)).transfers).to.be.empty;
      expect(await commonBridge.oldestLockedEventId()).to.be.eq(1);
      expect(await commonBridge.inputEventId()).to.be.eq(0);
    });

    it("remove transfers from 2", async () => {
      await commonBridge.pause();
      await commonBridge.removeLockedTransfers(2);

      expect((await commonBridge.getLockedTransfers(1)).transfers).to.not.be.empty;
      expect((await commonBridge.getLockedTransfers(2)).transfers).to.be.empty;
      expect(await commonBridge.oldestLockedEventId()).to.be.eq(1);
      expect(await commonBridge.inputEventId()).to.be.eq(1);
    });

    it("remove transfers from 2 when have unlocked transfer 1", async () => {
      await nextTimeframe();
      await commonBridge.unlockTransfers(1);

      await commonBridge.pause();
      await commonBridge.removeLockedTransfers(2);

      expect((await commonBridge.getLockedTransfers(1)).transfers).to.be.empty; // coz unlocked
      expect((await commonBridge.getLockedTransfers(2)).transfers).to.be.empty;
      expect(await commonBridge.oldestLockedEventId()).to.be.eq(2);
      expect(await commonBridge.inputEventId()).to.be.eq(1);
    });

    it("remove unlocked", async () => {
      await nextTimeframe();
      await commonBridge.unlockTransfers(1);

      await commonBridge.pause();
      await expect(commonBridge.removeLockedTransfers(1)).to.be.revertedWith("eventId must be >= oldestLockedEventId");
    });

    it("remove not locked", async () => {
      await commonBridge.pause();
      await expect(commonBridge.removeLockedTransfers(3)).to.be.revertedWith("eventId must be <= inputEventId");
    });

    it("trigger transfers event check", async () => {
      const beforeEventOutputEventId = await commonBridge.outputEventId();
      await commonBridge.addElementToQueue();

      const tx = await commonBridge.triggerTransfers();
      const receipt = await tx.wait();
      const events = await getEvents(receipt);

      const afterEventOutputEventId = await commonBridge.outputEventId();

      expect(events[0].event).eq("Transfer");
      expect(beforeEventOutputEventId.add("0x1")).eq(afterEventOutputEventId);
    });

    it("trigger transfers empty queue check", async () => {
      await expect(commonBridge.triggerTransfers())
          .to.be.revertedWith("Queue is empty");
    });


    it("skip transfers", async () => {
      await commonBridge.pause();
      await commonBridge.skipTransfers(4);

      expect(await commonBridge.oldestLockedEventId()).to.be.eq(4);
      expect(await commonBridge.inputEventId()).to.be.eq(3);
    });



  });
  it("unlock native coins", async () => { // separate from previous describe block coz of hindering `beforeEach`
    const wrapperAddress = await commonBridge.wrapperAddress();

    // lock some SAMB tokens on bridge
    const signature = await getSignature(relayS, wrapperAddress, START_TIMESTAMP, 50);
    await commonBridge.wrapWithdraw(user, signature, transferFee, bridgeFee, {value: transferFee + bridgeFee + 50});

    await commonBridge.lockTransfersTest([[ethers.constants.AddressZero, user, 25]], 1);
    await nextTimeframe();

    await expect(() => commonBridge.unlockTransfers(1))
      .to.changeEtherBalance(userS, 25);
  });



  it('Test calcTransferReceiptsHash', async () => {
    const receiptProof = require("./fixtures/transfer-event-proof.json");
    const transferProof = [
      receiptProof, 1,
      [["0xc4b907fc242097D47eFd47f36eaee5Da2C239aDd", "0x8FC84c829d9cB1982f2121F135624E25aac679A9", 10]]
    ];
    const sideBridgeAddress = "0xd34baced0bf45ad4752783ad610450d0167ef6c7";

    expect(await commonBridge.calcTransferReceiptsHashTest(transferProof, sideBridgeAddress))
      .to.eq("0x3cd6a7c9c4b79bd7231f9c85f7c6ef783b012faaadf908e54fb75c0b28ee2f88");
  });


  it('Test checkSignature', async () => {
    const hash = "0x1d0a6ca42217dc9f0560840b3eb91a3879b836cb7ec5a8055e265a520e6839d0";
    const signature = "0x5c1974f609035dc81319f058a8b9428b7ce26b366fadf9768b8ca19e3014c759467d732731a58a2ad9f3e9efedc56275427cd4a2fd7a6de59007b0bdb2e95f7d00";
    const needAddress = "0xc89C669357D161d57B0b255C94eA96E179999919";
    expect(ethers.utils.recoverAddress(ethers.utils.arrayify(hash), signature)).eq(needAddress);

    expect(await commonBridge.checkSignatureTest(hash, signature)).eq(needAddress);
  });
});


let currentTimeframe = Math.floor(START_TIMESTAMP / 14400);
const nextTimeframe = async (amount = 1) => {
  currentTimeframe += amount;
  const timestamp = currentTimeframe * 14400 + amount * 14400;
  await network.provider.send("evm_setNextBlockTimestamp", [timestamp]);

  return timestamp;
};

const getEvents = async (receipt: any) =>
  receipt.events?.filter((x: any) => x.event == "Transfer");

const withdrawArgs = (token: string, user: string, signature: any, unwrapSide = false, feeAddition = 0) => {
  return [token, user, 1, unwrapSide, signature, transferFee, bridgeFee, {value: +transferFee + bridgeFee + feeAddition}];
}

const getSignature = async (
    signer: Signer,
    tokenAddress: string,
    timestamp: number,
    amount: number
) => {
  const msg =  ethers.utils.arrayify(
    ethers.utils.solidityKeccak256(
      ["address", "uint", "uint", "uint", "uint"],
      [
        tokenAddress,
        ethers.utils.hexlify(Math.floor(timestamp / 1800)),
        ethers.utils.hexlify(transferFee),
        ethers.utils.hexlify(bridgeFee),
        ethers.utils.hexlify(amount),
      ]
    ))

  return await signer.signMessage(msg);
};
