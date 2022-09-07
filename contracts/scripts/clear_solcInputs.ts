import fs from "fs";
import path from "path";


async function main() {
  const basePath = path.resolve(__dirname, `../deployments/`);

  for (const stage of fs.readdirSync(basePath)) {
    for (const network of fs.readdirSync(basePath + "/" + stage)) {

      const networkPath = basePath + "/" + stage + "/" + network;
      const usefullSolcs = [];

      for (const deploymentName of fs.readdirSync(networkPath)) {
        if (!deploymentName.endsWith(".json")) continue;
        const deployment = require(networkPath + "/" + deploymentName);
        usefullSolcs.push(deployment.solcInputHash);
      }

      for (let solcName of fs.readdirSync(networkPath + "/solcInputs")) {
        if (!solcName.endsWith(".json")) continue;
        const solc = solcName.substring(0, solcName.length - 5);  // remove `.json` suffix

        if (usefullSolcs.includes(solc)) continue;

        console.log(`Removing ${networkPath}/solcInputs/${solcName}`);
        fs.unlinkSync(networkPath + "/solcInputs/" + solcName);
      }

    }
  }

}

main().then(() => process.exit(0)).catch(error => {
  console.error(error);
  process.exit(1);
});
