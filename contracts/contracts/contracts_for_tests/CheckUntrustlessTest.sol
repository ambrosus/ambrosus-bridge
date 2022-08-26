// SPDX-License-Identifier: MIT
pragma solidity 0.8.6;

import "../networks/ETH_AmbBridge_Untrustless.sol";

contract CheckUntrustlessTest is CheckUntrustless {
    constructor() {}

    event checkUntrustlessTestResult(bool result);

    function checkUntrustlessTest(uint eventId, CommonStructs.Transfer[] calldata transfers) public {
        emit checkUntrustlessTestResult(checkUntrustless_(eventId, transfers));
    }

    function setRelaysAndConfirmationsTest(address[] calldata toRemove, address[] calldata toAdd, uint _confirmations) public {
        _setRelaysAndConfirmations(toRemove, toAdd, _confirmations);
    }

}
