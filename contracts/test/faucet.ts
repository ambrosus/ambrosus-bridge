import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

describe("Faucet", () => {
  let ownerS: Signer;
  let owner: string;

  let userS: Signer;
  let user: string;

  let faucet: Contract;

  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, user} = await getNamedAccounts());
    userS = await ethers.getSigner(user);
    ownerS = await ethers.getSigner(owner);

    faucet = await ethers.getContract("Faucet", ownerS);
  });

  it("faucet", async () => {
    await replenish()
    await expect(faucet.faucet(user, 69, 420)).to
      .emit(faucet, "Faucet").withArgs(user, 69, 420)
      .changeEtherBalance(userS, 420);
  });

  it("faucet no money", async () => {
    await expect(faucet.faucet(user, 69, 420)).to
      .be.revertedWith("not enough funds");
  });

  it("faucet wrong acc", async () => {
    await expect(faucet.connect(userS).faucet(user, 69, 420)).to
      .be.revertedWith(`AccessControl: account ${user.toLowerCase()} is missing role ${ethers.constants.HashZero}`);
  });

  it("withdraw", async () => {
    await replenish()
    await expect(faucet.withdraw(user, 420)).to
      .changeEtherBalance(userS, 420);
  });

  it("withdraw wrong acc", async () => {
    await expect(faucet.connect(userS).withdraw(user, 420)).to
      .be.revertedWith(`AccessControl: account ${user.toLowerCase()} is missing role ${ethers.constants.HashZero}`);
  });

  it("withdraw no money", async () => {
    await expect(faucet.withdraw(user, 420)).to
      .be.revertedWith("not enough funds");
  });

  it("receive", async () => {
    expect(await faucet.provider.getBalance(faucet.address)).to.eq(0);
    await replenish()
    expect(await faucet.provider.getBalance(faucet.address)).to.eq(420);
  });

  async function replenish() {
    await ownerS.sendTransaction({to: faucet.address, value: 420, gasPrice: 0})
  }

});
