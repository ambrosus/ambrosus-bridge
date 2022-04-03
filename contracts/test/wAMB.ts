import { deployments, ethers, getNamedAccounts } from "hardhat";
import type { Contract, Signer } from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

describe("wAMB", () => {
    let ownerS: Signer;
    let owner: string;
    let user: string;
    let userS: Signer;

    let wAMB: Contract;

    const val = 1000;

    before(async () => {
        await deployments.fixture(["for_tests"]); // reset contracts state
        ({ owner, user } = await getNamedAccounts());
        ownerS = await ethers.getSigner(owner);
        userS = await ethers.getSigner(user);

        wAMB = await ethers.getContract("wAMB", ownerS);
    });

    describe("wrap/unwrap", () => {
        it("Test wrap", async () => {
            await wAMB.connect(userS).wrap({value: val});
            expect(await wAMB.balanceOf(user)).eq(val);
        });

        it("Test unwrap", async () => {
            await expect(
                () => wAMB.connect(userS).unwrap(val)
            ).to.changeEtherBalance(userS, 1000);
        });
    });
});