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
    await ethBridge.checkAuraTest(proof, 10, "0x08e0dB5952df058E18dbCD6F3d9433Cfd6bbC18B");
  });

  it("Test CheckAura (changes in VS)", async () => {
    const proof = require("./fixtures/auraProof-changeVs.json");
    await ethBridge.checkAuraTestVS(proof, 5, "0x48d5cE2A10438559a14D399ca510F4235315dc6e", [
        '0x45a9645fcd727C2CeE29b5945aB49D92564Af199',
        '0x4fE180D06096A1216d3ec97ceAcE28d9c255B348',
        '0xf46BA733b57Da081F27309496CD618cdaB4E12B0',
        '0xfCf23c040142999e1c9894b421544aaA5805a21B',
        '0x4940b222D0Ec0c737F1b4f84caeaC60c3021a256',
        '0xe3746c9406FF78C3854aECAF36769826a22F1C89'
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

  it('CheckSignature', async () => {
    const hash = "0x74dfc2c4994f9393f823773a399849e0241c8f4c549f7e019960bd64b938b6ae"
    const signature = "0x44c08b83a120ad90f645f645f3fe1bc49dd88e703fce665de1f941c0cede65d81968ac2ad0b5bb9db7cb32a23064c199c0ab5378957f99b1c361fc3ed3b209eb00";
    const needAddress = "0x90B2Ce3741188bCFCe25822113e93983ecacfcA0"
    expect(ethers.utils.recoverAddress(ethers.utils.arrayify(hash), signature)).eq(needAddress);
    await ethBridge.checkSignatureTest(needAddress, hash, signature)
  });

  it('Test bytesToUintTest', async () => {
    expect(await ethBridge.bytesToUintTest("0xdeadbeef")).to.be.equal(3735928559);
  });


});
