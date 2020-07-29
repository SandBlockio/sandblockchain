#### Sandblock Chain Core

------------

👋 Welcome to the official Sandblock Chain Core!

⚠️ This is beta software. Some features may not work as intended at the very beginning.

🤓 All contributions are more than welcome! Feel free to fork the repository and create a Pull Request!

------------

##### Getting started

    $ git clone git@github.com:SandBlockio/sandblock-chain.git
    $ cd sandblock-chain
    
Make sure your Go configuration looks like that one (especially the GO111Module):

    export GOPATH=$HOME/dev/go
    export GOBIN=$GOPATH/bin
    export GO111MODULE=on
    export PATH=$PATH:$GOBIN

Then compile (this will take care of installing the binaries in the PATH once compiled)

    $ make

Then init the genesis state (local is node name, could be anything u want)

    $ sbd init local --chain-id sandblockchain

Configure initial wallets (save the adresses for later)

    $ sbcli keys add enguerrand
    $ sbcli keys add fabrice

Add the initial wallets in the genesis block with initial coins and initial staking coins

    $ sbd add-genesis-account $(sbcli keys show enguerrand -a) 40000000sbc
    $ sbd add-genesis-account $(sbcli keys show fabrice -a) 2500000sbc

Apply the last configurations

    $ sbcli config chain-id sandblockchain
    $ sbcli config output json
    $ sbcli config indent true
    $ sbcli config trust-node true

Then we finally generate the genesis block (replace enguerrand with one of the previous wallet if u want)

    $ sbd gentx --name enguerrand --amount=10000000sbc

We validate the configurations

    $ sbd collect-gentxs
    $ sbd validate-genesis

And we can start :)

    $ sbd start

##### Creating the first branded token
We can see by doing that command that the wallet enguerrand only contains surprisecoin, no branded token:

    $ sbcli query account $(sbcli keys show enguerrand -a)

We will create a new branded token named "brandedtoken1" with 1000 of initial supply from the wallet enguerrand

    $ sbcli tx surprise create-token brandedtoken1 1000 --from enguerrand

If we request again we will see the wallet now has a new coin inside it :)

    $ sbcli query account $(sbcli keys show enguerrand -a)

##### Transfering branded tokens

Now let's transfer brandedtoken1 units to wallet fabrice (replace the address with the one from the wallet)

    $ sbcli tx surprise transfer-token brandedtoken1 cosmos168p32u2h4z4c0x4w9ahwfzntqcpnxs2ht38v9s 500 --from enguerrand

We will now see that wallet enguerrand and wallet fabrice both contains 500 units of the brandedtoken1

    $ sbcli query account $(sbcli keys show enguerrand -a)
    $ sbcli query account $(sbcli keys show fabrice -a)

##### Requesting infos about branded token
We can query the list of created branded tokens by using

    $ sbcli query surprise fetch

And then get informations about a given token with

    $ sbcli query surprise get brandedtoken1

##### Connecting a second node to the network
We can connect a second node to the network by initializing it:

    $ sbd init local-2 --chain-id sandblockchain

Then on the first node get the ID by tapping the command:

    $ sbcli status

Then on the second node, edit the config file:

    $ nano /.sandblockchain/config/config.toml

Still on the second, replace the line (edit values):

    persistent_peers = "first_node_id@firt_node_ip:26656"
    
Copy the genesis file from bootnode to your new node

After genesis file is copied, still on the second one, start it :)

    $ sbd start
    
##### Run a HTTP server
You can run the built-in HTTP server with the given instruction:

    $ sbcli rest-server --chain-id sandblockchain --trust-node

##### Interacting with the HTTP server
Lets create another token via the HTTP API, we start by getting our account sequence number (save it for later):

    $ curl -s http://localhost:1317/auth/accounts/$(sbcli keys show enguerrand -a)
    
Then we generate a raw transaction and export it on local file

    $ curl -XPOST -s http://localhost:1317/surprise/token --data-binary '{"base_req":{"from":"'$(sbcli keys show enguerrand -a)'","chain_id":"sandblockchain"},"name":"bernard","amount":"512"}' > unsignedTx.json

We can now sign it :) (replace sequence number here)

    $ sbcli tx sign unsignedTx.json --from enguerrand --offline --chain-id sandblockchain --sequence 3 --account-number 0 > signedTx.json

And finally broadcast it

    $ sbcli tx broadcast signedTx.json


##### Start the project using Docker
    $ docker build -t sandblockchain .
    $ docker run -it -p 26657:26657 -p 26656:26656 -v ~/.sbd:/home/sandblockchain/.sbd -v ~/.sbcli:/home/sandblockchain/.sbcli sandblockchain sbd start
