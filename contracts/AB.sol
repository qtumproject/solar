pragma solidity ^0.4.11;

contract AB {
  uint256 a;
  int256 b;

  uint256 fund;

  function AB(uint256 _a, int256 _b) public {
    a = _a;
    b = _b;
  }

  function burnFund() public payable {
    fund += msg.value;
  }

  function getBurnedFund() public constant returns(uint256) {
    return fund;
  }

  function setA(uint256 _a) public {
    a = _a;
  }

  function setB(int256 _b) public {
    b = _b;
  }

  function setAB(uint256 _a, int256 _b) public {
    a = _a;
    b = _b;
  }

  function getA() public constant returns(uint256) {
    return a;
  }

  function getB() public constant returns(int256) {
    return b;
  }

  function getAB() public constant returns(uint256, int256) {
    return (a, b);
  }
}