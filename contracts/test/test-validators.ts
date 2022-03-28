import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";


chai.should();
export const expect = chai.expect;


describe("Validators", () => {
  let ownerS: Signer;
  let owner: string;

  let ethBridge: Contract;
  let ambBridge: Contract;
  let mockERC20: Contract;

  before(async () => {
    await deployments.fixture(["ethbridge", "ambbridge", "mocktoken"]);
    ({owner} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    ethBridge = await ethers.getContract("EthBridge", ownerS);
    ambBridge = await ethers.getContract("AmbBridge", ownerS);
    mockERC20 = await ethers.getContract("MockERC20", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["ethbridge", "ambbridge"]); // reset contracts state
  });


  it('Test checkpow', async () => {
    /*
    const block = require("./data-pow/block-pow-checkpow.json");
    const receipt_proof = require("./data-pow/receipt-proof-checkpow.json");
    const transfer = {
      tokenAddress: "0xc4b907fc242097D47eFd47f36eaee5Da2C239aDd",
      toAddress: "0x8FC84c829d9cB1982f2121F135624E25aac679A9",
      amount: 10
    }
    const powProof = [
      [],
      []  //
    ]

    await ambBridge.CheckPoW_(powProof, ethers.utils.getAddress("0x5493cA7e444F606BeE4E17748d3eF95716156d76"));
// todo


    // all data from go relay, https://rinkeby.etherscan.io/address/0x295C2707319ad4BecA6b5bb4086617fD6F240CfE

    const rlpBlocks = [
      ["0xf90213a0e5dd87c5db902a9288fc919bd5090790dd465d01cf888aeef552e55c6cc28222a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794f927a40c8b7f6e07c5af7fa2155b4864a4112b13a0772f8983a7c608e1c85c7edeaf7659ba3fed07f65f66de21a9972c543ff94bfea0", "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421", "0xa056e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000085", "0x97e3b9f2ef", "0x822b67821388808455bb7c3799476574682f76312e302e302f6c696e75782f676f312e342e32a08ba868beedbfd6bf3d2161aedeedeff27fc18f0efb0a7a129896e5db11844c4d88d7e030fb5bf8ee9c", ],
    ];

    const proof = [
        "0xf9027220b9026e02f9026a0183236fc6b9010000000000000000000000000000000000000000000004002000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000002000000004000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000008000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f9015ff9015c9467374e81da10af874664fec38bacd9137774f58ef842a088d88376da47bd16f6e8e45064996475d1536ab348a4f812de2f22518e06ee2da034b75fc0c8cd8e0b812bc8a3e1250f2287c67acee53b934091d34f4e81a20733b90100",
        "0x",
        "0xf891a049f601f6a3f228fc1bb8b67d859421ca82a03760a7ea61dd48557b7e016a6a25a0",
        "0xa0dd79151e8e1cfaa71c6afa47276326e1b7712a072c8ed37b3c6410389a112360a095aa56b747376fc8f3966bab353cc1f021280c793c8465285d4a17d8867044c880808080808080808080808080",
        "0xf871a0543013b82f229db241cc109973d1c570b96a15014367f79d434a0a3ab88732cca0",
        "0x808080808080a0449d4b558ab1dd1bc88d1e456bf9c7cff9e09370c045086eb062515407e1fa3e8080808080808080"
    ];

    const events = [
      ["0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", 123],
      ["0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", 456],
    ];


    for (let block of rlpBlocks) {
      const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.concat(block)));
      const signature = await ownerS.signMessage(hash);
      block.push(signature);
    }


    await ambBridge.CheckPoW_(rlpBlocks, events, proof);

    const testData = [
        "0xf9025ea0",
        "0x1cd7bc73b74bd8ab52f9acc422de0a29be3e263a74f9cffe3bc00a61d25870fc",
        "0xa01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a0dde85e2b8d314f59e054fac21ddcd03df33e50c5efe151c6be1a66d582e5c0b7a0dd9284831cb8028a6270caec4aac6e1f79a29162933bbfed953b676ce9085d72a0fb9a806f63a0e01440fad77be564089d59e73970e803656d66b91f74689bd25cb901000a800002024800221080418020004014000000800000000000100402008102024020800000001020010000000000000000000000022a90200000008000040a89100088504000000010080028010040022414030c20450000002800a4000000012980808002000018000a000000240800400000002400004100101030900000000000104080442421000100420600601400008680000210814006040060402010000400000000001408240808000400002000a004808050044424040004004a0000000802000000060000049080400008804a000080a0004202009201004060000200001081108000000600010000000101020000149022020000200000044003028395405f8401c951118329a25584",
        "0x61b137fd",
        "0xb861d883010a09846765746888676f312e31362e36856c696e757800000000000000bb92c3d7e9e0eac3ecd1b82d27f23bf8b5247784aa86afff820d4ec2f1a7d25b16144979532612c3d5fe7d1fcea3d8d65a41e1c698061f816c1ba19981ae6f8200a0000000000000000000000000000000000000000000000000000000000000000088000000000000000009",
    ];
    await ambBridge.testCheckPow(testData);
    */
  });

  it('Test  checkpoa', async () => {
// todo
    /*

    // all data from go relay, https://rinkeby.etherscan.io/address/0x295C2707319ad4BecA6b5bb4086617fD6F240CfE

    const rlpBlocks = [
      [
        "0xf9023b",
        "0xf901f3",
        "0xa084d550d22045e85a33e05f2873a305efd9be708f9ebfe494dbe8827bceea5485a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347947a3599b88ea2068268b763dee8b8703d21700df3a0324cdf89715ca3676749dab5432bce67ca098900cb28666350269b07886a4368a04747c79887ca71c79b728cfcbdb33dcece385f5cd4e11b27e9c43f3be9f84db5a0",
        "0xaf38d78710df0eea0c9db288894ce557784cdf9ec9ed214d0f22e041e5725229",
        "0xb901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000090fffffffffffffffffffffffffffffffe83f478cd837a120082dff284",
        "0x61bb4115",
        "0x9441706f6c6c6f2076322e372e322d737461626c65",
        "0x84138bd9d1b841",
        "0xe437f2d7a70044b9b7c1d98cde917b28e7e0298ab26abdbeacd3fd81de6100af5a61014a397218b4fb36230867634db0e5e40d4c17dd8b69979e52eb2fb9217e00"
      ],
    ];

    const proof = [
      "0xf9027220b9026e02f9026a0183236fc6b9010000000000000000000000000000000000000000000004002000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000002000000004000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000008000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f9015ff9015c9467374e81da10af874664fec38bacd9137774f58ef842a088d88376da47bd16f6e8e45064996475d1536ab348a4f812de2f22518e06ee2da034b75fc0c8cd8e0b812bc8a3e1250f2287c67acee53b934091d34f4e81a20733b90100",
      "0x",
      "0xf891a049f601f6a3f228fc1bb8b67d859421ca82a03760a7ea61dd48557b7e016a6a25a0",
      "0xa0dd79151e8e1cfaa71c6afa47276326e1b7712a072c8ed37b3c6410389a112360a095aa56b747376fc8f3966bab353cc1f021280c793c8465285d4a17d8867044c880808080808080808080808080",
      "0xf871a0543013b82f229db241cc109973d1c570b96a15014367f79d434a0a3ab88732cca0",
      "0x808080808080a0449d4b558ab1dd1bc88d1e456bf9c7cff9e09370c045086eb062515407e1fa3e8080808080808080"
    ];

    const events = [
      ["0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", 123],
      ["0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", 456],
    ];


    for (let block of rlpBlocks) {
      const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.concat(block)));
      const signature = await ownerS.signMessage(hash);
      block.push(signature);
    }


    await ethBridge.CheckPoA_(rlpBlocks, events, proof);
     */
  });

  it('CheckReceiptsProof', async () => {
    // todo
    /*
        const receiptsRoot = "0x0e19377a0bfb933115babbc12b4c8dbdac491e6b59346d1d311204af82c4dd08"
        const eventToCheck = "0x000000000000000000000000a17d0240c839e5461fd897a943080e3420156d40"
        const proof = [
          "0xf9024f20b9024b02f902470183044f90b9010000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000008000000000000000000000000000000200000000000000010000000000008000000000000000000000000000000000000080000000000020000000000000400000800020000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000020000000000000800000000000000000000000004000000000000000000000000000002000000000000000000000000000000000000000000000000000020000010000000000000000000000000000000000000000000000000000000000000f9013cf89c94564c7dacc723f92f9287ade47f1fb553dc8986eaf884a08c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925a0",
          "0xa00000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000080f89c94564c7dacc723f92f9287ade47f1fb553dc8986eaf884a0ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa0000000000000000000000000a17d0240c839e5461fd897a943080e3420156d40a00000000000000000000000005dc90c1582ba6b2af0d539fc076c94c584653faca0000000000000000000000000000000000000000000000000000000000000000080",
          "0xf901f180a0712a86434306bd9a49808ad6b19fed18016356a6b3718db019017a29b4ddd690a0",
          "0xa0ffcd39ae80db397d871917099f4d16c1d964b5968f622953519798c04c1fb6c4a0477a3b710fb3df4c8be2e03b80a999fe369b3ac955cafbe56a74b1518e73b325a04381d92bff3bc6b706f88f6c9020c3fe02929156e6d5e4a704d826cf11158b21a04cb1b1f4e5e34957430326ae5f082014e7c8258adcc9342f2c711fedb2160e89a00d427a16885eb1878991c4283d0804f56faaf392795566505001d88c3ca865eda03d5c0206dd0d877bf10ec8bee3177c6dfa88caf923408962e8783c39b4084607a078ce17ac5542c38c0975dd6dbb73040758fe160e8c8889bd7a98a131cfee92d5a0bac40bf499bace323f8714e8bccb68ed46776307371ab18907f81b1ca496c0c7a00880578c4c5e2a5228cb589249762bb8fe69fdcbef276a6d67b00871abb4be94a0f1db54679b31c6c2197ef4ce99a043f810a85c1ee8d16ad09794073f8559d935a0c2d20c5ffebf3b01ad547d443449258a4405d681e8f6cbf7ae3edd94651a27c4a01eb60f9187783bf53e7900ddf527567ca2ea4291494a2fa442b482b72f5d842aa088ede6d13f18c64710d981a882f88d1bdb0bb94559aef758016a9c5ad4baa48e80",
          "0xf871a0",
          "0xa05621bdcbf22658e545e8288486f0a32cbfd514319b4721fd629b28f30eb1f4f4808080808080a0449d4b558ab1dd1bc88d1e456bf9c7cff9e09370c045086eb062515407e1fa3e8080808080808080"]


        await ethBridge.CheckReceiptsProof(proof, eventToCheck, receiptsRoot);
        await ambBridge.CheckReceiptsProof(proof, eventToCheck, receiptsRoot);
      */
  });

  it("Test GetValidatorSet", async () => {
    // deploy script contains two validators in constructor
    const validator1 = ethers.utils.getAddress("0x11112707319ad4beca6b5bb4086617fd6f240cfe");
    const validator2 = ethers.utils.getAddress("0x22222707319ad4beca6b5bb4086617fd6f240cfe");
    const expectedValidatorSet = [validator1, validator2];

    expect(await ethBridge.GetValidatorSet()).eql(expectedValidatorSet)
  })
});
