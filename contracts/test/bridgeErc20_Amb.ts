import {deployments, ethers, getNamedAccounts, getUnnamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";
import { parseUnits } from "ethers/lib/utils";

chai.should();
export const expect = chai.expect;

describe("BridgeERC20_Amb", () => {
  let ownerS: Signer;
  let userS: Signer;
  let bridge6S: Signer;
  let owner: string;
  let user: string;
  let bridge6: string;

  let ambERC20: Contract;

  before(async () => {
    await deployments.fixture(["for_tests"]);
    ({owner, user, bridge: bridge6} = await getNamedAccounts());
    ownerS = await ethers.getSigner(owner);
    userS = await ethers.getSigner(user);
    bridge6S = await ethers.getSigner(bridge6);

    ambERC20 = await ethers.getContract("BridgeERC20_AmbTest", ownerS);
  });

  beforeEach(async () => {
    await deployments.fixture(["for_tests"]); // reset contracts state

    // enable bridge with DECIMALS == 6
    await ambERC20.setSideTokenDecimals([bridge6], [6]);
  });

  [6, 18, 20].forEach((bridgeDecimals, i) => {
    let bridgeS: Signer;
    let bridge: string;

    beforeEach(async () => {
      // await deployments.fixture(["for_tests"]); // reset contracts state

      bridge = (await getUnnamedAccounts())[i+10];  // different addresses with offset to eliminate collisions
      bridgeS = await ethers.getSigner(bridge);
      await ambERC20.setSideTokenDecimals([bridge], [bridgeDecimals]);
      expect(await ambERC20.sideTokenDecimals(bridge)).to.eq(bridgeDecimals);

    });


    describe(`bridge with decimal == ${bridgeDecimals}`, async () => {

      describe("bridge is sender", () => {

        it("should mint (transfer)", async () => {
          // transfer 1 TOKEN from 1e6 BRIDGE to 1e18 USER
          // AMOUNT is 1e6 == 1 TOKEN in bridge network
          // USER MUST RECEIVE 1e18 == 1 TOKEN in user network
          await ambERC20.connect(bridgeS).transfer(owner, parseUnits("1", bridgeDecimals));
          expect(await ambERC20.balanceOf(owner)).eq(parseUnits("1", 18));
          expect(await ambERC20.balanceOf(bridge)).eq(0);
        });

        // `transferFrom` doesn't work with bridge as sender,
        // coz bridge should call `increaseAllowance` for it.
        // bridge contract use `transfer` for sending.
        // but we test it anyway


        it("should mint (transferFrom)", async () => {
          // transfer 1 TOKEN from 1e6 BRIDGE to 1e18 USER
          // AMOUNT is 1e6 == 1 TOKEN in bridge network
          // USER MUST RECEIVE 1e18 == 1 TOKEN in user network

          await ambERC20.connect(bridgeS).increaseAllowance(owner, parseUnits("1", 18));
          await ambERC20.connect(ownerS).transferFrom(bridge, owner, parseUnits("1", bridgeDecimals));
          expect(await ambERC20.balanceOf(owner)).eq(parseUnits("1", 18));
          expect(await ambERC20.balanceOf(bridge)).eq(0);
        });

      });


      describe("bridge is recipient", () => {

        it("should throw error if sender has insufficient bridgeBalance amount", async () => {
          await expect(ambERC20.transfer(bridge, 100000))
            .to.be.revertedWith("not enough locked tokens on bridge");
        });

        it("should burn (transfer)", async () => {
          // transfer 2 TOKEN from 1e18 USER to 1e6 BRIDGE
          // AMOUNT is 2e6 == 2 TOKEN in bridge network
          // USER MUST LOSE 2e18 == 2 TOKEN in user network

          await bridgeMint(bridgeS, owner, parseUnits("3", bridgeDecimals))

          await ambERC20.transfer(bridge, parseUnits("2", bridgeDecimals));
          expect(await ambERC20.balanceOf(owner)).eq(parseUnits("1", 18));
          expect(await ambERC20.balanceOf(bridge)).eq(parseUnits("0", 18));
        });

        it("should burn (transferFrom)", async () => {
          // transfer 2 TOKEN from 1e18 USER to 1e6 BRIDGE
          // AMOUNT is 2e6 == 2 TOKEN in bridge network
          // USER MUST LOSE 2e18 == 2 TOKEN in user network

          await bridgeMint(bridgeS, owner, parseUnits("3", bridgeDecimals))

          await ambERC20.increaseAllowance(bridge, parseUnits("2", 18));
          await ambERC20.connect(bridgeS).transferFrom(owner, bridge, parseUnits("2", bridgeDecimals));
          expect(await ambERC20.balanceOf(owner)).eq(parseUnits("1", 18));
          expect(await ambERC20.balanceOf(bridge)).eq(parseUnits("0", 18));
        });


      });

    });

  });


  describe("sender and recipient is not a bridge", () => {

    it("should just transfer (transfer)", async () => {
      // mint 3 TOKENS
      await bridgeMint(bridge6S, owner, parseUnits("3", 6))

      // transfer 2 TOKENS
      await ambERC20.transfer(user, parseUnits("2", 18));
      expect(await ambERC20.balanceOf(user)).eq(parseUnits("2", 18));
      expect(await ambERC20.balanceOf(owner)).eq(parseUnits("1", 18));
    });

    it("should just transfer (transferFrom called by owner)", async () => {
      // mint 3 TOKENS
      await bridgeMint(bridge6S, owner, parseUnits("3", 6))

      // transfer 2 TOKENS
      await ambERC20.increaseAllowance(owner, parseUnits("2", 18));
      await ambERC20.connect(ownerS).transferFrom(owner, user, parseUnits("2", 18));

      expect(await ambERC20.balanceOf(user)).eq(parseUnits("2", 18));
      expect(await ambERC20.balanceOf(owner)).eq(parseUnits("1", 18));
    });

    it("should just transfer (transferFrom called by bridge)", async () => {
      // mint 3 TOKENS
      await bridgeMint(bridge6S, owner, parseUnits("3", 6))

      // transfer 2 TOKENS
      await ambERC20.increaseAllowance(bridge6, parseUnits("2", 18));
      await ambERC20.connect(bridge6S).transferFrom(owner, user, parseUnits("2", 18));

      expect(await ambERC20.balanceOf(user)).eq(parseUnits("2", 18));
      expect(await ambERC20.balanceOf(owner)).eq(parseUnits("1", 18));
    });

  });

  const bridgeMint = (bridgeS: Signer, to: string, amount: any) => ambERC20.connect(bridgeS).transfer(to, amount);

});
