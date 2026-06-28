// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract DIDRegistry {
    struct Profile {
        string did;
        string github;
        string metadataCID;
        uint64 updatedAt;
        bool exists;
    }

    mapping(address => Profile) private _profiles;

    event DIDRegistered(address indexed owner, string did, string metadataCID);
    event ProfileUpdated(address indexed owner, string metadataCID);

    function registerDID(
        string calldata did,
        string calldata github,
        string calldata metadataCID
    ) external {
        require(!_profiles[msg.sender].exists, "DID already registered");
        require(bytes(did).length > 0, "DID required");

        _profiles[msg.sender] = Profile({
            did: did,
            github: github,
            metadataCID: metadataCID,
            updatedAt: uint64(block.timestamp),
            exists: true
        });

        emit DIDRegistered(msg.sender, did, metadataCID);
    }

    function updateProfile(string calldata github, string calldata metadataCID) external {
        require(_profiles[msg.sender].exists, "DID not registered");

        Profile storage profile = _profiles[msg.sender];
        profile.github = github;
        profile.metadataCID = metadataCID;
        profile.updatedAt = uint64(block.timestamp);

        emit ProfileUpdated(msg.sender, metadataCID);
    }

    function getProfile(address owner) external view returns (Profile memory) {
        require(_profiles[owner].exists, "Profile not found");
        return _profiles[owner];
    }

    function hasProfile(address owner) external view returns (bool) {
        return _profiles[owner].exists;
    }
}
