import { deployments, ethers, getNamedAccounts } from "hardhat";
import { BigNumber, Contract, Signer } from "ethers";

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

    let multisig: Contract;


    before(async () => {
        ({ owner, user, admin } = await getNamedAccounts());
        ownerS = await ethers.getSigner(owner);
        adminS = await ethers.getSigner(admin);
        userS = await ethers.getSigner(user);

        multisig = await ethers.getContract("MultiSigWallet", ownerS);
    });

    beforeEach(async () => {
        await deployments.fixture(["for_tests"]);
    });


    it("makeTransactionSolo", async () => {
        await multisig.connect(ownerS).submitTransaction(user, 1, "0x");

        expect(await multisig.getTransactionCount(true, true)).to.eq(1);
        expect((await multisig.getTransactionIds(0, 1, true, true))[0]).to.eq(BigNumber.from(0));

        expect(await multisig.getConfirmationCount(0)).to.eq(1);
        expect(await multisig.getConfirmations(0)).to.eql([owner] );

        expect(await multisig.isConfirmed(0)).to.eq(true);

        // can't confirm twice
        await expect(multisig.connect(ownerS).confirmTransaction(0)).to.be.reverted;
        // can't execute twice
        // await expect(multisig.connect(ownerS).executeTransaction(0)).to.be.reverted;

    });

    it("makeTransactionDuo", async () => {
        // can't change requirement to 2 while there is only 1 owner
        await submitTx(multisig.populateTransaction.changeRequirement(2));
        expect((await multisig.transactions(latestTxId())).executed).to.be.false;
        expect(await multisig.required()).to.eq(1);

        // successfully change requirement to 2
        await submitTx(multisig.populateTransaction.addOwner(admin));
        await submitTx(multisig.populateTransaction.changeRequirement(2));
        expect(await multisig.required()).to.eq(2);


        // create tx
        await multisig.connect(ownerS).submitTransaction(user, 1, "0x");
        const txId = await latestTxId();

        expect(await multisig.getConfirmationCount(txId)).to.eq(1);
        expect(await multisig.isConfirmed(txId)).to.eq(false);

        // can't confirm twice
        await expect(multisig.connect(ownerS).confirmTransaction(txId)).to.be.reverted;

        // can't confirm without permissions
        await expect(multisig.connect(userS).confirmTransaction(txId)).to.be.reverted;

        await multisig.connect(ownerS).revokeConfirmation(txId);
        // can't revoke twice
        await expect(multisig.connect(userS).revokeConfirmation(txId)).to.be.reverted;


        await multisig.connect(adminS).confirmTransaction(txId);
        await multisig.connect(ownerS).confirmTransaction(txId);
    });


    it("changeOwners", async () => {
        await expect(multisig.addOwner(admin)).to.be.reverted;

        await submitTx(multisig.populateTransaction.addOwner(admin));
        expect(await multisig.getOwners()).to.eql([owner, admin]);

        // can't add twice
        await submitTx(multisig.populateTransaction.addOwner(admin));
        expect((await multisig.transactions(latestTxId())).executed).to.be.false;
        expect(await multisig.getOwners()).to.eql([owner, admin]);


        await submitTx(multisig.populateTransaction.replaceOwner(admin, user));
        expect(await multisig.getOwners()).to.eql([owner, user]);



        await submitTx(multisig.populateTransaction.addOwner(admin));
        expect(await multisig.getOwners()).to.eql([owner, user, admin]);

        await submitTx(multisig.populateTransaction.removeOwner(user));
        expect(await multisig.getOwners()).to.eql([owner, admin]);

        // can't remove twice
        await submitTx(multisig.populateTransaction.removeOwner(user))
        expect((await multisig.transactions(latestTxId())).executed).to.be.false;
        expect(await multisig.getOwners()).to.eql([owner, admin]);



        await submitTx(multisig.populateTransaction.removeOwner(admin));
        expect(await multisig.getOwners()).to.eql([owner]);
    });


    async function submitTx(populateTx: any, value=0) {
        return multisig.connect(ownerS).submitTransaction(multisig.address, value, (await populateTx).data);
    }

    async function latestTxId() {
        return (await multisig.transactionCount()) - 1;
    }

});
