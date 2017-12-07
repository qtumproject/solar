pragma solidity ^0.4.11;

contract Foo {
  uint256 a;
  uint256 b;

  function Foo() {
  }

  function foo(uint256 _a, uint256 _b) {
    a = _a;
    b = _b;
  }
}