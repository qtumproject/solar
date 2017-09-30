# Prototype for Smart Contract deployment tool

Start qtumd in regtest mode:

```
qtumd -regtest -rpcuser=howard -rpcpassword=yeh
```

Use env variable to specify the local qtumd RPC node:

```
export SOLAR_RPC=http://howard:yeh@localhost:13889
```

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

  function setA(uint256 _a) {
    a = _a;
  }

  function getA() returns(uint256) {
    return a;
  }
}
```

We need to give it a name when deploying. Let's call it `daisy`:

```
$ solar deploy contracts/A.sol daisy
    deploy contracts/A.sol => daisy
🚀  All contracts confirmed
```

(On a real network it would take longer to deploy. For development locally it is instantenous.)

You should see the address and ABI saved in a JSON file named `solar.development.json`:

```
{
  "daisy": {
    "name": "A",
    "deployName": "daisy",
    "address": "77a4190bdb5a01df293b0dd921f1a87f5c180620",
    "txid": "5ef2aa0c2b1d7fd41e3cf794b20617d9d35a0fe508227eb01057213f5c36355c",
    "abi": [
      {
        "name": "getA",
        "type": "function",
        "payable": false,
        "inputs": [],
        "outputs": [
          {
            "name": "",
            "type": "uint256"
          }
        ],
        "constant": false
      },
      {
        "name": "setA",
        "type": "function",
        "payable": false,
        "inputs": [
          {
            "name": "_a",
            "type": "uint256"
          }
        ],
        "outputs": [],
        "constant": false
      }
    ],
    "bin": "6060604052341561000f57600080fd5b5b60b98061001e6000396000f300606060405263ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663d46300fd81146046578063ee919d50146068575b600080fd5b3415605057600080fd5b6056607d565b60405190815260200160405180910390f35b3415607257600080fd5b607b6004356084565b005b6000545b90565b60008190555b505600a165627a7a723058203431ad594c9688027a5ac39ec60fbb0786fc861d6d51417f600fe03b9412752a0029",
    "binhash": "42712271c9f5e5dcd27eaeb999bf4388eb80c55cd652980a7b22aa34f774d76b",
    "createdAt": "2017-09-30T16:40:15.656957558+08:00",
    "confirmed": true
  }
}
```

You can't reuse the same name twice. You'll get a warning:

```
$ solar deploy contracts/A.sol daisy
   deploy contracts/A.sol => daisy
❗️  deploy name already used: daisy
```

Add the flag `--force` to redeploy a contract:

```
$ solar deploy contracts/A.sol daisy --force
   deploy contracts/Foo.sol => foo
🚀  All contracts confirmed
```

In `solar.development.json` you should see that the address had changed.