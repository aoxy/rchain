# Exchange guide

This is a guide about how exchange can get infomation from [RChain](https://github.com/rchain/rchain) mainnet.


# Mainnet Servers

Currently RChain doesn't encourage you to start a node by yourself. But if you really want to do that , we can also provide some guides on that but that's another topic which is not discussed in this article.

Currently there are two kinds of nodes in RChain mainnet -> `Validator` and `Observer`.

## Validator

`Validator` is the node which packs all the deploys into a block and broadcast the block to the whole network. The transafer deploy should go to `Validator`. `Validator` is a staking node which needs real Rev to support and the private key to sign the block.

## Observer

`Observer` is a read-only node which provide more ways than `Validator` to get information from the node. But the observer doesn't have to ability to generate a block which it can only be `read-only` and it can not make any changes to the whole network. The transfer deploy which goes into `Observer` would never be packed into a block.

## Current Observer and Validator in Mainnet

Validator:

**node[0-29].root-shard.mainnet.rchain.coop**

(it means `node0.root-shard.mainnet.rchain.coop`, `node1.root-shard.mainnet.rchain.coop`, ...  and `node29.root-shard.mainnet.rchain.coop`)

Observer:

1. observer-asia.services.mainnet.rchain.coop
2. observer-eu.services.mainnet.rchain.coop
3. observer-us.services.mainnet.rchain.coop
4. observer-exch2.services.mainnet.rchain.coop(This is especially for exchange which we hope you use)



# Clients to interact with Mainnet

Currently, there are two available client.

1. python https://github.com/rchain/pyrchain
2. javascript https://github.com/rchain/exchange-api-njs

You can choose whatever client you want above. They work mostly similar.

There are example codes for both reposities.
1.  generate private key and public key-> [python](https://github.com/rchain/pyrchain/blob/master/examples/keys_example.py) , [javascript](https://github.com/rchain/exchange-api-njs/blob/master/examples/key_example.js)
2. sign a deploy with the private key -> [python](https://github.com/rchain/pyrchain/blob/master/examples/sign_verify_examples.py), [javascript](https://github.com/rchain/exchange-api-njs/blob/master/examples/sign_and_verify_example.js)
3. use client api to interact with rnode -> [python](https://github.com/rchain/pyrchain/blob/master/examples/grpc_api_example.py) , [javascript](https://github.com/rchain/exchange-api-njs/blob/master/examples/grpc_api_example.js)
4. Vault Api of rchain to do transfer and check balance -> [python](https://github.com/rchain/pyrchain/blob/master/examples/vault_example.py), [javascript](https://github.com/rchain/exchange-api-njs/blob/master/examples/vault_example.js)
5. Sign offline -> [python](https://github.com/rchain/pyrchain/blob/master/examples/transfer_sign_deploy.py) , [javascript](https://github.com/rchain/exchange-api-njs/blob/master/examples/transfer_sign.js)

## Get transfer message in one block

Getting transfer in a block is kind of tricky right now in RChain Mainnet. Both of the clients above actually support this functionality([python](https://github.com/rchain/pyrchain/blob/master/examples/transaction_example.py),[javaScript](https://github.com/rchain/exchange-api-njs/blob/master/examples/transaction_example.js)) 

But **the API above is very slow and the api would create great performance impact on the server which would make the service unstable**. To avoid that, currently coop provide a way which request the transfer info in http from a cache server.

**http://observer-exch2.services.mainnet.rchain.coop:7070/getTransaction/{blockHash}** (example :http://observer-exch2.services.mainnet.rchain.coop:7070/getTransaction/197823c8202d73dd6b02ce682d9bf6069f9f1a1781289c76ff94146b2921366a)

# Choosing the best validator to deploy

Currently the RChain network is proposed in sequence which means we got 30 nodes now and the blocks would be generated one by one. This leads to a problem that supposed that the network is proposed in sequence from 0 to 29, if the node1 is the next node to propose and you deploy the transfer in node0, it would take the whole round to process the deploy which is extremely slow. 

To solve this problem, coop provides an api which can help you know which node is the best node to propose and help you get the transfer on chain as fast as possible.

**https://status.rchain.coop/api/validators**

The feild in `nextToPropose` is the best validator to deploy when you query the request. It would be better to deploy anything in that server at that time.

