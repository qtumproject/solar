pragma solidity ^0.5.8;

contract AB {
  uint256 a;
  int256 b;

  uint256 fund;

  constructor(uint256 _a, int256 _b) public {
    a = _a;
    b = _b;
  }

  function burnFund() public payable {
    fund += msg.value;
  }

  function getBurnedFund() public view returns(uint256) {
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

  function getA() public view returns(uint256) {
    return a;
  }

  function getB() public view returns(int256) {
    return b;
  }

  function getAB() public view returns(uint256, int256) {
    return (a, b);
  }
}