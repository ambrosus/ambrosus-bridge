import path from "path";
import {execSync} from "child_process";


async function main() {
    const oldGasUsage = readReporterOutput(true)
    execSync(`CI=true yarn hardhat test`, {stdio: 'inherit'});
    const newGasUsage = readReporterOutput(false)

    const allKeys = Object.keys({...oldGasUsage, ...newGasUsage})
    for (const k of allKeys) {
        const [old, new_] = [oldGasUsage[k], newGasUsage[k]];
        if (old == new_) continue;

        console.log(k, ": ", old, " -> ", new_, '(', new_ - old, ')')
    }

}

function readReporterOutput(isOld: boolean): { [k: string]: number } {
    const output: { [k: string]: number } = {};
    const reporter = require(path.resolve(__dirname, `../gasReporterOutput${isOld?"Old":""}.json`))
    Object.values(reporter.info.methods).forEach((v: any) => {
        output["method " + v.contract + "." + v.fnSig] = avg(v.gasData);
    })
    Object.values(reporter.info.deployments).forEach((v: any) => {
        output["deploy " + v.name] = avg(v.gasData);
    })
    return output;
}

function avg(numbers: number[]): number {
    if (numbers.length == 0) return 0;
    return numbers.reduce((a, b) => a + b, 0) / numbers.length
}

main().catch(reason => {
    console.log(reason);
    process.exitCode = -1;
});
