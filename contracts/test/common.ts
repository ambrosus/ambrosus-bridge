import {deployments, ethers, getNamedAccounts, network} from "hardhat";
import type {Contract, ContractReceipt, ContractTransaction, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

const ADMIN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const RELAY_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("RELAY_ROLE"));
const BRIDGE_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("BRIDGE_ROLE"));

describe("Common tests", () => {
  let ownerS: Signer;
  let relayS: Signer;
  let userS: Signer;
  let owner: string;
  let relay: string;
  let user: string;


  let ethBridge: Contract;
  let ambBridge: Contract;
  let mockERC20: Contract;
  let ambBridgeTest: Contract;
  let ethBridgeTest: Contract;
  let wAmb: Contract;


  before(async () => {
    await deployments.fixture(["ethbridge", "ambbridge", "mocktoken", "ambbridgetest", "wamb"]);
    ({owner, relay, user} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);
    userS = await ethers.getSigner(user);

    ethBridge = await ethers.getContract("EthBridge", ownerS);
    ambBridge = await ethers.getContract("AmbBridge", ownerS);
    mockERC20 = await ethers.getContract("MockERC20", ownerS);
    ambBridgeTest = await ethers.getContract("AmbBridgeTest", ownerS);
    ethBridgeTest = await ethers.getContract("EthBridgeTest", ownerS);
    wAmb = await ethers.getContract("wAMB", ownerS);

    await ambBridge.grantRole(ADMIN_ROLE, owner);
    await ambBridgeTest.grantRole(ADMIN_ROLE, owner);
    await ethBridge.grantRole(ADMIN_ROLE, owner);

    await ambBridge.grantRole(RELAY_ROLE, relay);
    await ambBridgeTest.grantRole(RELAY_ROLE, relay);
    await ethBridge.grantRole(RELAY_ROLE, relay);

    await mockERC20.grantRole(BRIDGE_ROLE, ambBridge.address);
    await mockERC20.grantRole(BRIDGE_ROLE, ethBridge.address);
  });

  beforeEach(async () => {
    await deployments.fixture(["ethbridge", "ambbridge", "ambbridgetest"]); // reset contracts state
  });

  it("Really common", async () => {
    // todo if this shit really common why the fuck would you test it for all contracts? :thinking:
    const bridges = [
      {name: "ETH", bridge: ethBridge},
      {name: "AMB", bridge: ambBridge},
    ];
    for (const {name, bridge} of bridges) describe("Bridge " + name, async () => {


      it("TestWithdraw timeframe", async () => {
        await mockERC20.mint(owner, 10000);
        await mockERC20.increaseAllowance(bridge.address, 5000);

        await bridge.withdraw(mockERC20.address, user, 1, {value: 1000});
        await bridge.withdraw(mockERC20.address, user, 2, {value: 1000});
        await nextTimeframe();

        // will catch previous txs (because nextTimeframe happened)
        let tx1Amb: ContractTransaction = await bridge.withdraw(mockERC20.address, user, 1337, {value: 1000});
        await bridge.withdraw(mockERC20.address, user, 3, {value: 1000});
        await bridge.withdraw(mockERC20.address, user, 4, {value: 1000});
        await nextTimeframe();

        // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
        let tx2Amb: ContractTransaction = await bridge.withdraw(mockERC20.address, user, 1337, {value: 1000});
        await bridge.withdraw(mockERC20.address, user, 5, {value: 1000});

        let receipt1Amb: ContractReceipt = await tx1Amb.wait();
        let receipt2Amb: ContractReceipt = await tx2Amb.wait();

        // todo check erc20 balance changed
        // todo use truffle helpers for catch events

        let events1Amb: any = await getEvents(receipt1Amb);
        let events2Amb: any = await getEvents(receipt2Amb);

        // Checking that event_id increased
        expect(events2Amb[0].args.event_id).eq(events1Amb[0].args.event_id.add("1"));
      });

      describe("Token addresses", () => {
        const token1 = "0x0000000000000000000000000000000000000001";
        const token2 = "0x0000000000000000000000000000000000000002";
        const token3 = "0x0000000000000000000000000000000000000003";
        const token4 = "0x0000000000000000000000000000000000000004";

        it("add tokens", async () => {
          await bridge.tokensAdd(token1, token2);

          expect(await bridge.tokenAddresses(token1)).eq(token2);
        });

        it("add tokens batch", async () => {
          await bridge.tokensAddBatch([token1, token2], [token3, token4]);

          expect(await bridge.tokenAddresses(token1)).eq(token3);
          expect(await bridge.tokenAddresses(token2)).eq(token4);
        });

        it("remove tokens", async () => {
          await bridge.tokensAdd(token1, token2);

          await bridge.tokensRemove(token1);

          expect(await bridge.tokenAddresses(token1)).eq(ethers.constants.AddressZero);
        });

        it("remove tokens batch", async () => {
          await bridge.tokensAddBatch([token1, token2], [token3, token4]);

          await bridge.tokensRemoveBatch([token1, token2]);

          expect(await bridge.tokenAddresses(token1)).eq(ethers.constants.AddressZero);
          expect(await bridge.tokenAddresses(token2)).eq(ethers.constants.AddressZero);
        });
      });

      it("Test fee", async () => {
        await mockERC20.mint(owner, 5);
        await mockERC20.increaseAllowance(bridge.address, 5);

        await bridge.changeFeeRecipient(user);
        await expect(
          () => bridge.withdraw(mockERC20.address, owner, 5, {value: 1000})
        ).to.changeEtherBalance(userS, 1000);

      });

      it("Test changeMinSafetyBlocks", async () => {
        await bridge.changeMinSafetyBlocks(20);
        expect(await bridge.minSafetyBlocks()).eq(20);
      });

      it("Test changeTimeframeSeconds", async () => {
        await bridge.changeTimeframeSeconds(20000);
        expect(await bridge.timeframeSeconds()).eq(20000);
      });

      it("Test changeLockTime", async () => {
        await bridge.changeLockTime(2000);
        expect(await bridge.lockTime()).eq(2000);
      });

      it("Test Pausable", async () => {
        await expect(bridge.removeLockedTransfers(0)).to.be.revertedWith("Pausable: not paused");

        await bridge.pause()
        expect(await bridge.paused()).to.be.true;

        await expect(bridge.connect(relayS).unlockTransfers(0)).to.be.revertedWith("Pausable: paused");

        await bridge.unpause()
        expect(await bridge.paused()).to.be.false;
      });

    });
  });

  it("Test setSideBridge", async () => {
    await ambBridge.setSideBridge(mockERC20.address);
    expect(await ambBridge.sideBridgeAddress()).eq(mockERC20.address);
  });


  it("Test Transfer lock/unlock", async () => {
    await mockERC20.mint(ambBridgeTest.address, 900);

    let data1 = [
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 36],
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 37],
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 38],
    ];

    let data2 = [
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 39],
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 40],
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 41],
    ];

    await ambBridgeTest.connect(relayS).lockTransfersTest(data1, 1);
    await ambBridgeTest.connect(relayS).lockTransfersTest(data2, 2);

    for (let i = 0; i < 3; i++) {
      let answer = await ambBridgeTest.getLockedTransferTest(1, i);
      expect(answer[0]).eq(data1[i][0]); // Check tokenAddress is correct
      expect(answer[1]).eq(data1[i][1]); // Check toAddress is correct
      expect(parseInt(answer[2]._hex, 16)).eq(data1[i][2]); // Check amount is correct
    }

    await nextTimeframe();

    await ambBridgeTest.connect(relayS).unlockTransfersTest(1);
    await ambBridgeTest.connect(relayS).unlockTransfersTest(2);
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

  it('CheckSignature', async () => {
    const hash = "0x74dfc2c4994f9393f823773a399849e0241c8f4c549f7e019960bd64b938b6ae"
    const signature = "0x44c08b83a120ad90f645f645f3fe1bc49dd88e703fce665de1f941c0cede65d81968ac2ad0b5bb9db7cb32a23064c199c0ab5378957f99b1c361fc3ed3b209eb1b";
    const needAddress = "0x90B2Ce3741188bCFCe25822113e93983ecacfcA0"

    expect(ethers.utils.recoverAddress(ethers.utils.arrayify(hash), signature)).eq(needAddress);
    await ethBridgeTest.CheckSignatureTest(needAddress, hash, signature)
  });


  let currentTimeframe = Math.floor(Date.now() / 14400);
  const nextTimeframe = async (amount = 1) => {
    currentTimeframe += amount;
    const timestamp = currentTimeframe * 14400 + amount * 14400;
    await network.provider.send("evm_setNextBlockTimestamp", [timestamp]);
  };

  it('Test wrap_withdraw in AmbBridge', async () => {
    await ambBridge.setAmbWrapper(wAmb.address);

    await expect(() => ambBridge.wrap_withdraw(wAmb.address, user, 1, {value: 1050}))
        .to.changeTokenBalance(wAmb, ownerS, 50);

    await ambBridge.wrap_withdraw(wAmb.address, user, 2, {value: 1001});
    await nextTimeframe();

    // will catch previous txs (because nextTimeframe happened)
    let tx1Amb: ContractTransaction = await ambBridge.wrap_withdraw(wAmb.address, user, 1337, {value: 1001});
    await ambBridge.wrap_withdraw(wAmb.address, user, 3, {value: 1001});
    await ambBridge.wrap_withdraw(wAmb.address, user, 4, {value: 1001});
    await nextTimeframe();

    // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
    let tx2Amb: ContractTransaction = await ambBridge.wrap_withdraw(wAmb.address, user, 1337, {value: 1001});
    await ambBridge.wrap_withdraw(wAmb.address, user, 5, {value: 1001});

    let receipt1Amb: ContractReceipt = await tx1Amb.wait();
    let receipt2Amb: ContractReceipt = await tx2Amb.wait();

    let events1Amb: any = await getEvents(receipt1Amb);
    let events2Amb: any = await getEvents(receipt2Amb);

    // Checking that event_id increased
    expect(events2Amb[0].args.event_id).eq(events1Amb[0].args.event_id.add("1"));
  });

  const getEvents = async (receipt: any) => {
    return receipt.events?.filter((x: any) => {
      return x.event == "Transfer";
    });
  };
});
