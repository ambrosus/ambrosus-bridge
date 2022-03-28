import { deployments, ethers, getNamedAccounts, network } from "hardhat";
import type { Contract, ContractReceipt, ContractTransaction, Signer } from "ethers";
import { BigNumber } from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

const adminRole = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const relayRole = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("RELAY_ROLE"));

const FEE = 1000;

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
  let ethash: Contract;

  let ambBridgeTest: Contract;

  let tokenThisAddresses: string[];
  let tokenSideAddresses: string[];
  let token1: string;
  let token2: string;
  let token3: string;
  let token4: string;
  let token5: string;
  let token6: string;

  before(async () => {
    await deployments.fixture(["ethbridge", "ambbridge", "mocktoken", "ethash", "ambbridgetest"]);
    ({ owner, admin, relay, user1, user2, user3 } = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    adminS = await ethers.getSigner(admin);
    relayS = await ethers.getSigner(relay);

    ethBridge = await ethers.getContract("EthBridge", ownerS);
    ambBridge = await ethers.getContract("AmbBridge", ownerS);
    mockERC20 = await ethers.getContract("MockERC20", ownerS);
    ethash = await ethers.getContract("Ethash", ownerS);
    ambBridgeTest = await ethers.getContract("AmbBridgeTest", ownerS);

    await ambBridge.grantRole(adminRole, admin);
    await ambBridgeTest.grantRole(adminRole, admin);
    await ethBridge.grantRole(adminRole, admin);

    await ambBridge.grantRole(relayRole, relay);
    await ambBridgeTest.grantRole(relayRole, relay);
    await ethBridge.grantRole(relayRole, relay);
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
    const doWithdraw = async (bridge: Contract, tokenAddress: string, toAddress: string) => {
      await bridge.withdraw(tokenAddress, toAddress, 1, { value: FEE });
      await bridge.withdraw(tokenAddress, toAddress, 2, { value: FEE });
    }

    await doWithdraw(ambBridge, user1, user2);
    await doWithdraw(ethBridge, user1, user2);
    await nextTimeframe();

    // will catch previous txs (because nextTimeframe happened)
    let tx1Amb: ContractTransaction = await ambBridge.withdraw(user1, user2, 1337, { value: FEE });
    let tx1Eth: ContractTransaction = await ethBridge.withdraw(user1, user2, 1337, { value: FEE });
    await doWithdraw(ambBridge, user1, user2);
    await doWithdraw(ethBridge, user1, user2);
    await nextTimeframe();

    // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
    let tx2Amb: ContractTransaction = await ambBridge.withdraw(user1, user2, 1337, { value: FEE });
    let tx2Eth: ContractTransaction = await ethBridge.withdraw(user1, user2, 1337, { value: FEE });
    await ambBridge.withdraw(user1, user2, 5, { value: FEE });

    let receipt1Amb: ContractReceipt = await tx1Amb.wait();
    let receipt1Eth: ContractReceipt = await tx1Eth.wait();
    let receipt2Amb: ContractReceipt = await tx2Amb.wait();
    let receipt2Eth: ContractReceipt = await tx2Eth.wait();

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
    before(async () => {
      tokenThisAddresses = [
        ethers.utils.getAddress("0x195c2707319ad4beca6b5bb4086617fd6f240cfe"),
        ethers.utils.getAddress("0x295c2707319ad4beca6b5bb4086617fd6f240cfe"),
        ethers.utils.getAddress("0x395c2707319ad4beca6b5bb4086617fd6f240cfe"),
      ];
      tokenSideAddresses = [
        ethers.utils.getAddress("0x495c2707319ad4beca6b5bb4086617fd6f240cfe"),
        ethers.utils.getAddress("0x595c2707319ad4beca6b5bb4086617fd6f240cfe"),
        ethers.utils.getAddress("0x695c2707319ad4beca6b5bb4086617fd6f240cfe"),
      ];
      token1 = ethers.utils.getAddress("0x13372707319ad4beca6b5bb4086617fd6f240cfe");
      token2 = ethers.utils.getAddress("0x12282707319ad4beca6b5bb4086617fd6f240cfe");
      token3 = ethers.utils.getAddress("0x11192707319ad4beca6b5bb4086617fd6f240cfe");
      token4 = ethers.utils.getAddress("0x10002707319ad4beca6b5bb4086617fd6f240cfe");
      token5 = ethers.utils.getAddress("0x99992707319ad4beca6b5bb4086617fd6f240cfe");
      token6 = ethers.utils.getAddress("0x88882707319ad4beca6b5bb4086617fd6f240cfe");
    });

    it("should token for this address == to token for side address", async () => {
      for (let i = 0; i < tokenThisAddresses.length; i++) {
        expect(await ambBridge.tokenAddresses(tokenThisAddresses[i])).eq(tokenSideAddresses[i]);
        expect(await ethBridge.tokenAddresses(tokenThisAddresses[i])).eq(tokenSideAddresses[i]);
      }
    });

    it("add tokens", async () => {
      await ambBridge.connect(adminS).tokensAdd(token1, token2);
      await ethBridge.connect(adminS).tokensAdd(token1, token2);

      expect(await ambBridge.tokenAddresses(token1)).eq(token2);
      expect(await ethBridge.tokenAddresses(token1)).eq(token2);
    });

    it("add tokens batch", async () => {
      await ambBridge.connect(adminS).tokensAddBatch([token3, token4], [token5, token6]);
      await ethBridge.connect(adminS).tokensAddBatch([token3, token4], [token5, token6]);

      expect(await ambBridge.tokenAddresses(token3)).eq(token5);
      expect(await ambBridge.tokenAddresses(token4)).eq(token6);
      expect(await ethBridge.tokenAddresses(token3)).eq(token5);
      expect(await ethBridge.tokenAddresses(token4)).eq(token6);
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
      await ambBridge.connect(adminS).tokensAddBatch([token3, token4], [token5, token6]);
      await ethBridge.connect(adminS).tokensAddBatch([token3, token4], [token5, token6]);

      await ambBridge.connect(adminS).tokensRemoveBatch([token3, token4]);
      await ethBridge.connect(adminS).tokensRemoveBatch([token3, token4]);

      expect(await ambBridge.tokenAddresses(token3)).eq(ethers.constants.AddressZero);
      expect(await ambBridge.tokenAddresses(token4)).eq(ethers.constants.AddressZero);
      expect(await ethBridge.tokenAddresses(token3)).eq(ethers.constants.AddressZero);
      expect(await ethBridge.tokenAddresses(token4)).eq(ethers.constants.AddressZero);
    });
  });

  it("Test Ethash PoW", async () => {
    const blockPoW = require("../../relay/cmd/dump-test-data/BlockPoW-14257704.json");
    const epoch = require("../../relay/cmd/dump-test-data/epoch-475.json");

    await submitEpochData(ambBridge, epoch);
    expect(await ambBridge.isEpochDataSet(epoch.Epoch)).to.be.true;

    await ambBridge.verifyEthash(blockPoW);
  });

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

  it("Test fee", async () => {
    const prevBalance = await getAccountBalance(user3);
    const feeValue = 1000;

    await ambBridge.connect(adminS).changeFeeRecipient(user3);
    await ambBridge.withdraw(user1, user2, 5, {
      value: feeValue,
    });

    const curBalance = await getAccountBalance(user3);

    expect(curBalance).eq(prevBalance + feeValue);
  });

  it("Test fee with changing 'fee' variable", async () => {
    const prevBalance = await getAccountBalance(user3);
    const feeValue = 1000;

    await ambBridge.connect(adminS).changeFeeRecipient(user3);
    await ambBridge.connect(adminS).changeFee(feeValue);
    await ambBridge.withdraw(user1, user2, 5, {
      value: feeValue,
    });

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
    const expectedSideBridgeAddress = ethers.utils.getAddress("0x13372707319ad4beca6b5bb4086617fd6f240cfe");

    await ambBridge.connect(adminS).setSideBridge(expectedSideBridgeAddress);

    const realSideBridgeAddress = await ambBridge.sideBridgeAddress();

    expect(realSideBridgeAddress).eq(expectedSideBridgeAddress);
  });

  it("Test blockHash", async () => {
    const blockPoW = require("../../relay/cmd/dump-test-data/BlockPoW-14257704.json");
    const expectedBlockHash = "0xc4ca0efd5d528d67691abd9e10e9d4ca570f16235779e1f314b036caa5b455a1";

    const realBlockHash = await ambBridgeTest.blockHashTest(blockPoW);
    expect(realBlockHash).eq(expectedBlockHash);
  });

  let currentTimeframe = Math.floor(Date.now() / 14400);
  const nextTimeframe = async (amount = 1) => {
    currentTimeframe += amount;
    const timestamp = currentTimeframe * 14400 + amount * 14400;
    await network.provider.send("evm_setNextBlockTimestamp", [timestamp]);
  };

  const submitEpochData = async (ethashContractInstance: Contract, epoch_: any) => {
    let epoch = epoch_;
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
          epoch.Epoch,
          epoch.FullSizeIn128Resolution,
          epoch.BranchDepth,
          nodes,
          start,
          mnlen
        );
        start = start + mnlen;
        nodes = [];
      }
      index++;
    }
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
