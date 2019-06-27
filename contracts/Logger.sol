pragma solidity ^0.5.1;

contract Logger {
  event EventLog(uint256 a, uint256 b);
  function log(uint256 a, uint256 b) public {
    emit EventLog(a, b);
  }
}