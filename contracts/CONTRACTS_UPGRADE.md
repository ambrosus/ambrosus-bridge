# Contracts upgrade reminder
> Violating any of these will cause the upgraded version of the contract to have its storage values mixed up, and can lead to critical errors.
## DO NOT:
1. Change the type of existing variable.
2. Change the order in which variables are declared.
3. Remove existing variable.
4. Swap the order in which the base contracts are declared, or introducing new base contracts.
5. Add new variables before existing ones
6. Add new variables to base contracts without decreasing the length of the gap, if the child has any variables of its own.

## Keep in mind:
1. **Memory packing:** [Solidity doc](https://docs.soliditylang.org/en/v0.8.14/internals/layout_in_storage.html#storage-inplace-encoding)
2. If you need to introduce a new variable, make sure you always do it at the end.
3. If variable will be removed from the end of contract, storage will not be cleared.
   A subsequent update that adds a new variable will cause that variable to read the leftover value from the deleted one.
4. Function and event/structure definitions can be added (or removed) anywhere in a contract, as they don't take up storage space.
5. Gap is uint256 because the EVM only operates on 32 bytes at a time.
6. If you rename a variable, then it will keep the same value as before after upgrading. This may be the desired behavior if the new variable is semantically the same as the old one
7. Mappings and dynamically-sized array occupy only 32 bytes, the elements they contain are stored starting at a different storage slot that is computed using a Keccak-256 hash
