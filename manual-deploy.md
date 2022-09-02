## This is manual to deploy temporary relay instance for AMB-ETH bridge.

## Installation

##### Creating a Droplet

Typical DigitalOcean node (2 GB / 2 CPUs, 60 GB SSD disk, 3 TB transfer) should be good for running a relay

There is detailed step by step information how to setup droplet on digitalocean https://www.digitalocean.com/docs/droplets/how-to/create/

Our brief instructions:

Create an account and log in. Press 'Droplets' and then 'Create Droplet'. Use the OS Ubuntu and then choose what machine preferences and which data center suits you. Then either create a SSH key which you will use to access the instance or if you do not choose one you will get a password to your email. Write a hostname that suits you and launch the instance.


TODO: Info about private key generation?

### Configuration

These are required environment variables to run relay. To apply them, simply put them in `env.list` file, which should be located in the same directory from where you start relay. The content of `env.list` file:
<pre>
STAGE=&lt;dev/test/prod>
NETWORK=eth-untrustless
NETWORKS_AMB_PRIVATEKEY=&lt;amb private key>
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
