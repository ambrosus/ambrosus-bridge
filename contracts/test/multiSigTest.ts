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

        proxy = await ethers.getContract("proxyMultiSig", ownerS);
        implementation = await ethers.getContract("ProxyMultisigTest", ownerS);
    });

    beforeEach(async () => {
        await deployments.fixture(["for_tests"]);
    });

    it("Check Proxy", async () => {
        const lastProcessedBlock = "0x1111111111111111111111111111111111111111111111111111111111111111";

        const callData = implementation.interface.encodeFunctionData(
            implementation.interface.functions["updateLastProcessedBlock(bytes32)"],
            [lastProcessedBlock]
        );
        await proxy.upgradeToAndCall(implementation.address, callData);

        const tx = await proxy.connect(proxyAdminS).confirmTransaction(0);
        const receipt = await tx.wait();
        expect((await getEvents(receipt, "Upgraded"))[0].args[0]).eq(implementation.address);

        // Cannot use implementation.METHOD()
        const Factory = await ethers.getContractFactory("ProxyMultisigTest");
        const contract = await Factory.attach(proxy.address);

        expect(await contract.lastProcessedBlock()).eq(lastProcessedBlock);
    });

    it ("Non admin submitTransactionCall", async () => {
        await expect(proxy.connect(userS).submitTransaction(proxy.address, 0, "0x1234"))
            .to.be.revertedWith("");
    });
});

const getEvents = async (receipt: any, eventName: string) =>
    receipt.events?.filter((x: any) => x.event === eventName);
