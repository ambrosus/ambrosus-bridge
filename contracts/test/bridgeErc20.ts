import { deployments, ethers, getNamedAccounts } from "hardhat";
import type { Contract, Signer } from "ethers";

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
    ({ owner, user, bridge } = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    userS = await ethers.getSigner(user);
    bridgeS = await ethers.getSigner(bridge);

    mockERC20 = await ethers.getContract("BridgeERC20Test", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state
  });

  describe("transfer/transferFrom", () => {
    beforeEach(async () => {
      await mockERC20.setBridgeAddressesRole([bridge]);
      await mockERC20.mint(owner, 1000);
      await mockERC20.increaseAllowance(owner, 1000);
      await mockERC20.connect(bridgeS).increaseAllowance(bridge, 1000);
    });

    it("should mint if sender is bridge (transfer)", async () => {
      await mockERC20.connect(bridgeS).transfer(owner, 1);
      expect(await mockERC20.balanceOf(owner)).eq(1001);
      expect(await mockERC20.balanceOf(bridge)).eq(0);
    });

    it("should mint if sender is bridge (transferFrom)", async () => {
      await mockERC20.connect(bridgeS).transferFrom(bridge, owner, 1);
      expect(await mockERC20.balanceOf(owner)).eq(1001);
      expect(await mockERC20.balanceOf(bridge)).eq(0);
    });

    it("should burn if recipient is bridge (transfer)", async () => {
      await mockERC20.transfer(bridge, 1);
      expect(await mockERC20.balanceOf(owner)).eq(999);
      expect(await mockERC20.balanceOf(bridge)).eq(0);
    });

    it("should burn if recipient is bridge (transferFrom)", async () => {
      await mockERC20.transferFrom(owner, bridge, 1);
      expect(await mockERC20.balanceOf(owner)).eq(999);
      expect(await mockERC20.balanceOf(bridge)).eq(0);
    });

    it("should simple transfer if sender or recipient isn't bridge (transfer)", async () => {
      await mockERC20.transfer(user, 1);
      expect(await mockERC20.balanceOf(owner)).eq(999);
      expect(await mockERC20.balanceOf(user)).eq(1);
    });

    it("should simple transfer if sender or recipient isn't bridge (transferFrom)", async () => {
      await mockERC20.transferFrom(owner, user, 1);
      expect(await mockERC20.balanceOf(owner)).eq(999);
      expect(await mockERC20.balanceOf(user)).eq(1);
    });
  });
});
