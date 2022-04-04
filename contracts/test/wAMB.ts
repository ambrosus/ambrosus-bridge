import { deployments, ethers, getNamedAccounts } from "hardhat";
import type { Contract, Signer } from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

describe("wAMB", () => {
    let userS: Signer;
    let user: string;

    let wAMB: Contract;

    const val = 1000;

    before(async () => {
        await deployments.fixture(["for_tests"]);
        ({ user } = await getNamedAccounts());
        userS = await ethers.getSigner(user);

        wAMB = await ethers.getContract("wAMB", userS);
    });

    describe("wrap/unwrap", () => {
        it("Test wrap", async () => {
            await wAMB.wrap({value: val});
            expect(await wAMB.balanceOf(user)).eq(val);
        });

        it("Test unwrap", async () => {
            await expect(
                () => wAMB.unwrap(val)
            ).to.changeEtherBalance(userS, val);
        });
    });
});
