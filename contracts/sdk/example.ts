import { AllowanceException, Network, Token } from "./types";
import { ethers } from "ethers";
import { BridgeSDK } from "./index";

class MySdk extends BridgeSDK {
  getPairs(sourceNetwork: Network, destinationNetwork: Network): [Token, Token][] {
    const result = super.getPairs(sourceNetwork, destinationNetwork);
    return result.map(([from, to]) => [this.mapToken(from), this.mapToken(to)]);
  }

  mapToken(token: Token): Token {
    // todo logos
    // todo if isNativeCoin you can change name or symbol, for example:
    //   use networkToken.nativeCoin (it's a symbol of native coin)
    //   or create a mapping (networkName => {name, symbol, logo, ...}) for native coins

    if (token.isNativeCoin)
      token.name += " (NATIVE)";
    return token;
  }
}


async function main() {

  const config = {
    "bridges": {
      "eth": {
        "amb": "0x19caBC1E34Ab0CC5C62DaA1394f6022B38b75c78",
        "side": "0x0De2669e8A7A6F6CC0cBD3Cf2D1EEaD89e243208"
      },
      "bsc": {
        "amb": "0xe10eB55f6EeF66218BbE58B749428ec4A51D6659",
        "side": "0x92fa52d3043725D00Eab422440C4e9ef3ba180d3"
      }
    },
    "tokens": {
      "SAMB": {
        "isActive": true,
        "name": "Synthetic Amber",
        "symbol": "SAMB",
        "logo": "https://media-exp1.licdn.com/dms/image/C560BAQFuR2Fncbgbtg/company-logo_200_200/0/1636390910839?e=2159024400&v=beta&t=W0WA5w02tIEH859mVypmzB_FPn29tS5JqTEYr4EYvps",
        "networks": {
          "amb": {
            "address": "0x2b2d892C3fe2b4113dd7aC0D2c1882AF202FB28F",
            "denomination": 18,
            "isPrimary": true,
            "nativeCoin": "AMB"
          },
          "bsc": {
            "address": "0x23c1C1cc14270B7Bd63677d1fe4790891b17A33d",
            "denomination": 18,
            "isPrimary": false
          },
          "eth": {
            "address": "0x683aae5cD37AC94943D05C19E9109D5876113562",
            "denomination": 18,
            "isPrimary": false
          }
        }
      },
      "SAMB2": {
        "isActive": true,
        "name": "AirDAO",
        "symbol": "AMB",
        "logo": "https://media-exp1.licdn.com/dms/image/C560BAQFuR2Fncbgbtg/company-logo_200_200/0/1636390910839?e=2159024400&v=beta&t=W0WA5w02tIEH859mVypmzB_FPn29tS5JqTEYr4EYvps",
        "networks": {
          "amb": {
            "address": "0x2b2d892C3fe2b4113dd7aC0D2c1882AF202FB28F",
            "denomination": 18,
            "isPrimary": true,
            "nativeCoin": "AMB"
          },
          "bsc": {
            "address": "",
            "denomination": 18,
            "isPrimary": false
          },
          "eth": {
            "address": "0xf4fB9BF10E489EA3Edb03E094939341399587b0C",
            "denomination": 18,
            "isPrimary": false
          }
        }
      },
      "WETH": {
        "isActive": false,
        "name": "Wrapped Ether",
        "symbol": "WETH",
        "logo": "https://ethereum.org/static/bfc04ac72981166c740b189463e1f74c/448ee/eth-diamond-black-white.webp",
        "networks": {
          "amb": {
            "address": "0xe7c3607474E235Ec8deF1f0a63Ea983538eea182",
            "denomination": 18,
            "isPrimary": false
          },
          "eth": {
            "address": "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
            "denomination": 18,
            "isPrimary": true,
            "nativeCoin": "ETH"
          }
        }
      },
      "WBNB": {
        "isActive": false,
        "name": "Wrapped BNB",
        "symbol": "WBNB",
        "logo": "https://bscscan.com/token/images/binance_32.png",
        "networks": {
          "amb": {
            "address": "0xA96C522fA8Df99BB73A6E317A1afb0E3FA13b735",
            "denomination": 18,
            "isPrimary": false
          },
          "bsc": {
            "address": "0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c",
            "denomination": 18,
            "isPrimary": true,
            "nativeCoin": "BNB"
          }
        }
      },
      "USDC": {
        "isActive": true,
        "name": "USD Coin",
        "symbol": "USDC",
        "logo": "https://etherscan.io/token/images/centre-usdc_28.png",
        "networks": {
          "amb": {
            "address": "0xFF9F502976E7bD2b4901aD7Dd1131Bb81E5567de",
            "denomination": 18,
            "isPrimary": false
          },
          "bsc": {
            "address": "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d",
            "denomination": 18,
            "isPrimary": true
          },
          "eth": {
            "address": "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
            "denomination": 6,
            "isPrimary": true
          }
        }
      },
      "BUSD": {
        "isActive": false,
        "name": "BUSD Token",
        "symbol": "BUSD",
        "logo": "https://bscscan.com/token/images/busd_32.png",
        "networks": {
          "amb": {
            "address": "0x7A477aA8ED4884509387Dba81BA6F2B7C97597e2",
            "denomination": 18,
            "isPrimary": false
          },
          "bsc": {
            "address": "0xe9e7CEA3DedcA5984780Bafc599bD69ADd087D56",
            "denomination": 18,
            "isPrimary": true
          }
        }
      },
      "USDT": {
        "isActive": true,
        "name": "Tether USD",
        "symbol": "USDT",
        "logo": "https://etherscan.io/token/images/tether_32.png",
        "networks": {
          "amb": {
            "address": "0xfEE01F2D120250A0a59bfbF9C144F8ECC4425fCc",
            "denomination": 18,
            "isPrimary": false
          },
          "bsc": {
            "address": "0x55d398326f99059fF775485246999027B3197955",
            "denomination": 18,
            "isPrimary": true
          },
          "eth": {
            "address": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
            "denomination": 6,
            "isPrimary": true
          }
        }
      }
    },
    "ambFaucetAddress": "0x5B72903b7D2711134DD2120Bc3CA61DB6f10cB92"
  };
  const providerAmb = new ethers.providers.JsonRpcProvider("https://network.ambrosus.io");
  const providerEth = new ethers.providers.JsonRpcProvider("https://eth.llamarpc.com");
  const providerBsc = new ethers.providers.JsonRpcProvider("https://bsc.llamarpc.com");

  const signer = new ethers.Wallet("34d8e83fca265e9ab5bcc1094fa64e98692375bf8980d066a9edcf4953f0f2f5", providerAmb);

  const sdk = new MySdk(config, {
    relayUrls: {
      eth: "https://relay-eth.ambrosus.io/fees",
      bsc: "https://relay-bsc.ambrosus.io/fees"
    }
  });


  const ethAmb = {name: "eth->amb", pairs: sdk.getPairs("eth", "amb"), provider: providerEth};
  const ambEth = {name: "amb->eth", pairs: sdk.getPairs("amb", "eth"), provider: providerAmb};
  const bscAmb = {name: "bsc->amb", pairs: sdk.getPairs("bsc", "amb"), provider: providerBsc};
  const ambBsc = {name: "amb->bsc", pairs: sdk.getPairs("amb", "bsc"), provider: providerAmb};

  for (const { name, pairs, provider } of [ethAmb, ambEth, bscAmb, ambBsc]) {
    for (const [from, to] of pairs) {


      console.log(name, from.name, "->", to.name)

      const feeData = await sdk.getFeeData(from, to, "10000", true);

      try {
        await sdk.withdraw(from, to, signer.address, "10000", feeData, signer.connect(provider));
      } catch (e) {
        if (e instanceof AllowanceException) {
          await sdk.setAllowance(e.token, signer, e.spenderAddress, e.amount);
          // try again
          await sdk.withdraw(from, to, signer.address, "10000", feeData, signer.connect(provider));
        }
        throw e;
      }

      console.log()


    }
  }


}

/*
example output:

eth->amb Tether USD -> Tether USD
Need to increase allowance to 1000000
bridge.withdraw(0xdAC17F958D2ee523a2206206994597C13D831ec7, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000, false, 0x5a4fe32bc10384df40b73dd90478c2010118bdf66aa403631c1577b88835882845c7e3bab89345fb6d3ada9e22821825ca0dbca5fa0caf1d77b3288b238a05c700, 158020501059307, 0, { value: 158020501059307 })


eth->amb USD Coin -> USD Coin
Need to increase allowance to 1000000
bridge.withdraw(0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000, false, 0xbc14e2db76772e605cf3e23ae650b49b416f0f0ae812184d311e079d3e9784a70f6389597e8250555ea6d158fa8614f891ea63480327dc4e5e9c668e40d8393000, 158020501059307, 0, { value: 158020501059307 })

eth->amb AirDAO -> AirDAO
Need to increase allowance to 1000000000000000000
bridge.withdraw(0xf4fB9BF10E489EA3Edb03E094939341399587b0C, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0x9e4a9bd978c04ca16a7c55901f49beddbfb12c571409b02272ed4da99549730c69aaa2ab35d154ff94ed9193e90ab50022c659815527af23d64a823b258b0c4501, 158020501059307, 0, { value: 158020501059307 })

eth->amb AirDAO -> AirDAO (NATIVE)
Need to increase allowance to 1000000000000000000
bridge.withdraw(0xf4fB9BF10E489EA3Edb03E094939341399587b0C, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, true, 0x9e4a9bd978c04ca16a7c55901f49beddbfb12c571409b02272ed4da99549730c69aaa2ab35d154ff94ed9193e90ab50022c659815527af23d64a823b258b0c4501, 158020501059307, 0, { value: 158020501059307 })

eth->amb Synthetic Amber -> Synthetic Amber
Need to increase allowance to 1000000000000000000
bridge.withdraw(0x683aae5cD37AC94943D05C19E9109D5876113562, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0xe95d03fa248a46d43c8e92ffb7377156943cc38bd5fd74a4ff318be945535cf160f171c803204db4458e0b2ecb218b10da0350f309c30e26264377436ec3f74a01, 158020501059307, 0, { value: 158020501059307 })

amb->eth Tether USD -> Tether USD
Need to increase allowance to 1000000000000000000
bridge.withdraw(0xfEE01F2D120250A0a59bfbF9C144F8ECC4425fCc, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000, false, 0xd2a89974fc1475cc7e09b69183cb8332543c6f1a4b7cbb2df0036b093d45cb045f8011e0a3fe63cc648f6321d91bcca5c701bd9850d1af72262a1cb5c1408bbb00, 4723877201423040465425, 0, { value: 4723877201423040465425 })

amb->eth USD Coin -> USD Coin
Need to increase allowance to 1000000000000000000
bridge.withdraw(0xFF9F502976E7bD2b4901aD7Dd1131Bb81E5567de, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000, false, 0xa1b326116eb5992e87bb2d769abbccedad896e07e93091f8afeebdd05b5b56c42f3be36fc65bf6857a891ea6b7e67b729aebd214ec6219e3a75c9f07488dc6da01, 4723877201423040465425, 0, { value: 4723877201423040465425 })

amb->eth AirDAO -> AirDAO
Need to increase allowance to 1000000000000000000
bridge.withdraw(0x2b2d892C3fe2b4113dd7aC0D2c1882AF202FB28F, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0xa378fe61fb70da4ecc01f7834f3cc0d0f4c2af32ea66829b8693bed05ff149fb533d46c407ec79eefceefe3637bfd9f0eed483f0c107d96f44db11e66ddfef8101, 4723877201423040465425, 0, { value: 4723877201423040465425 })

amb->eth AirDAO (NATIVE) -> AirDAO
bridge.wrapWithdraw(0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 0xa378fe61fb70da4ecc01f7834f3cc0d0f4c2af32ea66829b8693bed05ff149fb533d46c407ec79eefceefe3637bfd9f0eed483f0c107d96f44db11e66ddfef8101, 4723877201423040465425, 0, { value: 4724877201423040465425 })

amb->eth Synthetic Amber -> Synthetic Amber
Need to increase allowance to 1000000000000000000
bridge.withdraw(0x2b2d892C3fe2b4113dd7aC0D2c1882AF202FB28F, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0xa378fe61fb70da4ecc01f7834f3cc0d0f4c2af32ea66829b8693bed05ff149fb533d46c407ec79eefceefe3637bfd9f0eed483f0c107d96f44db11e66ddfef8101, 4723877201423040465425, 0, { value: 4723877201423040465425 })

amb->eth Synthetic Amber (NATIVE) -> Synthetic Amber
bridge.wrapWithdraw(0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 0xa378fe61fb70da4ecc01f7834f3cc0d0f4c2af32ea66829b8693bed05ff149fb533d46c407ec79eefceefe3637bfd9f0eed483f0c107d96f44db11e66ddfef8101, 4723877201423040465425, 0, { value: 4724877201423040465425 })

bsc->amb Tether USD -> Tether USD
Need to increase allowance to 1000000000000000000
bridge.withdraw(0x55d398326f99059fF775485246999027B3197955, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0xb5a4c1a9e86f7fc0896eb9d7baaca93deb0e065fbe6f63344ac0bcd4ff1a76cb555d0d38828cb0e248b3f476806001dd9dbd335c0497cc08d9bf385ec5b2fa2701, 848777785990751, 0, { value: 848777785990751 })

bsc->amb USD Coin -> USD Coin
Need to increase allowance to 1000000000000000000
bridge.withdraw(0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0x7ca1930d7ba742063c4967e742b111a6fcc26087e0a0ffb5a08b5012b83e1495186673f29d00d9fb36ba28b20c66794d2431c5e8811ab333f0377bf5c7469c9601, 848777785990751, 0, { value: 848777785990751 })

bsc->amb Synthetic Amber -> Synthetic Amber
Need to increase allowance to 1000000000000000000
bridge.withdraw(0x23c1C1cc14270B7Bd63677d1fe4790891b17A33d, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0x2a08d40d3191e3a930160179a36041ac9670793e10612e8e21ebe6bae8d6a148192d7bf195f978c7181a8048335ee5affc401f6f9c73d598baae6fae2ec55de400, 848777785990751, 0, { value: 848777785990751 })

bsc->amb Synthetic Amber -> Synthetic Amber (NATIVE)
Need to increase allowance to 1000000000000000000
bridge.withdraw(0x23c1C1cc14270B7Bd63677d1fe4790891b17A33d, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, true, 0x2a08d40d3191e3a930160179a36041ac9670793e10612e8e21ebe6bae8d6a148192d7bf195f978c7181a8048335ee5affc401f6f9c73d598baae6fae2ec55de400, 848777785990751, 0, { value: 848777785990751 })

amb->bsc Tether USD -> Tether USD
Need to increase allowance to 1000000000000000000
bridge.withdraw(0xfEE01F2D120250A0a59bfbF9C144F8ECC4425fCc, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0x0b022a34ce9ed0265e316ae48a216e793e35be2ca54d503e91adb6092a825ed343b96ae8f892b4ee14755b683561ca9048529eef47bd3e0230f65b7df90828e300, 121715358895993008670, 0, { value: 121715358895993008670 })

amb->bsc USD Coin -> USD Coin
Need to increase allowance to 1000000000000000000
bridge.withdraw(0xFF9F502976E7bD2b4901aD7Dd1131Bb81E5567de, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0x85fd72dbad530f74c5e5a5fab9e4860ddb380acb89976aaa117559b64132ca891d9ee64ed520570de9a552aa07dc3f14b78a9194671e7c42dbd9ce423b79bd3201, 121715358895993008670, 0, { value: 121715358895993008670 })

amb->bsc Synthetic Amber -> Synthetic Amber
Need to increase allowance to 1000000000000000000
bridge.withdraw(0x2b2d892C3fe2b4113dd7aC0D2c1882AF202FB28F, 0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 1000000000000000000, false, 0x33ba355e51585860a96292abda5de4e14c3369c7b538ec1c6a26bbb01637d8c31cf3eee7c1ebfdeddbd5885d58a60e89e9038a8cffe0cb83a78bdc57a88cf5a901, 121715358895993008670, 0, { value: 121715358895993008670 })

amb->bsc Synthetic Amber (NATIVE) -> Synthetic Amber
bridge.wrapWithdraw(0x295C2707319ad4BecA6b5bb4086617fD6F240CfE, 0x33ba355e51585860a96292abda5de4e14c3369c7b538ec1c6a26bbb01637d8c31cf3eee7c1ebfdeddbd5885d58a60e89e9038a8cffe0cb83a78bdc57a88cf5a901, 121715358895993008670, 0, { value: 122715358895993008670 })

 */

main();
