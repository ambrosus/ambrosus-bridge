import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

const ADMIN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const RELAY_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("RELAY_ROLE"));

describe("Check Aura", () => {
  let ownerS: Signer;
  let relayS: Signer;
  let owner: string;
  let relay: string;

  let ethBridge: Contract;


  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, relay} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    relayS = await ethers.getSigner(relay);

    ethBridge = await ethers.getContract("EthBridgeTest", ownerS);

    await ethBridge.grantRole(ADMIN_ROLE, owner);
    await ethBridge.grantRole(RELAY_ROLE, relay);
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state
  });

  it("Test CheckAura (no changes in VS)", async () => {
    const proof = require("./fixtures/auraProof-staticVs.json");
    await ethBridge.CheckAuraTest(proof, 10, "0x08e0dB5952df058E18dbCD6F3d9433Cfd6bbC18B", "0x0000000000000000000000000000000000000F00");
  });

  it("Test CheckAura (changes in VS)", async () => {
    const proof = require("./fixtures/auraProof-changeVs.json");
    await ethBridge.CheckAuraTest(proof, 10, "0x08e0dB5952df058E18dbCD6F3d9433Cfd6bbC18B", "0x0000000000000000000000000000000000000F00");
  });

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
    await ethBridge.CheckSignatureTest(needAddress, hash, signature)
  });

  it('Test bytesToUintTest', async () => {
    expect(await ethBridge.bytesToUintTest("0xdeadbeef")).to.be.equal(3735928559);
  });


});
