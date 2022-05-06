import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;


const testData = require("../test/fixtures/posaProof.json");

describe("Check PoSA", () => {
    let ownerS: Signer;
    let owner: string;

    let bscBridge: Contract;


    before(async () => {
        await deployments.fixture(["for_tests"]);
        ({owner} = await getNamedAccounts());
        ownerS = await ethers.getSigner(owner);

        bscBridge = await ethers.getContract("CheckPoSATest", ownerS);
    });

    beforeEach(async () => {
        await deployments.fixture(["for_tests"]); // reset contracts state
    });

    it("Test CheckPoSA (no changes in VS)", async () => {
        // todo
    });

    it("Test CheckPoSA (changes in VS)", async () => {
        // todo
    });

    it("Test blockHash", async () => {
        const block = testData.blocks[0];

        const [bare, seal] = await bscBridge.blockHashTest(block);
        expect(bare).to.be.equal("0x1d0a6ca42217dc9f0560840b3eb91a3879b836cb7ec5a8055e265a520e6839d0");
        expect(seal).to.be.equal("0x12177fc4566157a0e6e71a60373b1c0587878c668587eea9ac2dfcaf2a44f080");

        await bscBridge.blockHashTestPaid(block);
    });

    it('CheckSignature', async () => {
        const hash = "0x1d0a6ca42217dc9f0560840b3eb91a3879b836cb7ec5a8055e265a520e6839d0";
        const signature = "0x5c1974f609035dc81319f058a8b9428b7ce26b366fadf9768b8ca19e3014c759467d732731a58a2ad9f3e9efedc56275427cd4a2fd7a6de59007b0bdb2e95f7d00";
        const needAddress = "0xc89C669357D161d57B0b255C94eA96E179999919";
        expect(ethers.utils.recoverAddress(ethers.utils.arrayify(hash), signature)).eq(needAddress);

        expect(await bscBridge.checkSignatureTest(hash, signature)).eq(needAddress);
    });
});