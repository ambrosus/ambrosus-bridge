## Bridges

Pair (amb and side) of bridge contract addresses for each network

## Tokens

List of tokens used by bridge  

Tokens with empty addresses will be deployed when calling the deployment script.   
**Be careful**, this can lead to problems

#### Our token: _`SAMB`_
- **amb** will be deployed as `SAMB` contract
- **eth**, **bsc** will be deployed as `BridgeERC20` contract

#### Wrapped native coins on other networks: _`WBNB`, `WETH`_
- **amb** will be deployed as `BridgeERC20` contract
- **eth**, **bsc** addresses **must** be found on etherscan/bscscan

#### Other tokens: _`WBTC`, `WDOGE`_
- **amb** will be deployed as `BridgeERC20` contract
- **eth**, **bsc** addresses **must** be found on etherscan/bscscan
