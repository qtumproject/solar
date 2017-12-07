pragma solidity ^0.4.11;

contract Bar {
  uint256 a;
  uint256 b;

  function Bar() {
  }

  function foo(uint256 _a, uint256 _b) {
    a = _a;
    b = _b;
  }
}