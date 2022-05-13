import {deployments, ethers, getNamedAccounts} from "hardhat";
import type {Contract, Signer} from "ethers";

import chai from "chai";

chai.should();
export const expect = chai.expect;


describe("MultiSig test", () => {
    let ownerS: Signer;
    let owner: string;
    let proxyAdminS: Signer;
    let proxyAdmin: string;
    let user: string;
    let userS: Signer;

    let proxy: Contract;
    let implementation: Contract;


    before(async () => {
        ({owner, user, proxyAdmin} = await getNamedAccounts());
        ownerS = await ethers.getSigner(owner);
        proxyAdminS = await ethers.getSigner(proxyAdmin);
        userS = await ethers.getSigner(user);

        proxy = await ethers.getContract("ProxyMultiSig", ownerS);
        implementation = await ethers.getContract("ProxyMultisigTest", ownerS);
    });

    beforeEach(async () => {
        await deployments.fixture(["for_tests"]);
    });

    it("Check Proxy", async () => {
        await proxy.upgradeTo(implementation.address);

        await expect(proxy.connect(proxyAdminS).confirmTransaction(0)).
            to.emit(proxy, "Upgraded")
            .withArgs(implementation.address);

        const callData = implementation.interface.encodeFunctionData(
            implementation.interface.functions["changeValue(bytes4)"], ["0x11223344"]
        );

        await proxy.submitTransaction(
            implementation.address,
            0,
            callData
        );

        await expect(proxy.connect(proxyAdminS).confirmTransaction(1))
            .to.emit(proxy, "Execution")
            .withArgs(1);

        expect(await implementation.value()).eq("0x11223344");

    });

    it ("Non admin submitTransactionCall", async () => {
        await expect(proxy.connect(userS).submitTransaction(proxy.address, 0, "0x1234"))
            .to.be.revertedWith("");
    });
});
