```
pragma solidity ^0.4.11;

contract A {
  uint256 a;

  function A(uint256 _a) payable {
    a = _a;
  }
}
```

60606040526040516020806065833981016040528080519150505b60008190

55 // sstore

5b // tag_2
  50 // pop
5b // tag_3
  60 36 // dataSize(sub_0)
    stack: [0x36]
  80 // dup1
    stack: [0x36 0x36]
  60 2f // dataOffset(sub_0)
    stack: [0x2f 0x36 0x36]
  60 00 // 0x0
    stack: [0x00 0x2f 0x36 0x36]
  39 // codecopy
    // mOff = 0x00
    // cOff = 0x2f
    // l = 0x36

    // refers to this chunk: it's sub_0 + auxdata
    // 60606040525b600080fd00a165627a7a7230582069822b9b795f99ef3ebd52a351d6e76a987c74de859ba35435e56e36102c93ca0029
  60 00 // 60 00
  f3 // return
  00 // stop

  // sub_0
  60 60
  60 40
  52
  60 00
  80
  fd
  00

auxdata:
  a165627a7a7230582069822b9b795f99ef3ebd52a351d6e76a987c74de859ba35435e56e36102c93ca0029


b = "60606040526040516020806065833981016040528080519150505b60008190555b505b603680602f6000396000f30060606040525b600080fd00a165627a7a7230582069822b9b795f99ef3ebd52a351d6e76a987c74de859ba35435e56e36102c93ca0029"

```
EVM assembly:
    /* "contracts/A.sol":26:105  contract A {... */

  // 60 60 (PUSH 60)
  // 60 40 (PUSH 40)
  // 52 (MSTORE)
  mstore(0x40, 0x60)

  // 60 40 (PUSH 40)
  // 51 (MLOAD)
  mload(0x40)
    stack: [0x60]
  // 60 20 (PUSH 20)
  0x20
    stack: [0x20 0x60]
  // 80
  dup1
    stack: [0x20 0x20 0x60]
  // 60 65 (PUSH 65)
  bytecodeSize // 0x65 = 101
    stack: [0x65 0x20 0x20 0x60]
  // 83
  dup4
    stack: [0x60 0x65 0x20 0x20 0x60]
  // 39
  codecopy
    // mOff = 0x60
    // cOff = 0x65
    // l = 0x20

    // 0x65 is exactly the length of the bytecode (including auxdata.) Expect 32 bytes of data following that for contract creation.

    stack: [0x20 0x60]
    memory: [
      0x40: 0x40
      0x60: @arg1
    ]
  // update Solidity's own memory allocation counter
  dup2
    stack: [0x60 0x20 0x60]
  add
    stack: [0x80 0x60]
  0x40
    stack: [0x40 0x80 0x60]
  mstore
    stack: [0x60]
    memory: [
      0x40: 0x80
      0x60: @arg1
    ]
  dup1
    stack: [0x60 0x60]
  dup1
    stack: [0x60 0x60 0x60]
  // load @arg1 onto stack
  mload
    stack: [@arg1 0x60 0x60]
  swap2
    stack: [0x60 0x60 @arg1 ]
  pop
    stack: [0x60 @arg1 ]
  pop
    stack: [@arg1]
tag_1:
    /* "contracts/A.sol":92:93  a */
  0x0
    stack: [0x0 @arg1]
    /* "contracts/A.sol":92:98  a = _a */
  dup2
    stack: [@arg1 0x0 @arg1]
  swap1
    stack: [0x0 @arg1 @arg1]
  sstore
    stack: [@arg1]
    store: {
      0x0 => @arg1
    }
tag_2:
  pop
    stack: []
tag_3: // 5b600080fd00
  // 60 00
  dataSize(sub_0)
  //
  dup1
  dataOffset(sub_0)
  0x0
  codecopy
  0x0
  return
stop
```