// map new config structure to old one

import {readConfig, Token} from "../deploy/utils/config";
import path from "path";
import fs from "fs";

interface LegacyToken {
  isActive: boolean;
  name: string;
  symbol: string;
  logo: string;
  denomination: number;
  addresses: { [net: string]: string }
  primaryNets: string[];
  nativeAnalog: string | null;
  decimals: { [net: string]: number }
}


interface LegacyConfig {
  tokens: { [symb: string]: LegacyToken };
  bridges: { [net: string]: { amb: string, side: string } };
  ambFaucetAddress: string;
}


function mapToLegacyConfig(stage: string) {
  const newConfig = readConfig(stage);

  const legacyConfig: LegacyConfig = {
    bridges: newConfig.bridges,
    ambFaucetAddress: newConfig.ambFaucetAddress,
    tokens: {}
  }

  for (const [tokenName, newToken] of Object.entries(newConfig.tokens)) {
    legacyConfig.tokens[tokenName] = mapToLegacyToken(newToken);
  }

  const savePath = path.resolve(__dirname, `../configs/legacy/${stage}.json`)
  fs.writeFileSync(savePath, JSON.stringify(legacyConfig, null, 2))
}

function mapToLegacyToken(newToken: Token) {
  const legacyToken: LegacyToken = {
    isActive: newToken.isActive,
    name: newToken.name,
    symbol: newToken.symbol,
    denomination: 0,
    decimals: {},
    logo: newToken.logo,
    primaryNets: [],
    addresses: {},
    nativeAnalog: ""
  }

  for (const [netName, net] of Object.entries(newToken.networks)) {
    legacyToken.denomination = net.denomination;
    legacyToken.decimals[netName] = net.denomination;
    legacyToken.addresses[netName] = net.address;
    if (net.isPrimary) legacyToken.primaryNets.push(netName);
    if (net.nativeCoin) legacyToken.nativeAnalog = net.nativeCoin;
  }

  return legacyToken;
}

mapToLegacyConfig("test")
