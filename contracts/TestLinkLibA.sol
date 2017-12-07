pragma solidity ^0.4.18;

import "./TestLibA.sol";

contract TestLinkLibA {
  function a() public view returns (uint256) {
    return TestLibA.a();
  }
}

/*
solc --allow-paths=. --metadata --optimize --bin contracts/TestLinkLib.sol \
  --libraries 'contracts/TestLibA.sol:TestLibA:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa,contracts/TestLibB.sol:TestLibB:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb' | tee TestLinkLib.out
*/