import path from "path";
import fs from "fs";
import {ethers} from "ethers";
import { isAddress } from "ethers/lib/utils";




interface Token {
  name: string;
  symbol: string;
  denomination: number;
  addresses: { [net: string]: string }
  primaryNet: string;
  nativeAnalog: string | null;
}


export interface Config {
  tokens: { [symb: string]: Token };
  bridges: { [net: string]: { amb: string, side: string } };

  save(): void;
  getTokenPairs(thisNet: string, sideNet: string): { [k: string]: string }
  bridgesInNet(net: string): string[]
}



export function readConfig(stage: string): Config {
  const tokenPath = path.resolve(__dirname, `../../configs/${stage}.json`)
  const config = require(tokenPath);

  config.save = () => fs.writeFileSync(tokenPath, JSON.stringify(config, null, 2));
  config.getTokenPairs = (thisNet: string, sideNet: string) => getTokenPairs(thisNet, sideNet, config)
  config.bridgesInNet = (net: string) => bridgesInNet(net, config)

  return config;
}



function getTokenPairs(thisNet: string, sideNet: string, configFile: Config): { [k: string]: string } {
  const tokenPair: { [k: string]: string } = {};

  for (const token of Object.values(configFile.tokens)) {
    // token must be deployed on both sides
    if (!isAddress(token.addresses[thisNet]) || !isAddress(token.addresses[sideNet]))
      continue

    tokenPair[token.addresses[thisNet]] = token.addresses[sideNet];

    if (token.primaryNet === sideNet && token.nativeAnalog)   // native token for sideNet
      tokenPair[ethers.constants.AddressZero] = token.addresses[thisNet];

  }

  return tokenPair;
}


// get all deployed bridges in `net` network;
// for amb it's array of amb addresses for each network pair (such "amb-eth" or "amb-bsc")
// for other networks is array of one address
function bridgesInNet(net: string, configFile: Config): string[] {
  const bridges = (net == "amb") ?
    Object.values(configFile.bridges).map(i => i.amb) :
    [configFile.bridges[net].side];
  return bridges.filter(i => !!i);  // filter out empty strings
}
