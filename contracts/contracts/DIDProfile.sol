// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract DIDProfile {
    struct Profile {
        string did;
        string metadataURI;
        uint64 updatedAt;
        bool exists;
    }

    mapping(address => Profile) private _profiles;

    event ProfileRegistered(address indexed owner, string did, string metadataURI);
    event ProfileUpdated(address indexed owner, string metadataURI);

    function registerProfile(string calldata did, string calldata metadataURI) external {
        require(!_profiles[msg.sender].exists, "Profile already registered");
        require(bytes(did).length > 0, "DID required");

        _profiles[msg.sender] = Profile({
            did: did,
            metadataURI: metadataURI,
            updatedAt: uint64(block.timestamp),
            exists: true
        });

        emit ProfileRegistered(msg.sender, did, metadataURI);
    }

    function updateProfile(string calldata metadataURI) external {
        require(_profiles[msg.sender].exists, "Profile not registered");
        Profile storage profile = _profiles[msg.sender];
        profile.metadataURI = metadataURI;
        profile.updatedAt = uint64(block.timestamp);
        emit ProfileUpdated(msg.sender, metadataURI);
    }

    function getProfile(address owner) external view returns (Profile memory) {
        require(_profiles[owner].exists, "Profile not found");
        return _profiles[owner];
    }

    function hasProfile(address owner) external view returns (bool) {
        return _profiles[owner].exists;
    }
}
