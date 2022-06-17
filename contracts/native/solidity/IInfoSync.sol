//SPDX-License-Identifier: Unlicense
pragma solidity ^0.5.0;
pragma experimental ABIEncoderV2;

interface InfoSync {
  event SyncRootInfoEvent(uint64 chainID, uint32 height, uint256 BlockHeight);
  function name() external view returns(string memory);
  function syncRootInfo(uint64 chainID, bytes[] calldata rootInfos) external returns(bool);
  function getLatestHeight(uint64 chainID) external returns(uint32);
  function getInfo(uint64 chainID, uint32 height) external returns(bytes memory);
}