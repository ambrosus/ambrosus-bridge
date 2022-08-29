## This is manual to deploy relay instances for both ETH and BSC bridges.

### Prerequisites
In order to run relay on your machine you need to have docker installed.

Also you need to be authorized to pull image from private docker registry. In this case you need to authorize using this [instructions]("https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry")
### Relay ETH

Command to run eth-relay:

<pre>docker run -d --name eth-relay --restart unless-stopped --env-file ./env.list -v ethash:/root/ethash -p 127.0.0.1:8080:8080 -p 127.0.0.1:9090:9090 ghcr.io/ambrosus/ambrosus-bridge</pre>

These are required env variables to run relay. To apply them, simply put them in `env.list` file, which should be located in the same directory where you run the command. The content of `env.list` file:
<pre>
STAGE=&lt;dev/test/prod>
NETWORK=eth
NETWORKS_AMB_PRIVATEKEY=&lt;amb private key>
NETWORKS_ETH_PRIVATEKEY=&lt;eth private key>
FEEAPI_AMB_PRIVATEKEY=&lt;feeapi amb private key>
FEEAPI_SIDE_PRIVATEKEY=&lt;feeapi side private key>
NETWORK_ETH_ETHASHDIR=/root/ethash
</pre>
After container is started you can remove this file.


### Relay BSC

Command to run bsc-relay:

<pre>docker run -d --name bsc-relay --restart unless-stopped --env-file ./env.list -p 127.0.0.1:8080:8080 -p 127.0.0.1:9091:9091 ghcr.io/ambrosus/ambrosus-bridge</pre>

Content of `env.list` file:

<pre>
STAGE=&lt;dev/test/prod>
NETWORK=bsc
NETWORKS_AMB_PRIVATEKEY=&lt;amb private key>
NETWORKS_BSC_PRIVATEKEY=&lt;bsc private key>
FEEAPI_AMB_PRIVATEKEY=&lt;feeapi amb private key>
FEEAPI_SIDE_PRIVATEKEY=&lt;feeapi side private key>
</pre>