pragma solidity ^0.4.11;

contract A {
  uint256 a;

  function setA(uint256 _a) {
    a = _a;
  }

  function getA() constant returns (uint256) {
    return a;
  }
}
