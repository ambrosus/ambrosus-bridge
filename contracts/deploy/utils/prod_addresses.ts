import path from "path";


interface constructorConfig {
  adminAddress: string,
  relayAddress: string,
  transferFeeRecipient: string,
  bridgeFeeRecipient: string,
  multisig: {
    admins: string[],
    threshold: number
  }
}


export function getAddresses(network: string): constructorConfig {
  return readConfig()[network];
}

function readConfig(): { [net: string]: constructorConfig } {
  return require(path.resolve(__dirname, `../../configs/prod_addresses.json`));
}
