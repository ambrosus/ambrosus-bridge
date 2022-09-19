import sys
import json
from pathlib import Path

ABI_FILES_PATH = Path.cwd().parent.parent.resolve() / "contracts" / "abi"


def main():
    abis = (
        json.load((ABI_FILES_PATH / f).open())
        for f in sys.argv[1:]
    )

    res = []
    already_used = []

    for network_abi in abis:
        for i in network_abi:
            if i["type"] in ("constructor", "receive", "fallback"):
                continue
            if i["name"] in already_used:
                continue
            res.append(i)
            already_used.append(i["name"])

    print(json.dumps(res, indent=2))


if __name__ == "__main__":
    main()
