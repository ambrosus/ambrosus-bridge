## This is manual to deploy temporary relay instance for AMB-ETH bridge.

### Prerequisites
In order to run relay on your machine you need to have docker installed.

Also you need to be authorized to pull image from private docker registry. In this case you need to authorize using this [instructions]("https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry").

TODO: Info about private key generation?

### Configuration

These are required env variables to run relay. To apply them, simply put them in `env.list` file, which should be located in the same directory where you run the command. The content of `env.list` file:
<pre>
STAGE=&lt;dev/test/prod>
NETWORK=eth
NETWORKS_AMB_PRIVATEKEY=&lt;amb private key>
NETWORKS_ETH_PRIVATEKEY=&lt;eth private key>
</pre>
After container is started you can remove this file.

### Run

Command to run eth-relay:

<pre>
docker run -d \
--name eth-relay \
--restart unless-stopped \
--env-file ./env.list \
ghcr.io/ambrosus/ambrosus-bridge
</pre>
