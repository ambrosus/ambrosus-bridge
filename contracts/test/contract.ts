import {deployments, ethers, getNamedAccounts, network} from "hardhat";
import type {Contract, ContractReceipt, ContractTransaction, Signer} from "ethers";
import rlp from 'rlp'

import chai from "chai";
import {BigNumber} from "ethers";


chai.should();
export const expect = chai.expect;


describe("Contract", () => {
  let ownerS: Signer;
  let owner: string;

  let ethBridge: Contract;
  let ambBridge: Contract;
  let mockERC20: Contract;
  let ethash: Contract;

  let ambBridgeTest: Contract;

  before(async () => {
    await deployments.fixture(["ethbridge", "ambbridge", "mocktoken", "ethash", "ambbridgetest"]);
    ({owner} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);

    ethBridge = await ethers.getContract("EthBridge", ownerS);
    ambBridge = await ethers.getContract("AmbBridge", ownerS);
    mockERC20 = await ethers.getContract("MockERC20", ownerS);
    ethash = await ethers.getContract("Ethash", ownerS);
    ambBridgeTest = await ethers.getContract("AmbBridgeTest", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["ethbridge", "ambbridge", "ethash", "ambbridgetest"]); // reset contracts state
  });


  it('TestSign', async () => {
    // amb main net block 16021709

    const hash = "0x6a6249e2d7c2aca375164f7eb233dfa23a8c0ff7046520849894ee23dd026eaf"
    const msg = ethers.utils.arrayify(hash);
    const signature = "0xe437f2d7a70044b9b7c1d98cde917b28e7e0298ab26abdbeacd3fd81de6100af5a61014a397218b4fb36230867634db0e5e40d4c17dd8b69979e52eb2fb9217e1b";

    const address = ethers.utils.recoverAddress(msg, signature);
    expect(address).eq("0x7a3599b88EA2068268B763deE8B8703d21700DF3");
  });


  it("TestWithdraw timeframe", async () => {
    let [addr1, addr2] = await ethers.getSigners();

    await ambBridge.withdraw(addr1.address, addr2.address, 1, {value: 1000});
    await ambBridge.withdraw(addr1.address, addr2.address, 2, {value: 1000});
    await ethBridge.withdraw(addr1.address, addr2.address, 1, {value: 1000});
    await ethBridge.withdraw(addr1.address, addr2.address, 2, {value: 1000});
    await nextTimeframe();

    // will catch previous txs (because nextTimeframe happened)
    let tx1Amb: ContractTransaction = await ambBridge.withdraw(addr1.address, addr2.address, 1337, {value: 1000});
    let tx1Eth: ContractTransaction = await ethBridge.withdraw(addr1.address, addr2.address, 1337, {value: 1000});
    await ambBridge.withdraw(addr1.address, addr2.address, 3, {value: 1000});
    await ambBridge.withdraw(addr1.address, addr2.address, 4, {value: 1000});
    await ethBridge.withdraw(addr1.address, addr2.address, 3, {value: 1000});
    await ethBridge.withdraw(addr1.address, addr2.address, 4, {value: 1000});
    await nextTimeframe();

    // will catch previous txs started from tx1Amb/tx1Eth (because nextTimeframe happened)
    let tx2Amb: ContractTransaction = await ambBridge.withdraw(addr1.address, addr2.address, 1337, {value: 1000});
    let tx2Eth: ContractTransaction = await ethBridge.withdraw(addr1.address, addr2.address, 1337, {value: 1000});
    await ambBridge.withdraw(addr1.address, addr2.address, 5, {value: 1000});

    let receipt1Amb: ContractReceipt = await tx1Amb.wait();
    let receipt1Eth: ContractReceipt = await tx1Eth.wait();
    let receipt2Amb: ContractReceipt = await tx2Amb.wait();
    let receipt2Eth: ContractReceipt = await tx2Eth.wait();

    const getEvents = async (receipt: any) => {
      return receipt.events?.filter((x: any) => {
        return x.event == "Transfer"
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

  it("Test TokenAddresses", async () => {
    let [_, addr2] = await ethers.getSigners();
    let tokenThisAddresses = [ethers.utils.getAddress("0x195c2707319ad4beca6b5bb4086617fd6f240cfe"), ethers.utils.getAddress("0x295c2707319ad4beca6b5bb4086617fd6f240cfe"), ethers.utils.getAddress("0x395c2707319ad4beca6b5bb4086617fd6f240cfe")];
    let tokenSideAddresses = [ethers.utils.getAddress("0x495c2707319ad4beca6b5bb4086617fd6f240cfe"), ethers.utils.getAddress("0x595c2707319ad4beca6b5bb4086617fd6f240cfe"), ethers.utils.getAddress("0x695c2707319ad4beca6b5bb4086617fd6f240cfe")];
    for (let i = 0; i < tokenThisAddresses.length; i++) {
      expect(await ambBridge.tokenAddresses(tokenThisAddresses[i])).eq(tokenSideAddresses[i]);
      expect(await ethBridge.tokenAddresses(tokenThisAddresses[i])).eq(tokenSideAddresses[i]);
    }

    let hashAdmin = await ambBridge.ADMIN_ROLE();
    await ambBridge.grantRole(hashAdmin, addr2.address);
    await ethBridge.grantRole(hashAdmin, addr2.address);

    let first = ethers.utils.getAddress("0x13372707319ad4beca6b5bb4086617fd6f240cfe");
    let second = ethers.utils.getAddress("0x12282707319ad4beca6b5bb4086617fd6f240cfe")
    let third = ethers.utils.getAddress("0x11192707319ad4beca6b5bb4086617fd6f240cfe")
    let fourth = ethers.utils.getAddress("0x10002707319ad4beca6b5bb4086617fd6f240cfe")
    let fifth = ethers.utils.getAddress("0x99992707319ad4beca6b5bb4086617fd6f240cfe")
    let sixth = ethers.utils.getAddress("0x88882707319ad4beca6b5bb4086617fd6f240cfe")
    await ambBridge.connect(addr2).tokensAdd(first, second);
    expect(await ambBridge.tokenAddresses(first)).eq(second);
    await ethBridge.connect(addr2).tokensAdd(first, second);
    expect(await ethBridge.tokenAddresses(first)).eq(second);

    // batch
    await ambBridge.connect(addr2).tokensAddBatch([third, fourth], [fifth, sixth]);
    expect(await ambBridge.tokenAddresses(third)).eq(fifth);
    expect(await ambBridge.tokenAddresses(fourth)).eq(sixth);
    await ethBridge.connect(addr2).tokensAddBatch([third, fourth], [fifth, sixth]);
    expect(await ethBridge.tokenAddresses(third)).eq(fifth);
    expect(await ethBridge.tokenAddresses(fourth)).eq(sixth);

    await ambBridge.connect(addr2).tokensRemove(first);
    expect(await ambBridge.tokenAddresses(first)).eq("0x0000000000000000000000000000000000000000");
    await ethBridge.connect(addr2).tokensRemove(first);
    expect(await ethBridge.tokenAddresses(first)).eq("0x0000000000000000000000000000000000000000");

    // batch
    await ambBridge.connect(addr2).tokensRemoveBatch([third, fourth]);
    expect(await ambBridge.tokenAddresses(third)).eq("0x0000000000000000000000000000000000000000");
    expect(await ambBridge.tokenAddresses(fourth)).eq("0x0000000000000000000000000000000000000000");
    await ethBridge.connect(addr2).tokensRemoveBatch([third, fourth]);
    expect(await ethBridge.tokenAddresses(third)).eq("0x0000000000000000000000000000000000000000");
    expect(await ethBridge.tokenAddresses(fourth)).eq("0x0000000000000000000000000000000000000000");
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
    let [addr1, addr2, addr3, addr4] = await ethers.getSigners();

    let hashAdmin = await ambBridge.ADMIN_ROLE();
    await ambBridge.grantRole(hashAdmin, addr3.address);

    const prevBalance = await getAccountBalance(addr4);
    const feeValue = 1000;

    await ambBridge.connect(addr3).changeFeeRecipient(addr4.address);
    await ambBridge.withdraw(addr1.address, addr2.address, 5, {value: feeValue});

    const curBalance = await getAccountBalance(addr4);

    expect(curBalance).eq(prevBalance + feeValue);
  });

  it("Test fee with changing 'fee' variable", async () => {
    let [addr1, addr2, addr3, addr4] = await ethers.getSigners();

    let hashAdmin = await ambBridge.ADMIN_ROLE();
    await ambBridge.grantRole(hashAdmin, addr3.address);

    const prevBalance = await getAccountBalance(addr4);
    const feeValue = 4000;

    await ambBridge.connect(addr3).changeFeeRecipient(addr4.address);
    await ambBridge.connect(addr3).changeFee(feeValue);
    await ambBridge.withdraw(addr1.address, addr2.address, 5, {value: feeValue});

    const curBalance = await getAccountBalance(addr4);

    expect(curBalance).eq(prevBalance + feeValue);
  });


  it("Test Transfer lock/unlock", async () => {
    let [_, __, addr3] = await ethers.getSigners();

    let hashRelay = await ambBridgeTest.RELAY_ROLE();
    await ambBridgeTest.grantRole(hashRelay, addr3.address);

    await mockERC20.mint(ambBridgeTest.address, 900);

    let data1 = [
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 36],
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 37],
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 38]
    ]

    let data2 = [
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 39],
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 40],
      [mockERC20.address, "0x5B38Da6a701c568545dCfcB03FcB875f56beddC4", 41]
    ]

    await ambBridgeTest.connect(addr3).lockTransfersTest(data1, 1);
    await ambBridgeTest.connect(addr3).lockTransfersTest(data2, 2);

    for (let i = 0; i < 3; i++) {
      let answer  = await ambBridgeTest.getLockedTransferTest(1, i);
      expect(answer[0]).eq(data1[i][0]); // Check tokenAddress is correct
      expect(answer[1]).eq(data1[i][1]); // Check toAddress is correct
      expect(parseInt(answer[2]._hex, 16)).eq(data1[i][2]); // Check amount is correct
    }

    await nextTimeframe();

    await ambBridgeTest.connect(addr3).unlockTransfersTest(1);
    await ambBridgeTest.connect(addr3).unlockTransfersTest(2);
  });

  let currentTimeframe = Math.floor(Date.now() / 14400);
  const nextTimeframe = async (amount = 1) => {
    currentTimeframe += amount;
    const timestamp = currentTimeframe * 14400 + amount * 14400;
    await network.provider.send("evm_setNextBlockTimestamp", [timestamp]);
  }

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
        await ethashContractInstance.setEpochData(epoch.Epoch, epoch.FullSizeIn128Resolution, epoch.BranchDepth, nodes, start, mnlen);
        start = start + mnlen;
        nodes = [];
      }
      index++;
    }
  };

  const getAccountBalance = async (addr: any) => {
    return parseInt(await ethers.provider.getBalance(addr.address).then(
        (BN: BigNumber) => {return BN._hex}), 16);
  }
});
