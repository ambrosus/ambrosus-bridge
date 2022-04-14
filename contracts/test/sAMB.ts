import { deployments, ethers, getNamedAccounts } from "hardhat";
import type { Contract, Signer } from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

describe("sAMB", () => {
    let userS: Signer;
    let user: string;

    let sAMB: Contract;

    const val = 1000;

    before(async () => {
        await deployments.fixture(["for_tests"]);
        ({ user } = await getNamedAccounts());
        userS = await ethers.getSigner(user);

        sAMB = await ethers.getContract("sAMB", userS);
    });

    describe("wrap/unwrap", () => {
        it("Test wrap", async () => {
            await sAMB.deposit({value: val});
            expect(await sAMB.balanceOf(user)).eq(val);
        });

        it("Test unwrap", async () => {
            await expect(
                () => sAMB.withdraw(val)
            ).to.changeEtherBalance(userS, val);
        });
    });
});
