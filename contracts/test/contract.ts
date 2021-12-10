import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";


chai.should();
export const expect = chai.expect;


describe("Contract", () => {
  let ownerS: Signer;
  let owner: string;

  let ethBridge: Contract;
  let ambBridge: Contract;

  before(async () => {
    await deployments.fixture(["ethbridge", "ambbridge"]);
    ({owner} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);

    ethBridge = await ethers.getContract("EthBridge", ownerS);
    ambBridge = await ethers.getContract("AmbBridge", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["ethbridge", "ambbridge"]); // reset contracts state
  });



  it('TestAll', async () => {


    // all data from go relay, https://rinkeby.etherscan.io/address/0x295C2707319ad4BecA6b5bb4086617fD6F240CfE


    const block = ["0xf9025ea0bbc7f0e38929ec0a39456a5bae6cb5a6b78f8221d3847c1f769870538e45b8bda01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347940000000000000000000000000000000000000000a03a1574ced84331e0c923ee8b59e3157c9e6c31b0bbb7229e075d344f77c294a4a086528e2027ea96e35f4f0d8b6f267142c3a08514e70e64b810ee7851c30e7d2aa0",
      "0x0e19377a0bfb933115babbc12b4c8dbdac491e6b59346d1d311204af82c4dd08",
      "0xb9010000000001020010340000000080000000000000004005402000402010008008420150000000000000000000000808000000000020082260800810004000240001100000000090000000000028002000000004010001040000000800808002080802800000060040001080000400000800420000100400000104040010100000080802100002440000400100000000000000000480000000000402800008004000020000001040000c000001080000000004000080040100040000182200080100000404020000000000010010003000080250000000810000000000200800e00100100c00000000000001000100000000a0010000100000004000520000000081018395405e8401c9c3808325a6c784",
      "0x61b137ee",
      "0xb861696e667572612e696f000000000000000000000000000000000000000000000010541caf38c3485e0d9459ccd31e6031b32bda1465b75f534ffa4209242818490faa202b60842f53150f7f76d447f968421d181a03f4bdc2d7e58f0be5643e6e01a0000000000000000000000000000000000000000000000000000000000000000088000000000000000009",
    ];
    const proof = ["0xf9027220b9026e02f9026a0183236fc6b9010000000000000000000000000000000000000000000004002000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000002000000004000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000008000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f9015ff9015c9467374e81da10af874664fec38bacd9137774f58ef842a088d88376da47bd16f6e8e45064996475d1536ab348a4f812de2f22518e06ee2da034b75fc0c8cd8e0b812bc8a3e1250f2287c67acee53b934091d34f4e81a20733b90100",
      "0x",
      "0xf891a049f601f6a3f228fc1bb8b67d859421ca82a03760a7ea61dd48557b7e016a6a25a0",
      "0xa0dd79151e8e1cfaa71c6afa47276326e1b7712a072c8ed37b3c6410389a112360a095aa56b747376fc8f3966bab353cc1f021280c793c8465285d4a17d8867044c880808080808080808080808080",
      "0xf871a0543013b82f229db241cc109973d1c570b96a15014367f79d434a0a3ab88732cca0",
      "0x808080808080a0449d4b558ab1dd1bc88d1e456bf9c7cff9e09370c045086eb062515407e1fa3e8080808080808080",
    ];

    const events = [
      ["0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", 123],
      ["0x295c2707319ad4beca6b5bb4086617fd6f240cfe", "0x295c2707319ad4beca6b5bb4086617fd6f240cfe", 456],

    ];


    const hash = ethers.utils.arrayify(ethers.utils.keccak256(ethers.utils.concat(block)));
    const signature = await ownerS.signMessage(hash);

    block.push(signature);


    await ethBridge.TestAll([block], events, proof);

  });

  it('TestReceiptsProof', async () => {

    const receiptsRoot = "0x0e19377a0bfb933115babbc12b4c8dbdac491e6b59346d1d311204af82c4dd08"
    const eventToCheck = "0x000000000000000000000000a17d0240c839e5461fd897a943080e3420156d40"
    const proof = [
      "0xf9024f20b9024b02f902470183044f90b9010000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000008000000000000000000000000000000200000000000000010000000000008000000000000000000000000000000000000080000000000020000000000000400000800020000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000020000000000000800000000000000000000000004000000000000000000000000000002000000000000000000000000000000000000000000000000000020000010000000000000000000000000000000000000000000000000000000000000f9013cf89c94564c7dacc723f92f9287ade47f1fb553dc8986eaf884a08c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925a0",
      "0xa00000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000080f89c94564c7dacc723f92f9287ade47f1fb553dc8986eaf884a0ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3efa0000000000000000000000000a17d0240c839e5461fd897a943080e3420156d40a00000000000000000000000005dc90c1582ba6b2af0d539fc076c94c584653faca0000000000000000000000000000000000000000000000000000000000000000080",
      "0xf901f180a0712a86434306bd9a49808ad6b19fed18016356a6b3718db019017a29b4ddd690a0",
      "0xa0ffcd39ae80db397d871917099f4d16c1d964b5968f622953519798c04c1fb6c4a0477a3b710fb3df4c8be2e03b80a999fe369b3ac955cafbe56a74b1518e73b325a04381d92bff3bc6b706f88f6c9020c3fe02929156e6d5e4a704d826cf11158b21a04cb1b1f4e5e34957430326ae5f082014e7c8258adcc9342f2c711fedb2160e89a00d427a16885eb1878991c4283d0804f56faaf392795566505001d88c3ca865eda03d5c0206dd0d877bf10ec8bee3177c6dfa88caf923408962e8783c39b4084607a078ce17ac5542c38c0975dd6dbb73040758fe160e8c8889bd7a98a131cfee92d5a0bac40bf499bace323f8714e8bccb68ed46776307371ab18907f81b1ca496c0c7a00880578c4c5e2a5228cb589249762bb8fe69fdcbef276a6d67b00871abb4be94a0f1db54679b31c6c2197ef4ce99a043f810a85c1ee8d16ad09794073f8559d935a0c2d20c5ffebf3b01ad547d443449258a4405d681e8f6cbf7ae3edd94651a27c4a01eb60f9187783bf53e7900ddf527567ca2ea4291494a2fa442b482b72f5d842aa088ede6d13f18c64710d981a882f88d1bdb0bb94559aef758016a9c5ad4baa48e80",
      "0xf871a0",
      "0xa05621bdcbf22658e545e8288486f0a32cbfd514319b4721fd629b28f30eb1f4f4808080808080a0449d4b558ab1dd1bc88d1e456bf9c7cff9e09370c045086eb062515407e1fa3e8080808080808080"]


    await ethBridge.TestReceiptsProof(proof, eventToCheck, receiptsRoot);


  });


  it('TestBloom', async () => {

    const topicHash = "0x000000000000000000000000a17d0240c839e5461fd897a943080e3420156d40"
    const bloom = "0x00000001020010340000000080000000000000004005402000402010008008420150000000000000000000000808000000000020082260800810004000240001100000000090000000000028002000000004010001040000000800808002080802800000060040001080000400000800420000100400000104040010100000080802100002440000400100000000000000000480000000000402800008004000020000001040000c000001080000000004000080040100040000182200080100000404020000000000010010003000080250000000810000000000200800e00100100c00000000000001000100000000a0010000100000004000520000000081"

    await ethBridge.TestBloom(bloom, topicHash);

  });


});
