import {network, deployments, ethers} from "hardhat";
import {Event} from "ethers";

const DEFAULT_ROLE = ethers.constants.HashZero
const ADMIN_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("ADMIN_ROLE"));
const RELAY_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("RELAY_ROLE"));
const WATCHDOG_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("WATCHDOG_ROLE"));
const FEE_PROVIDER_ROLE = ethers.utils.keccak256(ethers.utils.toUtf8Bytes("FEE_PROVIDER_ROLE"));

const roleNames = {
  [DEFAULT_ROLE]: "DEFAULT_ROLE",
  [ADMIN_ROLE]: "ADMIN_ROLE",
  [RELAY_ROLE]: "RELAY_ROLE",
  [WATCHDOG_ROLE]: "WATCHDOG_ROLE",
  [FEE_PROVIDER_ROLE]: "FEE_PROVIDER_ROLE",
}

async function main() {
  const BRIDGE_NAME = "ETH_EthBridge";

  const bridge = await ethers.getContract(BRIDGE_NAME);

  const deployBlock = (await deployments.get(BRIDGE_NAME)).receipt!.blockNumber;
  const logsLimit = (network.tags["bsc"]) ? 49_999 : undefined; // logs limit for bsc

  const events = await getEvents(bridge, bridge.filters.RoleGranted(), deployBlock, logsLimit);

  const possibleRoles: { [role: string]: string[] } = {};
  events.forEach(event => {
    const role = event.args!.role;

    if (!possibleRoles[role]) possibleRoles[role] = [];
    possibleRoles[role].push(event.args!.account)
  })

  const roles: { [role: string]: string[] } = {};
  for (const [role, accounts] of Object.entries(possibleRoles)) {
    const realAccounts = []
    for (const acc of accounts) {
      if (await bridge.hasRole(role, acc))
        realAccounts.push(acc)
      await sleep(200);
    }

    roles[roleNames[role]] = realAccounts;
  }

  console.log(roles)


}


async function getEvents(contract: any, filter: any, firstBlock: number, logsLimit?: number): Promise<Event[]> {
  const finishBlock = await contract.provider.getBlockNumber();

  const blockRange = [];
  if (logsLimit === undefined) {
    blockRange.push([firstBlock, finishBlock]); // from first to latest block
  } else {

    for (let fromBlock = firstBlock; fromBlock < finishBlock; fromBlock += logsLimit) {
      const toBlock = Math.min(fromBlock + logsLimit, finishBlock);
      blockRange.push([fromBlock, toBlock]);
    }
  }

  const events = [];
  for (const [fromBlock, toBlock] of blockRange) {
    // console.debug("fetching logs", fromBlock, "-", toBlock);
    const res = await withRetries(5, 1000, () => contract.queryFilter(filter, fromBlock, toBlock));
    await sleep(250);
    events.push(...res);
  }

  return events;
}


export async function withRetries(maxAttempts: number, delay: number, func: any) {
  for (let i = 0; i < maxAttempts; i++) {
    try {
      return await func();
    } catch (e) {
      console.error(e);
      await sleep(delay);
    }
  }
  throw new Error("max attempts reached");
}

export function sleep(ms: number) {
  return new Promise((resolve) => {
    setTimeout(resolve, ms);
  });
}


main().then(() => process.exit(0)).catch(error => {
  console.error(error);
  process.exit(1);
});
