// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";


// this contract is already deployed on mainnets as USDC and BUSD tokens, but
// it has an unpleasant feature: allowance need to be set with ANOTHER token denomination
// and user will see counterintuitive amount in increaseAllowance() function.
// SO, FRONT SHOULD CHECK FOR USDC AND BUSD TOKENS AND USE THIS LEGACY LOGIC FOR THEM
contract BridgeERC20_Amb_OLD is ERC20, Ownable {

    // decimals of token in side network
    // example:
    // 0xBSC_Amb => 18 (AMB contract of BSC bridge  => 18 decimals)
    // 0xETH_Amb => 6  (AMB contract of ETH bridge  => 6 decimals)
    // now, token will auto convert self _decimals to side _decimals (or vice versa) on bridge transfer
    // NOTE: value 0 means that address is not a bridge; DON'T SET NON ZERO VALUES FOR NON BRIDGE ADDRESSES
    mapping(address => uint8) public sideTokenDecimals;

    mapping(address => uint) public bridgeBalances;  // locked tokens on the side bridge

    uint8 _decimals;

    constructor(
        string memory name_, string memory symbol_, uint8 decimals_,
        address[] memory bridgeAddresses_, uint8[] memory sideTokenDecimals_
    ) ERC20(name_, symbol_) Ownable() {
        _setSideTokenDecimals(bridgeAddresses_, sideTokenDecimals_);
        _decimals = decimals_;
    }

    function decimals() public view override returns (uint8) {
        return _decimals;
    }

    function setSideTokenDecimals(address[] memory bridgeAddresses_, uint8[] memory sideTokenDecimals_) public onlyOwner() {
        _setSideTokenDecimals(bridgeAddresses_, sideTokenDecimals_);
    }

    // todo check if we need this func
    function _setSideTokenDecimals(address[] memory bridgeAddresses_, uint8[] memory sideTokenDecimals_) private {
        require(bridgeAddresses_.length == sideTokenDecimals_.length, "wrong array lengths");
        for (uint i = 0; i < bridgeAddresses_.length; i++)
            sideTokenDecimals[bridgeAddresses_[i]] = sideTokenDecimals_[i];
    }

    function _transfer(
        address sender,
        address recipient,
        uint amount
    ) internal virtual override {
        // todo events
        if (sideTokenDecimals[sender] != 0) { // sender is bridge
             // user transfer tokens to ambrosus => need to mint it

            // we receive tokens from network where token have sideTokenDecimals[sender] decimals
            // convert amount with SIDE network decimals form to SELF decimals form
            uint amount_this = _convertDecimals(amount, sideTokenDecimals[sender], _decimals);


            // bridge mint money to user; same amount locked on side bridge
            bridgeBalances[sender] += amount_this;

            _mint(recipient, amount_this);
        } else if (sideTokenDecimals[recipient] != 0) { // recipient is bridge
            // user withdraw tokens from ambrosus => need to burn it

            // we transfer tokens to network where token have sideTokenDecimals[sender] decimals
            // convert amount with SIDE network decimals form to SELF decimals form
            uint amount_this = _convertDecimals(amount, sideTokenDecimals[recipient], _decimals);


            // user burn tokens; side bridge must have enough tokens to send
            require(bridgeBalances[recipient] >= amount_this, "not enough locked tokens on bridge");
            bridgeBalances[recipient] -= amount_this;

            _burn(sender, amount_this);
        } else {
            super._transfer(sender, recipient, amount);
        }
    }

    function _convertDecimals(uint256 amount, uint8 dFrom, uint8 dTo) internal pure returns (uint256) {
        if (dTo == dFrom)
            return amount;
        if (dTo > dFrom)
            return amount * (10 ** (dTo - dFrom));
        else
            return amount / (10 ** (dFrom - dTo));
    }

}
