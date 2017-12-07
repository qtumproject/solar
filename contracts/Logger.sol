pragma solidity ^0.4.11;

contract Logger {
  event EventLog(uint256 a, uint256 b);
  function log(uint256 a, uint256 b) {
    EventLog(a, b);
  }
}