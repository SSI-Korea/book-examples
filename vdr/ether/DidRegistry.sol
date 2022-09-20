// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.7.0 <0.9.0;

contract SSIKRDid {
    mapping (string => string) _mapDids;

    function CreateDid(string memory _did, string memory _document) public {
        _mapDids[_did] = _document;
    }

    function ResolveDid(string memory _did) public view returns (string memory){
        return _mapDids[_did];
    }
}