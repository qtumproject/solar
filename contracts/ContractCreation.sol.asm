
======= contracts/ContractCreation.sol:C1 =======
EVM assembly:
    /* "contracts/ContractCreation.sol":47:455  contract C1 {... */
  mstore(0x40, 0x60)
    /* "contracts/ContractCreation.sol":159:197  function C1() {... */
  jumpi(tag_1, iszero(callvalue))
  0x0
  dup1
  revert
tag_1:
tag_2:
    /* "contracts/ContractCreation.sol":184:192  new C2() */
  tag_4
  jump	// in(tag_5)
tag_4:
  mload(0x40)
  dup1
  swap2
  sub
  swap1
  0x0
  create
  dup1
  iszero
  iszero
  tag_6
  jumpi
  0x0
  dup1
  revert
tag_6:
    /* "contracts/ContractCreation.sol":179:181  c2 */
  0x0
    /* "contracts/ContractCreation.sol":179:192  c2 = new C2() */
  dup1
  sload
  not(sub(exp(0x2, 0xa0), 0x1))
  and
  sub(exp(0x2, 0xa0), 0x1)
  swap3
  swap1
  swap3
  and
  swap2
  swap1
  swap2
  or
  swap1
  sstore
    /* "contracts/ContractCreation.sol":159:197  function C1() {... */
tag_3:
    /* "contracts/ContractCreation.sol":47:455  contract C1 {... */
  jump(tag_7)
tag_5:
  mload(0x40)
  dataSize(sub_1)
  dup1
  dataOffset(sub_1)
  dup4
  codecopy
  add
  swap1
  jump	// out
tag_7:
  dataSize(sub_0)
  dup1
  dataOffset(sub_0)
  0x0
  codecopy
  0x0
  return
stop

sub_0: assembly {
        /* "contracts/ContractCreation.sol":47:455  contract C1 {... */
      mstore(0x40, 0x60)
      and(div(calldataload(0x0), 0x100000000000000000000000000000000000000000000000000000000), 0xffffffff)
      0xf207564e
      dup2
      eq
      tag_2
      jumpi
    tag_1:
      0x0
      dup1
      revert
        /* "contracts/ContractCreation.sol":201:453  function register(uint value) {... */
    tag_2:
      jumpi(tag_3, iszero(callvalue))
      0x0
      dup1
      revert
    tag_3:
      tag_4
      calldataload(0x4)
      jump(tag_5)
    tag_4:
      stop
    tag_5:
        /* "contracts/ContractCreation.sol":422:440  knownValues.insert */
      linkerSymbol("c792e0a7f2dc81eac267450bf298a6ed51dab146debecfb69b118d9f48e7147a")
      0x831cb739
        /* "contracts/ContractCreation.sol":422:433  knownValues */
      0x1
        /* "contracts/ContractCreation.sol":441:446  value */
      dup4
        /* "contracts/ContractCreation.sol":422:447  knownValues.insert(value) */
      mstore(add(0x20, mload(0x40)), 0x0)
      mload(0x40)
      0x100000000000000000000000000000000000000000000000000000000
      0xffffffff
      dup6
      and
      mul
      dup2
      mstore
      0x4
      dup2
      add
      swap3
      swap1
      swap3
      mstore
      0x24
      dup3
      add
      mstore
      0x44
      add
      0x20
      mload(0x40)
      dup1
      dup4
      sub
      dup2
      dup7
      dup1
      extcodesize
      iszero
      iszero
      tag_7
      jumpi
      0x0
      dup1
      revert
    tag_7:
      sub(gas, 0x2c6)
      delegatecall
      iszero
      iszero
      tag_8
      jumpi
      0x0
      dup1
      revert
    tag_8:
      pop
      pop
      pop
      mload(0x40)
      dup1
      mload
      swap1
      pop
        /* "contracts/ContractCreation.sol":414:448  require(knownValues.insert(value)) */
      iszero
      iszero
      tag_9
      jumpi
      0x0
      dup1
      revert
    tag_9:
        /* "contracts/ContractCreation.sol":201:453  function register(uint value) {... */
    tag_6:
      pop
      jump	// out

    auxdata: 0xa165627a7a723058204fcd4b0c11c213cee7f309784b43b8253f4a5c2cc2d1fa1700fedb4f1a0623130029
}

sub_1: assembly {
        /* "contracts/ContractCreation.sol":457:472  contract C2 {... */
      mstore(0x40, 0x60)
      jumpi(tag_1, iszero(callvalue))
      0x0
      dup1
      revert
    tag_1:
    tag_2:
      dataSize(sub_0)
      dup1
      dataOffset(sub_0)
      0x0
      codecopy
      0x0
      return
    stop

    sub_0: assembly {
            /* "contracts/ContractCreation.sol":457:472  contract C2 {... */
          mstore(0x40, 0x60)
        tag_1:
          0x0
          dup1
          revert

        auxdata: 0xa165627a7a723058209bba8d8112e7296566df47ae75df062c528c3b4cd0befc39224ee63a02d90c510029
    }
}
Binary: 
6060604052341561000f57600080fd5b5b610018610054565b604051809103906000f080151561002e57600080fd5b60008054600160a060020a031916600160a060020a03929092169190911790555b610063565b60405160528061018e83390190565b61011c806100726000396000f300606060405263ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663f207564e8114603c575b600080fd5b3415604657600080fd5b604f6004356051565b005b73__contracts/Set.sol:Set_________________63831cb7396001836000604051602001526040517c010000000000000000000000000000000000000000000000000000000063ffffffff85160281526004810192909252602482015260440160206040518083038186803b151560c857600080fd5b6102c65a03f4151560d857600080fd5b50505060405180519050151560ec57600080fd5b5b505600a165627a7a723058204fcd4b0c11c213cee7f309784b43b8253f4a5c2cc2d1fa1700fedb4f1a062313002960606040523415600e57600080fd5b5b603680601c6000396000f30060606040525b600080fd00a165627a7a723058209bba8d8112e7296566df47ae75df062c528c3b4cd0befc39224ee63a02d90c510029

======= contracts/ContractCreation.sol:C2 =======
EVM assembly:
    /* "contracts/ContractCreation.sol":457:472  contract C2 {... */
  mstore(0x40, 0x60)
  jumpi(tag_1, iszero(callvalue))
  0x0
  dup1
  revert
tag_1:
tag_2:
  dataSize(sub_0)
  dup1
  dataOffset(sub_0)
  0x0
  codecopy
  0x0
  return
stop

sub_0: assembly {
        /* "contracts/ContractCreation.sol":457:472  contract C2 {... */
      mstore(0x40, 0x60)
    tag_1:
      0x0
      dup1
      revert

    auxdata: 0xa165627a7a723058209bba8d8112e7296566df47ae75df062c528c3b4cd0befc39224ee63a02d90c510029
}
Binary: 
60606040523415600e57600080fd5b5b603680601c6000396000f30060606040525b600080fd00a165627a7a723058209bba8d8112e7296566df47ae75df062c528c3b4cd0befc39224ee63a02d90c510029

======= contracts/Set.sol:Set =======
EVM assembly:
    /* "contracts/Set.sol":84:655  library Set {... */
  mstore(0x40, 0x60)
  jumpi(tag_1, iszero(callvalue))
  0x0
  dup1
  revert
tag_1:
tag_2:
  dataSize(sub_0)
  dup1
  dataOffset(sub_0)
  0x0
  codecopy
  0x0
  return
stop

sub_0: assembly {
        /* "contracts/Set.sol":84:655  library Set {... */
      mstore(0x40, 0x60)
      and(div(calldataload(0x0), 0x100000000000000000000000000000000000000000000000000000000), 0xffffffff)
      0x483b8a14
      dup2
      eq
      tag_2
      jumpi
      dup1
      0x6ce8e081
      eq
      tag_3
      jumpi
      dup1
      0x831cb739
      eq
      tag_4
      jumpi
    tag_1:
      0x0
      dup1
      revert
        /* "contracts/Set.sol":346:540  function remove(Data storage self, uint value)... */
    tag_2:
      tag_5
      calldataload(0x4)
      calldataload(0x24)
      jump(tag_6)
    tag_5:
      mload(0x40)
      swap1
      iszero
      iszero
      dup2
      mstore
      0x20
      add
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      return
        /* "contracts/Set.sol":544:653  function contains(Data storage self, uint value)... */
    tag_3:
      tag_5
      calldataload(0x4)
      calldataload(0x24)
      jump(tag_8)
    tag_7:
      mload(0x40)
      swap1
      iszero
      iszero
      dup2
      mstore
      0x20
      add
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      return
        /* "contracts/Set.sol":148:342  function insert(Data storage self, uint value)... */
    tag_4:
      tag_5
      calldataload(0x4)
      calldataload(0x24)
      jump(tag_10)
    tag_9:
      mload(0x40)
      swap1
      iszero
      iszero
      dup2
      mstore
      0x20
      add
      mload(0x40)
      dup1
      swap2
      sub
      swap1
      return
        /* "contracts/Set.sol":346:540  function remove(Data storage self, uint value)... */
    tag_6:
        /* "contracts/Set.sol":408:412  bool */
      0x0
        /* "contracts/Set.sol":429:446  self.flags[value] */
      dup2
      dup2
      mstore
      0x20
      dup4
      swap1
      mstore
      0x40
      dup2
      keccak256
      sload
      0xff
      and
        /* "contracts/Set.sol":428:446  !self.flags[value] */
      iszero
        /* "contracts/Set.sol":424:470  if (!self.flags[value])... */
      iszero
      tag_12
      jumpi
      pop
        /* "contracts/Set.sol":465:470  false */
      0x0
        /* "contracts/Set.sol":458:470  return false */
      jump(tag_11)
        /* "contracts/Set.sol":424:470  if (!self.flags[value])... */
    tag_12:
      pop
        /* "contracts/Set.sol":511:516  false */
      0x0
        /* "contracts/Set.sol":491:508  self.flags[value] */
      dup2
      dup2
      mstore
      0x20
      dup4
      swap1
      mstore
      0x40
      swap1
      keccak256
        /* "contracts/Set.sol":491:516  self.flags[value] = false */
      dup1
      sload
      not(0xff)
      and
      swap1
      sstore
      0x1
        /* "contracts/Set.sol":346:540  function remove(Data storage self, uint value)... */
    tag_11:
      swap3
      swap2
      pop
      pop
      jump	// out
        /* "contracts/Set.sol":544:653  function contains(Data storage self, uint value)... */
    tag_8:
        /* "contracts/Set.sol":608:612  bool */
      0x0
        /* "contracts/Set.sol":631:648  self.flags[value] */
      dup2
      dup2
      mstore
      0x20
      dup4
      swap1
      mstore
      0x40
      swap1
      keccak256
      sload
      0xff
      and
        /* "contracts/Set.sol":544:653  function contains(Data storage self, uint value)... */
    tag_13:
      swap3
      swap2
      pop
      pop
      jump	// out
        /* "contracts/Set.sol":148:342  function insert(Data storage self, uint value)... */
    tag_10:
        /* "contracts/Set.sol":210:214  bool */
      0x0
        /* "contracts/Set.sol":230:247  self.flags[value] */
      dup2
      dup2
      mstore
      0x20
      dup4
      swap1
      mstore
      0x40
      dup2
      keccak256
      sload
      0xff
      and
        /* "contracts/Set.sol":226:269  if (self.flags[value])... */
      iszero
      tag_15
      jumpi
      pop
        /* "contracts/Set.sol":264:269  false */
      0x0
        /* "contracts/Set.sol":257:269  return false */
      jump(tag_11)
        /* "contracts/Set.sol":226:269  if (self.flags[value])... */
    tag_15:
      pop
        /* "contracts/Set.sol":294:304  self.flags */
      0x0
        /* "contracts/Set.sol":294:311  self.flags[value] */
      dup2
      dup2
      mstore
      0x20
      dup4
      swap1
      mstore
      0x40
      swap1
      keccak256
        /* "contracts/Set.sol":294:318  self.flags[value] = true */
      dup1
      sload
      not(0xff)
      and
        /* "contracts/Set.sol":314:318  true */
      0x1
        /* "contracts/Set.sol":294:318  self.flags[value] = true */
      swap1
      dup2
      or
      swap1
      swap2
      sstore
        /* "contracts/Set.sol":148:342  function insert(Data storage self, uint value)... */
    tag_14:
      swap3
      swap2
      pop
      pop
      jump	// out

    auxdata: 0xa165627a7a72305820ab6966591b460129b1564b88ae42f8514f9f49ad295e19eeeb78a46dc7a73ec40029
}
Binary: 
6060604052341561000f57600080fd5b5b6101818061001f6000396000f300606060405263ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663483b8a1481146100535780636ce8e08114610075578063831cb73914610097575b600080fd5b6100616004356024356100b9565b604051901515815260200160405180910390f35b6100616004356024356100f9565b604051901515815260200160405180910390f35b610061600435602435610112565b604051901515815260200160405180910390f35b60008181526020839052604081205460ff1615156100d9575060006100f3565b506000818152602083905260409020805460ff1916905560015b92915050565b60008181526020839052604090205460ff165b92915050565b60008181526020839052604081205460ff1615610131575060006100f3565b506000818152602083905260409020805460ff191660019081179091555b929150505600a165627a7a72305820ab6966591b460129b1564b88ae42f8514f9f49ad295e19eeeb78a46dc7a73ec40029
