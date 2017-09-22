# Prototype for Smart Contract deployment tool

Use env variable to specify the local qtumd RPC node:

```
export SOLAR_RPC=http://howard:yeh@localhost:13889
```

Specify an environment.

```
# The environment is `development` by default if you don't explicitly specify one
export SOLAR_ENV=development
```

To deploy a contract:

```
$ solar deploy contracts/Foo.sol
   deploy contracts/Foo.sol => foo
🚀  All contracts confirmed
```

You should see the address and ABI saved in a JSON file named `solar.development.json`:

```
{
  "foo": {
    "name": "Foo",
    "deployName": "foo",
    "address": "08227338906c17f8dcb1014c73a7ffc684c2376d",
    "txid": "5b5f95c2768b6945d6f0f98b4bb7a8621dc678f0f272ced0b75b5214e2c80b31",
    "abi": [
      {
        "name": "foo",
        "type": "function",
        "payable": false,
        "inputs": [
          {
            "name": "_a",
            "type": "uint256"
          },
          {
            "name": "_b",
            "type": "uint256"
          }
        ],
        "outputs": [],
        "constant": false
      },
      {
        "name": "",
        "type": "constructor",
        "payable": false,
        "inputs": [],
        "outputs": null,
        "constant": false
      }
    ],
    "bin": "60606040523415600e57600080fd5b5b5b5b608f8061001f6000396000f300606060405263ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166304bc52f88114603c575b600080fd5b3415604657600080fd5b60526004356024356054565b005b600082905560018190555b50505600a165627a7a723058204634a088cbc9a52bb7f659390c4460c7db41b5e1053814c26685b3fed6b07e280029",
    "binhash": "bc89293df211d6ac48ac8b84c8254f76aff602f703a941c7d8624f3ce6b1122e",
    "createdAt": "2017-09-22T19:59:36.714950514+08:00",
    "confirmed": true
  }
}
```

Add the flag `--force` to redeploy a contract:

```
$ solar deploy contracts/Foo.sol --force
   deploy contracts/Foo.sol => foo
🚀  All contracts confirmed
```