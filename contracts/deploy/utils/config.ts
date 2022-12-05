import path from "path";
import fs from "fs";
import {ethers} from "ethers";
import { isAddress } from "ethers/lib/utils";


export interface Token {
  isActive: boolean;
  name: string;
  symbol: string;
  networks: { [net: string]: TokenNetwork }
}

export interface TokenNetwork {
  address: string;
  denomination: number;
  isPrimary: boolean;
  nativeCoin?: string;
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

export function isTokenWrappable(token: Token, netName: string): boolean {
    return !!token.networks[netName].nativeCoin
}

export function isTokenPrimary(token: Token, netName: string): boolean {
    return token.networks[netName].isPrimary
}


function getTokenPairs(thisNet: string, sideNet: string, configFile: Config): { [k: string]: string } {
  const tokenPair: { [k: string]: string } = {};

  for (const token of Object.values(configFile.tokens)) {
    if (!token.isActive) continue;
    // token must be deployed on both sides
    if (!isAddress(token.networks[thisNet]?.address) || !isAddress(token.networks[sideNet]?.address))
      continue

    tokenPair[token.networks[thisNet].address] = token.networks[sideNet].address;


    // if token is wrapper for side-network native coin - save it address with 0x0 key,
    // so if user wants to transfer synthetic wrapper and receive native coin - he can do it
    // by calling withdraw(tokenThisAddress=thisNet, unwrapSide=true, ...).
    // bridge will check if tokens[0x0] == tokenThisAddress and emit Transfer with zero token address
    if (isTokenPrimary(token, sideNet) && isTokenWrappable(token, sideNet))
      tokenPair[ethers.constants.AddressZero] = token.networks[thisNet].address;


  }

  return tokenPair;
}
