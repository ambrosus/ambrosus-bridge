import * as hre from "hardhat";
import { ethers } from "hardhat";
import { BytesLike } from "@ethersproject/bytes";
import { NonceManager } from "@ethersproject/experimental";
import { Contract } from "ethers";

import * as path from "path";
import * as fs from "fs";

const tokensFilepath = "./tokens.json";
const tokensFile = require(tokensFilepath);
const contractName = "ERC20Token";

async function main() {
  let tokens: any;

  const hreConfig: any = hre.network.config;
  const signer = getSigner(hreConfig.url, hreConfig.accounts[0]);

  // setting network specific vars
  switch (hre.network.name) {
    case "amb":
      tokens = tokensFile.amb.tokens;
      break;
    case "rinkeby": case "mainnet":
      tokens = tokensFile.eth.tokens;
      break;
    default:
      return console.log("Unknown network: %s", hre.network.name);
  }

  // main
  for (const token of tokens) {
    if (token.address) continue;

    console.log("Deploying ERC20 token (%s, %s) into %s network", token.name, token.symbol, hre.network.name);
    const deployedToken = await deployToken(token, signer);

    token.address = deployedToken.address;
    fs.writeFileSync(path.resolve(__dirname, tokensFilepath), JSON.stringify(tokensFile, null, 2));
  }
}

async function deployToken(token: any, network: any): Promise<Contract> {
  const deployedToken = await deploy(contractName, network, token.name, token.symbol, token.bridgeAddresses);
  await deployedToken.deployed();
  return deployedToken;
}

async function deploy(contractName: string, network: any, ...args: Array<any>) {
  const factory = await ethers.getContractFactory(contractName, network);
  return await factory.deploy(...args);
}

function getSigner(url: string, pk: BytesLike) {
  const provider = new ethers.providers.JsonRpcProvider(url);
  const signer = new ethers.Wallet(pk, provider);
  return new NonceManager(signer);
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });
