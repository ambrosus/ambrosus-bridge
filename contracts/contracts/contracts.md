## Contracts

### Networks
**Entry point for bridge contracts**  
Contains a pair of 2 contracts for each bridge.  
These contracts are inherited from `common/CommonBridge.sol` and some check from `checks` (currently all contracts are `CheckUntrustless2`)  

### Common
Common structs and code for bridges

### Checks
Functions that validate information submitted in bridge.
- `CheckUntrustless` - is like onchain multisig for relays - relays push their data into `checkUntrustless_` method and 
when all relays signs are collected - this function return `true` and the data is considered valid. 
Currently not used in favor of `CheckUntrustless2` 
- `CheckUntrustless2` - just empty contract. There is no checks, coz relays must create MPC signature
which will indicate their consensus and will be accepted by the contract (coz it address has RELAY role)
- `SignatureCheck` - utility function to check ECDSA sign.

### Tokens
- `IWrapper` - interface for wrapper contracts for native coins (like WETH in ethereum network and SAMB in ambrosus)
- `sAMB` - ERC20 wrapper for native amb coin
- `BridgeERC20` - contract for **synthetic tokens** (i.e. bridge can mint and burn them) for network != AMB
- `BridgeERC20_Amb` - contract for **synthetic tokens** (i.e. bridge can mint and burn them) for AMB network.  
It different from the `BridgeERC20` in that **it can have different decimals** depending on which bridge 
(ex: ETH or BSC) are using for transferring him.    
For example, if USDT has 6 decimals in ETH network and 18 decimals in BSC and AMB networks, it will works like this:
  - `ETH -> AMB: 1 USDT (1e6 ) -> 1 USDT (1e18)`
  - `AMB -> BSC: 1 USDT (1e18) -> 1 USDT (1e18)`
  - `ETH -> ETH: 1 USDT (1e6 ) -> 1 USDT (1e6 )`
- `BridgeERC20_Amb_OLD` - like previous one, but has unpleasant feature: allowance need to be set with   
ANOTHER token denomination and user will see counterintuitive amount in increaseAllowance() function

## Contracts_For_Test
Wrappers for contracts with some methods has their public `Test`-suffixed version.
- `MintableERC20` - deployed in testnets instead of primary (i.e. not synthetics) tokens. Has `mint()` method
