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
        const proof = require("./fixtures/posaProof-staticEpoch.json");
        await bscBridge.checkPoSATest(proof, 10, "0x944867B67cB2C28302C26df12B8aA01cb32F53Dc",
            ["0x049153b8dae0a232ac90d20c78f1a5d1de7b7dc5",
            "0x1284214b9b9c85549ab3d2b972df0deef66ac2c9",
            "0x35552c16704d214347f29fa77f77da6d75d7c752",
            "0x980a75ecd1309ea12fa2ed87a8744fbfc9b863d5",
            "0xa2959d3f95eae5dc7d70144ce1b73b403b7eb6e0",
            "0xb71b214cb885500844365e95cd9942c7276e7fd8",
            "0xf474cf03cceff28abc65c9cbae594f725c80e12d"], 100374, 97);
    });

    it("Test CheckPoSA (one change in VS)", async () => {
        const proof = require("./fixtures/posaProof-oneEpochChange.json");
        await bscBridge.checkPoSATest(proof, 10, "0x0953C80f775d36DC7CfbF32E6e2905FF040A354c",
            ["0x049153b8dae0a232ac90d20c78f1a5d1de7b7dc5",
                "0x1284214b9b9c85549ab3d2b972df0deef66ac2c9",
                "0x35552c16704d214347f29fa77f77da6d75d7c752",
                "0x96c5d20b2a975c050e4220be276ace4892f4b41a",
                "0x980a75ecd1309ea12fa2ed87a8744fbfc9b863d5",
                "0xa2959d3f95eae5dc7d70144ce1b73b403b7eb6e0",
                "0xb71b214cb885500844365e95cd9942c7276e7fd8",
                "0xf474cf03cceff28abc65c9cbae594f725c80e12d"], 100778, 97);
    });

    it("Test CheckPoSA (many changes in VS)", async () => {
        const proof = require("./fixtures/posaProof-manyEpochChanges.json");
        await bscBridge.checkPoSATest(proof, 10, "0x905D371a767024d485ecAb0c79871fbf9a14487d",
            ["0x049153b8dae0a232ac90d20c78f1a5d1de7b7dc5",
                "0x1284214b9b9c85549ab3d2b972df0deef66ac2c9",
                "0x35552c16704d214347f29fa77f77da6d75d7c752",
                "0x96c5d20b2a975c050e4220be276ace4892f4b41a",
                "0x980a75ecd1309ea12fa2ed87a8744fbfc9b863d5",
                "0xa2959d3f95eae5dc7d70144ce1b73b403b7eb6e0",
                "0xb71b214cb885500844365e95cd9942c7276e7fd8",
                "0xf474cf03cceff28abc65c9cbae594f725c80e12d"], 100782, 97, {gasLimit: 40000000});
    });

    it("Test submitValidatorSetChangesPoSA", async () => {
        const proof = require("./fixtures/posaProof-partEpochChanges.json");
        await bscBridge.checkPoSATest(proof, 0, ethers.constants.AddressZero,
            ["0x049153b8dae0a232ac90d20c78f1a5d1de7b7dc5",
                "0x1284214b9b9c85549ab3d2b972df0deef66ac2c9",
                "0x35552c16704d214347f29fa77f77da6d75d7c752",
                "0x980a75ecd1309ea12fa2ed87a8744fbfc9b863d5",
                "0xa2959d3f95eae5dc7d70144ce1b73b403b7eb6e0",
                "0xb71b214cb885500844365e95cd9942c7276e7fd8",
                "0xf474cf03cceff28abc65c9cbae594f725c80e12d"], 100500, 97, {gasLimit: 40000000});
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
