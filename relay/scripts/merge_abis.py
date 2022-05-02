import json
from pathlib import Path

ABI_FILES_PATH = Path.cwd().parent.parent.resolve() / "contracts" / "abi"


def main():
    amb_abi_file = open(ABI_FILES_PATH / "AmbBridge.json")
    eth_abi_file = open(ABI_FILES_PATH / "EthBridge.json")

    amb_abi = json.load(amb_abi_file)
    eth_abi = json.load(eth_abi_file)

    res = []
    already_used = []

    for network_abi in (amb_abi, eth_abi):
        for i in network_abi[1:-1]:  # skip "constructor" and "receive" func
            if i["name"] in already_used:
                continue
            res.append(i)
            already_used.append(i["name"])

    print(json.dumps(res, indent=2))


if __name__ == "__main__":
    main()
