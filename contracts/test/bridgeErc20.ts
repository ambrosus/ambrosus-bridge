import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;

describe("BridgeERC20", () => {
  let ownerS: Signer;
  let userS: Signer;
  let bridgeS: Signer;
  let owner: string;
  let user: string;
  let bridge: string;

  let mockERC20: Contract;

  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, user, bridge} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    userS = await ethers.getSigner(user);
    bridgeS = await ethers.getSigner(bridge);

    mockERC20 = await ethers.getContract("BridgeERC20Test", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state

    await mockERC20.setBridgeAddressesRole([bridge]);
    // await mockERC20.mint(owner, 1000);
    // await mockERC20.increaseAllowance(owner, 1000);
    // await mockERC20.connect(bridgeS).increaseAllowance(bridge, 1000);
  });

  describe("bridge is sender", () => {

    it("should mint (transfer)", async () => {
      await mockERC20.connect(bridgeS).transfer(owner, 1);
      expect(await mockERC20.balanceOf(owner)).eq(1);
      expect(await mockERC20.balanceOf(bridge)).eq(0);
    });

    // `transferFrom` doesn't work with bridge as sender,
    // coz bridge should call `increaseAllowance` for it.
    // bridge contract use `transfer` for sending, so it's ok

  });


  describe("bridge is recipient", () => {

    it("should throw error if sender has insufficient bridgeBalance amount", async () => {
      await expect(mockERC20.transfer(bridge, 100000))
        .to.be.revertedWith("not enough locked tokens on bridge");
    });

    it("should burn (transfer)", async () => {
      await bridgeMint(owner, 3)

      await mockERC20.transfer(bridge, 2);
      expect(await mockERC20.balanceOf(owner)).eq(1);
      expect(await mockERC20.balanceOf(bridge)).eq(0);
    });

    it("should burn (transferFrom)", async () => {
      await bridgeMint(owner, 3)

      await mockERC20.increaseAllowance(bridge, 2);
      await mockERC20.connect(bridgeS).transferFrom(owner, bridge, 2);
      expect(await mockERC20.balanceOf(owner)).eq(1);
      expect(await mockERC20.balanceOf(bridge)).eq(0);
    });



  });


  describe("sender and recipient is not a bridge", () => {

    it("should just transfer (transfer)", async () => {
      await bridgeMint(owner, 3)
      await mockERC20.transfer(user, 2);
      expect(await mockERC20.balanceOf(user)).eq(2);
      expect(await mockERC20.balanceOf(owner)).eq(1);
    });

    it("should just transfer (transferFrom called by owner)", async () => {
      await bridgeMint(owner, 3)

      await mockERC20.increaseAllowance(owner, 2);
      await mockERC20.transferFrom(owner, user, 2);

      expect(await mockERC20.balanceOf(user)).eq(2);
      expect(await mockERC20.balanceOf(owner)).eq(1);
    });

    it("should just transfer (transferFrom called by bridge)", async () => {
      await bridgeMint(owner, 3)

      await mockERC20.increaseAllowance(bridge, 2);
      await mockERC20.connect(bridgeS).transferFrom(owner, user, 2);

      expect(await mockERC20.balanceOf(user)).eq(2);
      expect(await mockERC20.balanceOf(owner)).eq(1);
    });

  });

  const bridgeMint = (to: string, amount: number) => mockERC20.connect(bridgeS).transfer(to, amount);

});

