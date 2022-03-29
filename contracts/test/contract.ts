import {deployments, ethers, getNamedAccounts, network} from "hardhat";
import type {Contract, ContractReceipt, ContractTransaction, Signer} from "ethers";
import {BigNumber} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

const ADMIN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const RELAY_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("RELAY_ROLE"));

describe("Contract", () => {
  let ownerS: Signer;

  let adminS: Signer;
  let relayS: Signer;
  let owner: string;
  let admin: string;
  let relay: string;
  let user1: string;
  let user2: string;
  let user3: string;

  let ethBridge: Contract;
  let ambBridge: Contract;
  let mockERC20: Contract;
  let ambBridgeTest: Contract;


  before(async () => {
    await deployments.fixture(["ethbridge", "ambbridge", "mocktoken", "ethash", "ambbridgetest"]);
    ({owner, admin, relay, user1, user2, user3} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    adminS = await ethers.getSigner(admin);
    relayS = await ethers.getSigner(relay);

    ethBridge = await ethers.getContract("EthBridge", ownerS);
    ambBridge = await ethers.getContract("AmbBridge", ownerS);
    mockERC20 = await ethers.getContract("MockERC20", ownerS);
    ambBridgeTest = await ethers.getContract("AmbBridgeTest", ownerS);

    await ambBridge.grantRole(ADMIN_ROLE, admin);
    await ambBridgeTest.grantRole(ADMIN_ROLE, admin);
    await ethBridge.grantRole(ADMIN_ROLE, admin);

    await ambBridge.grantRole(RELAY_ROLE, relay);
    await ambBridgeTest.grantRole(RELAY_ROLE, relay);
    await ethBridge.grantRole(RELAY_ROLE, relay);
  });

  beforeEach(async () => {
    await deployments.fixture(["ethbridge", "ambbridge", "ethash", "ambbridgetest"]); // reset contracts state
  });

  it("TestSign", async () => {
    // amb main net block 16021709

    const hash = "0x6a6249e2d7c2aca375164f7eb233dfa23a8c0ff7046520849894ee23dd026eaf";
    const msg = ethers.utils.arrayify(hash);
    const signature =
      "0xe437f2d7a70044b9b7c1d98cde917b28e7e0298ab26abdbeacd3fd81de6100af5a61014a397218b4fb36230867634db0e5e40d4c17dd8b69979e52eb2fb9217e1b";

    const address = ethers.utils.recoverAddress(msg, signature);
    expect(address).eq("0x7a3599b88EA2068268B763deE8B8703d21700DF3");
  });

  it("TestWithdraw timeframe", async () => {
    await ambBridge.withdraw(mockERC20.address, user1, 1, {value: 1000});
    await ambBridge.withdraw(mockERC20.address, user1, 2, {value: 1000});
    await ethBridge.withdraw(mockERC20.address, user1, 1, {value: 1000});
    await ethBridge.withdraw(mockERC20.address, user1, 2, {value: 1000});
    await nextTimeframe();

    // will catch previous txs (because nextTimeframe happened)
    let tx1Amb: ContractTransaction = await ambBridge.withdraw(mockERC20.address, user1, 1337, {value: 1000});
    let tx1Eth: ContractTransaction = await ethBridge.withdraw(mockERC20.address, user1, 1337, {value: 1000});
    await ambBridge.withdraw(mockERC20.address, user1, 3, {value: 1000});
    await ambBridge.withdraw(mockERC20.address, user1, 4, {value: 1000});
    await ethBridge.withdraw(mockERC20.address, user1, 3, {value: 1000});
    await ethBridge.withdraw(mockERC20.address, user1, 4, {value: 1000});
    await nextTimeframe();

    // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
    let tx2Amb: ContractTransaction = await ambBridge.withdraw(mockERC20.address, user1, 1337, {value: 1000});
    let tx2Eth: ContractTransaction = await ethBridge.withdraw(mockERC20.address, user1, 1337, {value: 1000});
    await ambBridge.withdraw(mockERC20.address, user1, 5, {value: 1000});

    let receipt1Amb: ContractReceipt = await tx1Amb.wait();
    let receipt1Eth: ContractReceipt = await tx1Eth.wait();
    let receipt2Amb: ContractReceipt = await tx2Amb.wait();
    let receipt2Eth: ContractReceipt = await tx2Eth.wait();

    // todo check erc20 balance changed
    // todo use truffle helpers for catch events

    const getEvents = async (receipt: any) => {
      return receipt.events?.filter((x: any) => {
        return x.event == "Transfer";
      });
    };

    let events1Amb: any = await getEvents(receipt1Amb);
    let events1Eth: any = await getEvents(receipt1Eth);
    let events2Amb: any = await getEvents(receipt2Amb);
    let events2Eth: any = await getEvents(receipt2Eth);

    // Checking that event_id increased
    expect(events2Amb[0].args.event_id).eq(events1Amb[0].args.event_id.add("1"));
    expect(events2Eth[0].args.event_id).eq(events1Eth[0].args.event_id.add("1"));
  });

  describe("Token addresses", () => {
    const token1 = "0x0000000000000000000000000000000000000001";
    const token2 = "0x0000000000000000000000000000000000000002";
    const token3 = "0x0000000000000000000000000000000000000003";
    const token4 = "0x0000000000000000000000000000000000000004";

    // const tokenThisAddresses = [token1, token2];
    // const tokenSideAddresses = [token3, token4];

    it("add tokens", async () => {
      await ambBridge.connect(adminS).tokensAdd(token1, token2);
      await ethBridge.connect(adminS).tokensAdd(token1, token2);

      expect(await ambBridge.tokenAddresses(token1)).eq(token2);
      expect(await ethBridge.tokenAddresses(token1)).eq(token2);
    });

    it("add tokens batch", async () => {
      await ambBridge.connect(adminS).tokensAddBatch([token1, token2], [token3, token4]);
      await ethBridge.connect(adminS).tokensAddBatch([token1, token2], [token3, token4]);

      expect(await ambBridge.tokenAddresses(token1)).eq(token3);
      expect(await ambBridge.tokenAddresses(token2)).eq(token4);
      expect(await ethBridge.tokenAddresses(token1)).eq(token3);
      expect(await ethBridge.tokenAddresses(token2)).eq(token4);
    });

    it("remove tokens", async () => {
      await ambBridge.connect(adminS).tokensAdd(token1, token2);
      await ethBridge.connect(adminS).tokensAdd(token1, token2);

      await ambBridge.connect(adminS).tokensRemove(token1);
      await ethBridge.connect(adminS).tokensRemove(token1);

      expect(await ambBridge.tokenAddresses(token1)).eq(ethers.constants.AddressZero);
      expect(await ethBridge.tokenAddresses(token1)).eq(ethers.constants.AddressZero);
    });

    it("remove tokens batch", async () => {
      await ambBridge.connect(adminS).tokensAddBatch([token1, token2], [token3, token4]);
      await ethBridge.connect(adminS).tokensAddBatch([token1, token2], [token3, token4]);

      await ambBridge.connect(adminS).tokensRemoveBatch([token1, token2]);
      await ethBridge.connect(adminS).tokensRemoveBatch([token1, token2]);

      expect(await ambBridge.tokenAddresses(token1)).eq(ethers.constants.AddressZero);
      expect(await ambBridge.tokenAddresses(token2)).eq(ethers.constants.AddressZero);
      expect(await ethBridge.tokenAddresses(token1)).eq(ethers.constants.AddressZero);
      expect(await ethBridge.tokenAddresses(token2)).eq(ethers.constants.AddressZero);
    });
  });

  it("Test fee", async () => {
    await mockERC20.mint(owner, 5);
    await mockERC20.increaseAllowance(ambBridge.address, 5);

    const prevBalance = await getAccountBalance(user3);
    const feeValue = 1000;

    await ambBridge.connect(adminS).changeFeeRecipient(user3);
    await ambBridge.withdraw(mockERC20.address, user2, 5, {value: feeValue});

    const curBalance = await getAccountBalance(user3);

    expect(curBalance).eq(prevBalance + feeValue);
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

  it("Test changeMinSafetyBlocks", async () => {
    const expectedMinSafetyBlocks = 20;

    await ambBridge.connect(adminS).changeMinSafetyBlocks(expectedMinSafetyBlocks);

    const realMinSafetyBlocks = await ambBridge.minSafetyBlocks();

    expect(realMinSafetyBlocks).eq(expectedMinSafetyBlocks);
  });

  it("Test changeTimeframeSeconds", async () => {
    const expectedTimeframeSeconds = 20000;

    await ambBridge.connect(adminS).changeTimeframeSeconds(expectedTimeframeSeconds);

    const realTimeframeSeconds = await ambBridge.timeframeSeconds();

    expect(realTimeframeSeconds).eq(expectedTimeframeSeconds);
  });

  it("Test changeLockTime", async () => {
    const expectedLockTime = 2000;

    await ambBridge.connect(adminS).changeLockTime(expectedLockTime);

    const realLockTime = await ambBridge.lockTime();

    expect(realLockTime).eq(expectedLockTime);
  });

  it("Test setSideBridge", async () => {
    const expectedSideBridgeAddress = "0x00000000000000000000000000000000000b00BA";

    await ambBridge.connect(adminS).setSideBridge(expectedSideBridgeAddress);

    const realSideBridgeAddress = await ambBridge.sideBridgeAddress();

    expect(realSideBridgeAddress).eq(expectedSideBridgeAddress);
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


  let currentTimeframe = Math.floor(Date.now() / 14400);
  const nextTimeframe = async (amount = 1) => {
    currentTimeframe += amount;
    const timestamp = currentTimeframe * 14400 + amount * 14400;
    await network.provider.send("evm_setNextBlockTimestamp", [timestamp]);
  };

  const getAccountBalance = async (addr: string) => {
    return parseInt(
      await ethers.provider.getBalance(addr).then((BN: BigNumber) => {
        return BN._hex;
      }),
      16
    );
  };
});
