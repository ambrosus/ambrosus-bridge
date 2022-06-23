import inquirer from "inquirer";
import {readConfig} from "../deploy/utils/config";
import {execSync} from "child_process";


// todo checkbox choose
// todo redeploy

const NETWORKS = ["eth", "bsc"]

async function main() {
  await chooseStage()
}

async function chooseStage() {
  const stage = await choose(["dev", "test", "main", "integr"], "Choose stage:");
  await chooseAction(stage);
}

async function chooseAction(stage: string) {
  const actions = [
    "full deploy all networks",
    ...NETWORKS.map((n) => ({name: `full deploy ${n}-amb`, value: n})),
    "tokens", "bridges",
  ]
  const action = await choose(actions, "Choose action:");

  if (action == "tokens") await tokens(stage);
  else if (action == "bridges") await bridges(stage);
  else if (action == "full deploy all networks") full_deploy_all_networks(stage);
  else full_deploy(stage, action)  // action == network name (eth/bsc/...)
}


async function tokens(stage: string) {
  const {tokens} = readConfig(stage);
  for (const token of Object.values(tokens)) {
    const deployedOn = Object.entries(token.addresses)
      .map(([network, address]) => `\t ${network}  ${address || "Not deployed"}`)
      .join("\n")
    console.log(`${token.name} ${token.symbol} \n ${deployedOn}`);
  }

  const actions = [
    ...NETWORKS.map((n) => ({name: "deploy " + n, value: n})),
    new inquirer.Separator(),
    {name: "deploy all", value: "all"}
  ];

  const network = await choose(actions, "Choose action:");
  if (network == "all") NETWORKS.forEach((n) => deploy_tokens(stage, n))
  else deploy_tokens(stage, network)

}


async function bridges(stage: string) {
  const {bridges} = readConfig(stage);

  for (const [pairName, addresses] of Object.entries(bridges))
    console.log(`amb-${pairName} 
\t amb  ${addresses.amb || "Not deployed"}
\t ${pairName}  ${addresses.side || "Not deployed"}`);


  const pairsToChoose = Object.keys(bridges).map((pairName) => ({name: `amb-${pairName}`, "value": pairName}))
  const pair = await choose(pairsToChoose, "Choose pair:");
  await bridge(stage, pair);

}


async function bridge(stage: string, pair: string) {
  const {bridges} = readConfig(stage);
  const bridge = bridges[pair];

  if (!bridge.amb || !bridge.side) {
    await choose(["deploy all"], "Choose action:");
    deploy_bridges(stage, pair);
    return
  }

  const action = await choose(["upgrade amb", "upgrade side"], "Choose action:");
  if (action == "upgrade amb") upgrade_bridge(stage, pair, "amb")
  if (action == "upgrade side") upgrade_bridge(stage, pair, pair)

}


// actions

function full_deploy_all_networks(stage: string) {
  for (const network of NETWORKS)
    full_deploy(stage, network);
}

function full_deploy(stage: string, network: string) {
    deploy_tokens(stage, network);
    deploy_bridges(stage, network);
}

function deploy_tokens(stage: string, network: string) {
  deploy(stage, "amb", "tokens");
  deploy(stage, network, "tokens");
}

function deploy_bridges(stage: string, network: string) {
  const bridgeTag = `bridges_${network}`

  deploy(stage, "amb", bridgeTag);
  deploy(stage, network, bridgeTag);

  // setSideBridge to newly deployment
  deploy(stage, "amb", bridgeTag);
  // add bridges addresses to deployed tokens
  add_bridges_to_tokens(stage, network);
}


function upgrade_bridge(stage: string, network: string, bridge: string) {
  deploy(stage, network, bridge);
  if (stage == "main") console.log("!! DON'T FORGET TO CONFIRM UPGRADE IN MULTISIG !!");
}


function add_bridges_to_tokens(stage: string, network: string) {
  deploy(stage, "amb", "tokens_add_bridges");
  deploy(stage, network, "tokens_add_bridges");
}


const deploy = (stage: string, network: string, tags: string) => exec(`yarn hardhat deploy --network ${stage}/${network} --tags ${tags}`);
const exec = (cmd: string) => execSync(cmd, {stdio: 'inherit'});


// dialog

const confirm = async (message: string): Promise<boolean> =>
  prompt({
    type: 'list',
    message: message,
    choices: [
      {name: "NO", value: false},
      {name: "YES", value: true}
    ]
  });

const choose = async (choices: Array<any>, message: string) => prompt({
  type: 'list',
  message: message,
  choices: choices
});

const prompt = async (q: any) => (await inquirer.prompt([{...q, name: 'value'}])).value;


// utils


main().catch(reason => {
  console.log(reason);
  process.exitCode = -1;
});
