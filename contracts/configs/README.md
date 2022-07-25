## Bridges

Pair (amb and side) of bridge contract addresses for each network

## Tokens

List of tokens used by bridge  

**Only** tokens with addresses == `DEPLOY` will be deployed when calling the deployment script.

#### Our token: _`SAMB`_
- **amb** will be deployed as `SAMB` contract
- **eth**, **bsc** will be deployed as `BridgeERC20` contract

#### Wrapped native coins on other networks: _`WBNB`, `WETH`_
- **amb** will be deployed as `BridgeERC20` contract
- **eth**, **bsc** addresses **must** be found on etherscan/bscscan

#### Other tokens: _`WBTC`, `WDOGE`_
- **amb** will be deployed as `BridgeERC20` contract
- **eth**, **bsc** addresses **must** be found on etherscan/bscscan


## Prod Addresses
Values from `prod_addressed.json` file will be used for bridge deploy.  

- **adminAddress** - can configure bridge after deploy
- **relayAddress** - collects information from another network and sends according to the contract; must be same address that relay service use
- **transferFeeRecipient** - receives fee for covering expenses
- **bridgeFeeRecipient** - receives bridge profit
- **multisig.admins** - their signatures are required to update the contract versions
- **multisig.threshold** - how many admin votes required to apply transaction 

Each bridge contract have separate object in cfg.
Bridge names described [there](https://github.com/ambrosus/ambrosus-bridge#contracts) and in `description` field of the cfg file.
