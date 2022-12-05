import fs from "fs";
import path from "path";


function main() {
  const args = process.argv.slice(2);
  const abiPath = path.resolve(__dirname, `../abi/`);

  function fileNames() {
    if (args.length == 1) {
      if (args[0] == "--all")
        return fs.readdirSync(abiPath).filter(file => file.endsWith(".json"))
      if (args[0] == "--relay") {
        return [
          "ETH_AmbBridge.json", "ETH_EthBridge.json",
          "BSC_AmbBridge.json", "BSC_BscBridge.json"
        ];
      }
      if (args[0] == "--multisig") {
        return [
          "ETH_AmbBridge.json", "ETH_EthBridge.json",
          "BSC_AmbBridge.json", "BSC_BscBridge.json",
          "ProxyMultiSig.json"
        ];
      }
    }
    return args
  }

  const filePaths = fileNames().map(file => path.resolve(abiPath, file));

  merge(filePaths);
}

function merge(filePaths: string[]) {

  const mergedAbi = [];
  const names = new Set();

  for (const filePath of filePaths) {
    const abi = require(filePath) as any[];


    for (const item of abi) {
      if (["constructor", "receive", "fallback"].includes(item.type)) continue;
      if (names.has(item.name)) continue;

      mergedAbi.push(item);
      names.add(item.name);

    }
  }

  console.log(JSON.stringify(mergedAbi, null, 2));

}

main()
