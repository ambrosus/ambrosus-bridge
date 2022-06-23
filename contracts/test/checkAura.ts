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
    await ethBridge.checkAuraTest(proof, 5, "0x48d5cE2A10438559a14D399ca510F4235315dc6e", [
        '0x45a9645fcd727C2CeE29b5945aB49D92564Af199',
        '0x4fE180D06096A1216d3ec97ceAcE28d9c255B348',
        '0xf46BA733b57Da081F27309496CD618cdaB4E12B0',
        '0xfCf23c040142999e1c9894b421544aaA5805a21B',
        '0x4940b222D0Ec0c737F1b4f84caeaC60c3021a256',
        '0xe3746c9406FF78C3854aECAF36769826a22F1C89'
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
