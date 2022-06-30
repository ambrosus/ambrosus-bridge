import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;


describe("Check Aura", () => {
  let ownerS: Signer;
  let owner: string;

  let ethBridge: Contract;


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);

    ethBridge = await ethers.getContract("CheckAuraTest", ownerS);

  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state
  });

  it("Test CheckAura (no changes in VS)", async () => {
    const proof = require("./fixtures/auraProof-staticVs.json");
    await ethBridge.checkAuraTest(proof, 10, "0x08e0dB5952df058E18dbCD6F3d9433Cfd6bbC18B", [
      "0x4c9785451bb2CA3E91B350C06bcB5f974cA33F79",
      "0x90B2Ce3741188bCFCe25822113e93983ecacfcA0",
      "0xAccdb7a2268BC4Af0a1898e725138888ba1Ca6Fc"
    ]);
  });

  it("Test CheckAura (changes in VS)", async () => {
    const proof = require("./fixtures/auraProof-changeVs.json");
    await ethBridge.checkAuraTest(proof, 2, "0xA7CB71C94f3fC5Ee2232f8E449461177bc23527B", [
          "0xdecA85befcC43ed1891758E37c35053aFF935AC1",
          "0x427933454115d6D55E8e24821d430F944d3eD936",
          "0x87a3d2CcacDe32f366Bd01bcbeB202643cD38A4E",
          "0x4682b2553F68a6C6d0182ac83425A1D0A0547337",
          "0xa45899BD58c4dE692883B3430B2e4a4CCE087c07",
          "0xA1c203F8B88F902b92cc96817382EC0b5dDAA77C",
          "0x716963005bf5b517cC7ACb4c8D99d7Dc1dC9A7c8",
          "0xaD5caf4A4B68eD66C2CD3A7d730Aee5747f31DFe",
          "0x6DD23d8c5c42c98194771218fB2aD465a8CFd55d",
          "0xf4B075fDF227219fF2f72fE87641aDCdFDc019BC",
          "0xc1E639642a242396C420C4880ABB3599Fb69d242",
          "0x4137e5c2D3a17E931F96Ef4eAe7F34985d4e6FED",
          "0x7B10BAEfA1bF7eDF72e1705b1e52dc66926a3bd8",
          "0xe58647A6dB76F3e7eaA6c2dB86239E73e99ebE03",
          "0x38A1835e04befEd507F3eF6b25D61f3E4BfbF9a1",
        ], {gasLimit: 40000000}
    );
  });


  it("Test submitValidatorSetChangesAura", async () => {
      const proof = require("./fixtures/auraProof-partChangeVs.json");
      await ethBridge.checkAuraTest(proof, 0, ethers.constants.AddressZero, [
              '0xF4Fc27eBDf978BC19f0E1cBdDC6875494b305AC4',
              '0x7eC889B72C145d0Ae82AAA1b816fe611b9Cf16B7',
              '0x7138bb1131C12e8e6687Cd29a1993F8E97991829',
              '0x271CE3a4c3778557A49c15152B2eB5151eD6eA0D',
          ], {gasLimit: 40000000}
      );
  });


  // it("Test CheckAura (changes in VS NEW)", async () => {
  //   const proof = require("./fixtures/auraProof-changeVs-NEW.json");
  //   await ethBridge.checkAuraTest(proof, 10, "0x08e0dB5952df058E18dbCD6F3d9433Cfd6bbC18B");
  // });


  it("Test blockHash", async () => {
    const block = require("./fixtures/BlockPoA-48879.json");

    const [bare, seal] = await ethBridge.blockHashTest(block);
    expect(bare).to.be.equal("0x36d67412a4917d85fc9334644fafb5e69ef71361c6ba17a9089d36e75918e3b3");
    expect(seal).to.be.equal("0x8579595d2c25916e0a465c24618e33df81e67e06be9b03fc433dd4a2114c4cf5");

    await ethBridge.blockHashTestPaid(block);
  });

  it('Test bytesToUintTest', async () => {
    expect(await ethBridge.bytesToUintTest("0xdeadbeef")).to.be.equal(3735928559);
  });


});
