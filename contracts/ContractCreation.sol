pragma solidity ^0.5.1;

import "./Set.sol";

contract C1 {
  address c2;

  using Set for Set.Data; // this is the crucial change
  Set.Data knownValues;

  constructor() {
    c2 = new C2();
  }

  function register(uint value) public {
    // Here, all variables of type Set.Data have
    // corresponding member functions.
    // The following function call is identical to
    // Set.insert(knownValues, value)
    require(knownValues.insert(value));
  }
}

contract C2 {
}