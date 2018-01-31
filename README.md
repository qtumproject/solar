# Install

```
go get -u github.com/qtumproject/solar/cli/solar
```

`solar` assumes that the [Solidity compiler](https://github.com/ethereum/solidity) is already installed.

# Prototype for Smart Contract deployment tool

## QTUM

Start qtumd in regtest mode:

```
qtumd -regtest -rpcuser=howard -rpcpassword=yeh
```

Use env variable to specify the local qtumd RPC node:

```
export QTUM_RPC=http://howard:yeh@localhost:13889
```

## QTUM Docker

You can run qtumd with docker, which comes bundled with solar (and `solc`):

```
docker run -it --rm \
  --name myapp \
  -v `pwd`:/dapp \
  -p 3889:3889 \
  hayeah/qtumportal
```

Then you enter into the container by running:

```
docker exec -it myapp sh
```

## Ethereum

Start eth in private network

```
geth --rpc --rpcapi "eth,miner,personal" --datadir "./" --nodiscover console
```

then open a new tab

```
export ETH_RPC=http://accountAddress:accountPassword@localhost:8545
```

`solar` will let you enter the account and password, if you does not set them

Specify an environment.

```
# The environment is `development` by default if you don't explicitly specify one
export SOLAR_ENV=development
```

# Deploy Contract

Suppose we have the following contract in `contracts/Foo.sol`:

```
pragma solidity ^0.4.11;

contract A {
  uint256 a;

  function A(uint256 _a) {
    a = _a;
  }
}
```

Use the `deploy` command to create the contract. The constructor parameters may be passed in as a JSON array:

```
$ solar deploy contracts/A.sol '[100]'
🚀  All contracts confirmed
deployed contracts/A.sol => 2672931f2f07b67383a4aec80fb6504727145e8c
```

> Note: On a real network it would take longer to deploy. For development locally it is instantenous.

You should see the address and ABI saved in a JSON file named `solar.development.json`:

```json
{
  "contracts": {
    "contracts/A.sol": {
      "source": "contracts/A.sol:A",
      "abi": [
        {
          "name": "",
          "type": "constructor",
          "payable": false,
          "inputs": [
            {
              "name": "_a",
              "type": "uint256",
              "indexed": false
            }
          ],
          "outputs": null,
          "constant": false,
          "anonymous": false
        }
      ],
      "bin": "60606040523415600e57600080fd5b604051602080606783398101604052808051600055505060358060326000396000f3006060604052600080fd00a165627a7a72305820cd047550e6b1a360fc0f2de526c2f6f0150d249efca0b5820d6050c935e129370029",
      "binhash": "861924fb216a600b22bc38d51d26b708bbd2d189a3433f0ec4862da6a3d78b9a",
      "name": "A",
      "deployName": "contracts/A.sol",
      "address": "2672931f2f07b67383a4aec80fb6504727145e8c",
      "txid": "dc36ebb365033f557367c88e4bad2f4c726a609e8a4cc0d5751ff4cab9187a51",
      "createdAt": "2018-01-26T11:36:09.650368039+08:00",
      "confirmed": true,
      "sender": "qYqieW18XiFU1sYCrHFsiBpvxQVXpcwC3R",
      "senderHex": "a7a0cff24ecf5089a5b5a814b9a6be942ade51c5"
    }
  },
  "libraries": {},
  "related": {}
}
```

The contract's deploy name defaults to the contract's file name. You can't deploy using the same name twice.

```
$ solar deploy contracts/A.sol
❗️  deploy contract name already used: contracts/A.sol
```

Add the flag `--force` to redeploy a contract:

```
$ solar deploy contracts/A.sol '[200]' --force
🚀  All contracts confirmed
   deployed contracts/A.sol => 3b566019c5e69e3c33c4609b51597b4b83b00e74
```

In `solar.development.json` you should see that the address had changed.

## Naming A Contract

By default the file name is used as the deploy name. We can change the deploy name to `Ace` by appending it after the file name:

```
$ solar deploy contracts/A.sol:Ace '[300]'
🚀  All contracts confirmed
   deployed Ace => 38dfbad58cb340815e649242fc7514fe1fa54e7d
```

The deploy name must be unique. Deploy should fail if the a name is used twice:

```
$ solar deploy contracts/A.sol:Ace '[300]'
❗️  deploy contract name already used: Ace
```

## Deploy A Library

Sometimes your contract may depend on a library, which requires linking when deploying. The most frequently used library is maybe `SafeMath.lib`.

```
solar deploy contracts/SafeMathLib.sol --lib

🚀  All contracts confirmed
   deployed contracts/SafeMathLib.sol => 26fe40c433b4d109299660284eaa475b95462342
```

Then deploying other contracts that reference `SafeMathLib` will automatically link to this deployed instance.

# deploy

Help message for `deploy`:

```
solar help deploy

usage: solar deploy [<flags>] <target> [<jsonParams>]

Compile Solidity contracts.

Flags:
  --help                     Show context-sensitive help (also try --help-long and --help-man).
  --qtum_rpc=QTUM_RPC        RPC provider url
  --qtum_sender=QTUM_SENDER  (qtum) Sender UXTO Address
  --eth_rpc=ETH_RPC          RPC provider url
  --env="development"        Environment name
  --repo=REPO                Path of contracts repository
  --optimize                 [solc] should Enable bytecode optimizer
  --allow-paths=""           [solc] Allow a given path for imports. A list of paths can be supplied by separating them with a comma.
  --force                    Overwrite previously deployed contract with the same deploy name
  --lib                      Deploy the contract as a library
  --no-confirm               Don't wait for network to confirm deploy
  --no-fast-confirm          (dev) Don't generate block to confirm deploy immediately
  --gasLimit=3000000         gas limit for creating a contract

Args:
  <target>        Solidity contracts to deploy.
  [<jsonParams>]  Constructor params as a json array
```
