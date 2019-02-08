/*
 * Copyright (c) 2018 vChain, Inc. All Rights Reserved.
 * This software is released under GPL3.
 * The full license information can be found under:
 * https://www.gnu.org/licenses/gpl-3.0.en.html
 */
 
pragma solidity ^0.4.24;

contract Proof {


  struct FileDetails {
      uint timestamp;
      string owner;
  }

  mapping (string => FileDetails) files;

  event logFileAddedStatus(bool status, uint timestamp, string owner, string fileHash);

  function set(string owner, string fileHash) public {
      if (files[fileHash].timestamp == 0) {
          files[fileHash] = FileDetails(block.timestamp, owner);

          emit logFileAddedStatus(true, block.timestamp, owner, fileHash);
      } else {
          emit logFileAddedStatus(false, block.timestamp, owner, fileHash);
      }
  }

  function get(string fileHash) view public returns (uint timestamp, string owner) {
      return (files[fileHash].timestamp, files[fileHash].owner);

  }

  function get_timestamp() public view returns (uint timestamp) {
      return block.timestamp;
  }

}