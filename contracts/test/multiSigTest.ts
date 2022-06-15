import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;


describe("MultiSig test", () => {
    let ownerS: Signer;
    let owner: string;
    let adminS: Signer;
    let admin: string;
    let user: string;
    let userS: Signer;

    let proxy: Contract;
    let implementation: Contract;
    let mockERC20: Contract;


    before(async () => {
        ({owner, user, admin} = await getNamedAccounts());
        ownerS = await ethers.getSigner(owner);
        adminS = await ethers.getSigner(admin);
        userS = await ethers.getSigner(user);

        proxy = await ethers.getContract("ProxyMultiSig", ownerS);
        implementation = await ethers.getContract("ProxyMultisigTest", ownerS);
        mockERC20 = await ethers.getContract("BridgeERC20Test", ownerS);
    });

    beforeEach(async () => {
        await deployments.fixture(["for_tests"]);
    });

    it("upgradeTo", async () => {
        await proxy.upgradeTo(implementation.address);

        await expect(proxy.connect(adminS).confirmTransaction(0)).
            to.emit(proxy, "Upgraded")
            .withArgs(implementation.address);

        const callData = implementation.interface.encodeFunctionData(
            implementation.interface.functions["changeValue(bytes4)"], ["0x11223344"]
        );

        await proxy.submitTransaction(
            proxy.address,
            0,
            callData
        );

        await expect(proxy.connect(adminS).confirmTransaction(1))
            .to.emit(proxy, "Execution")
            .withArgs(1);

        expect(await implementation.attach(proxy.address).value()).eq("0x11223344");

    });

    it ("UpgradeToAndCall", async () => {
        const callData = implementation.interface.encodeFunctionData(
            implementation.interface.functions["changeValue(bytes4)"], ["0x44332211"]
        );

        await proxy.upgradeToAndCall(implementation.address, callData);

        await expect(proxy.connect(adminS).confirmTransaction(0)).
        to.emit(proxy, "Upgraded")
            .withArgs(implementation.address);

        expect(await implementation.attach(proxy.address).value()).eq("0x44332211");

    });

    it ("MultiSig submitTransaction", async () => {
        await mockERC20.mint(proxy.address, 1000);

        const callData = mockERC20.interface.encodeFunctionData(
            mockERC20.interface.functions["transfer(address,uint256)"], [admin, 500]
        );

        await proxy.submitTransaction(
            mockERC20.address,
            0,
            callData
        );

        await expect(() => proxy.connect(adminS).confirmTransaction(0))
            .to.changeTokenBalance(mockERC20, adminS, 500);
    });

    it ("Non admin submitTransaction Call", async () => {
        await expect(proxy.connect(userS).submitTransaction(proxy.address, 0, "0x1234"))
            .to.be.revertedWith("");
    });
});
