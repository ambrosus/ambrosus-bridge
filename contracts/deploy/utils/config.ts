import path from "path";
import fs from "fs";
import {ethers} from "ethers";
import { isAddress } from "ethers/lib/utils";




export interface Token {
  isActive: boolean;
  name: string;
  symbol: string;
  denomination: number;
  addresses: { [net: string]: string }
  primaryNets: string[];
  nativeAnalog: string | null;
  decimals: { [net: string]: number }
}


export interface Config {
  tokens: { [symb: string]: Token };
  bridges: { [net: string]: { amb: string, side: string } };

  save(): void;
  getTokenPairs(thisNet: string, sideNet: string): { [k: string]: string }
}



export function readConfig(stage: string): Config {
  const tokenPath = path.resolve(__dirname, `../../configs/${stage}.json`)
  const config = require(tokenPath);

  config.save = () => fs.writeFileSync(tokenPath, JSON.stringify(config, null, 2));
  config.getTokenPairs = (thisNet: string, sideNet: string) => getTokenPairs(thisNet, sideNet, config)

  return config;
}

export function isTokenWrappable(token: Token): boolean {
    return token.nativeAnalog != ""
}

export function isTokenPrimary(token: Token, netName: string): boolean {
    return token.primaryNets.includes(netName)
}


function getTokenPairs(thisNet: string, sideNet: string, configFile: Config): { [k: string]: string } {
  const tokenPair: { [k: string]: string } = {};

  for (const token of Object.values(configFile.tokens)) {
    if (!token.isActive) continue;
    // token must be deployed on both sides
    if (!isAddress(token.addresses[thisNet]) || !isAddress(token.addresses[sideNet]))
      continue

    tokenPair[token.addresses[thisNet]] = token.addresses[sideNet];


    // if token is wrapper for side-network native coin - save it address with 0x0 key,
    // so if user wants to transfer synthetic wrapper and receive native coin - he can do it
    // by calling withdraw(tokenThisAddress=thisNet, unwrapSide=true, ...).
    // bridge will check if tokens[0x0] == tokenThisAddress and emit Transfer with zero token address
    if (isTokenPrimary(token, sideNet) && isTokenWrappable(token))
      tokenPair[ethers.constants.AddressZero] = token.addresses[thisNet];


  }

  return tokenPair;
}
