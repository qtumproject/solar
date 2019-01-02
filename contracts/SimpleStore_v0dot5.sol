pragma solidity ^0.5.2;

contract SimpleStore_v0dot5 {
  constructor(uint _value) public {
    value = _value;
  }

    function set(uint newValue) public {
        value = newValue;
    }

    function get() public view returns (uint) {
        return value;
    }

    uint value;
}
