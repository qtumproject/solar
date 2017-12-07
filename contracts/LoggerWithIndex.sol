pragma solidity ^0.4.11;

contract LoggerWithIndex {
  event EventLog(uint256 a, uint256 indexed b, uint256 c);
  function log(uint256 a, uint256 b, uint256 c) {
    EventLog(a, b, c);
  }
}