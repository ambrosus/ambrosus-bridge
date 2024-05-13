import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

describe("MintableERC20", () => {
  let userS: Signer;
  let user: string;
  let mintableErc20: Contract;

  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({user} = await getNamedAccounts());
    userS = await ethers.getSigner(user);
    mintableErc20 = await ethers.getContract("MintableERC20");
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state
  });

  it("should mint", async () => {
    await mintableErc20.mint(user, 1);
    expect(await mintableErc20.balanceOf(user)).eq(1);
  });

  it("decimals", async () => {
    expect(await mintableErc20.decimals()).eq(18);
  });

});

