pragma solidity ^0.4.18;

import "./TestLibA.sol";
import "./TestLibB.sol";

contract TestLinkLib {
  function a() public view returns (uint256) {
    return TestLibA.a();
  }

  function b() public view returns (uint256) {
    return TestLibB.b();
  }
}

/*
solc --allow-paths=. --metadata --optimize --bin contracts/TestLinkLib.sol \
  --libraries 'contracts/TestLibA.sol:TestLibA:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa,contracts/TestLibB.sol:TestLibB:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' | tee TestLinkLib.out
*/