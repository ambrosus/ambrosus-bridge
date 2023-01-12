import prod_addresses from '../../configs/prod_addresses';

interface constructorConfig {
  adminAddress: string,
  relayAddress: string,
  masterRelayAddress: string,
  feeProviderAddress: string,
  watchdogsAddresses: string[],
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

export function readConfig(): { [net: string]: constructorConfig } {
  return prod_addresses;
}
