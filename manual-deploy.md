## This is manual to deploy temporary relay instance for AMB-ETH bridge.

### Installation

##### Creating a Droplet

Typical DigitalOcean node (2 GB / 2 CPUs, 60 GB SSD disk, 3 TB transfer) should be good for running a relay

There is detailed step by step information how to setup droplet on digitalocean https://www.digitalocean.com/docs/droplets/how-to/create/

Our brief instructions:

Create an account and log in. Press 'Droplets' and then 'Create Droplet'. Use the OS Ubuntu and then choose what machine preferences and which data center suits you. Then either create a SSH key which you will use to access the instance or if you do not choose one you will get a password to your email. Write a hostname that suits you and launch the instance.

### Obtaining private key

You can export the private key from your metamask wallet using this guide:
https://metamask.zendesk.com/hc/en-us/articles/360015289632-How-to-export-an-account-s-private-key

### Run

Command to run eth-relay:

<pre>
curl -s https://raw.githubusercontent.com/ambrosus/ambrosus-bridge/dev/relay.sh > relay.sh && sudo chmod +x relay.sh && sudo ./relay.sh
</pre>

If any pop-ups occur, select "ok".  
