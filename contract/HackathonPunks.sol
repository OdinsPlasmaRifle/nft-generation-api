/**
 *Submitted for verification at Etherscan.io on 2021-08-27
*/

// SPDX-License-Identifier: MIT

pragma solidity ^0.8.0;

import "./ERC721Enumerable.sol";
import "./ERC721.sol";


contract HackathonPunks is ERC721Enumerable {
    
    uint256 blocknum;
    uint numAttributes;
    uint numPerAttribute;
    string description;
    
    string[] private accessories = [
    	"earring-gold",
    	"earring-silver",
    	"horns-pink"
    ];

    string[] private background = [
    	"blue",
    	"pink",
    	"purple"
    ];
    
    string[] private eyes = [
    	"blue",
    	"brown",
    	"green"
    ];
    
    string[] private hair = [
    	"blonde",
    	"blue",
    	"purple"
    ];
    
    string[] private lips = [
    	"orange",
    	"pink",
    	"red"
    ];
    
    string[] private outfits = [
    	"1",
    	"2",
    	"3"
    ];
    
    string[] private skin = [
    	"F-1",
    	"F-2",
    	"F-3"
    ];
    
    function getAttributes(uint256 tokenId) public view returns (string[7] memory) {
        uint[7] memory set = _generateRandomNumberSet(3, 6, tokenId);
        return [
            accessories[set[0]],
            background[set[1]],
            eyes[set[2]],
            hair[set[3]],
            lips[set[4]],
            outfits[set[5]],
            skin[set[6]]    
        ];
    }
    
    function _generateRandomNumberSet(uint mod, uint iterations, uint extra_entropy) internal view returns(uint[7] memory) {
        /*
        Number generator that creates a set of random numbers.
        */
        uint[7] memory numset;
        bytes32 blockHash = blockhash(blocknum-1);
        for (uint index = 0; index < iterations; index++) {
            numset[index] = (uint256(keccak256(abi.encodePacked(blockHash))) + extra_entropy) % mod;
            blockHash = blockhash(blocknum-(index+1));
        }
        return numset;
    }
    
    constructor(string memory n, string memory desc, string memory sym, uint attributes, uint perAtribute) ERC721(n, sym) {
        blocknum = block.number;
        description = desc;
        numAttributes = attributes;
        numPerAttribute = perAtribute;
    }
}