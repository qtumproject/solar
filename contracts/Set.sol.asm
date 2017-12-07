
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

    auxdata: 0xa165627a7a72305820209b4f10b5f95986503880077a003580f1e2b9db9beb89b6758c92961ccd47dc0029
}
Binary: 
6060604052341561000f57600080fd5b5b6101818061001f6000396000f300606060405263ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663483b8a1481146100535780636ce8e08114610075578063831cb73914610097575b600080fd5b6100616004356024356100b9565b604051901515815260200160405180910390f35b6100616004356024356100f9565b604051901515815260200160405180910390f35b610061600435602435610112565b604051901515815260200160405180910390f35b60008181526020839052604081205460ff1615156100d9575060006100f3565b506000818152602083905260409020805460ff1916905560015b92915050565b60008181526020839052604090205460ff165b92915050565b60008181526020839052604081205460ff1615610131575060006100f3565b506000818152602083905260409020805460ff191660019081179091555b929150505600a165627a7a72305820209b4f10b5f95986503880077a003580f1e2b9db9beb89b6758c92961ccd47dc0029
