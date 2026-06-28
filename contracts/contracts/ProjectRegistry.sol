// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

contract ProjectRegistry {
    struct Project {
        address owner;
        string name;
        string ipfsCID;
        bytes32 contentHash;
        string githubUrl;
        address contractAddr;
        uint64 createdAt;
        bool exists;
    }

    uint256 public projectCount;
    mapping(uint256 => Project) private _projects;
    mapping(address => uint256[]) private _ownerProjects;

    event ProjectAdded(
        uint256 indexed projectId,
        address indexed owner,
        string ipfsCID,
        bytes32 contentHash
    );

    function addProject(
        string calldata name,
        string calldata ipfsCID,
        bytes32 contentHash,
        string calldata githubUrl,
        address contractAddr
    ) external returns (uint256 projectId) {
        require(bytes(name).length > 0, "Name required");
        require(bytes(ipfsCID).length > 0, "IPFS CID required");

        projectCount++;
        projectId = projectCount;

        _projects[projectId] = Project({
            owner: msg.sender,
            name: name,
            ipfsCID: ipfsCID,
            contentHash: contentHash,
            githubUrl: githubUrl,
            contractAddr: contractAddr,
            createdAt: uint64(block.timestamp),
            exists: true
        });

        _ownerProjects[msg.sender].push(projectId);

        emit ProjectAdded(projectId, msg.sender, ipfsCID, contentHash);
    }

    function getProject(uint256 projectId) external view returns (Project memory) {
        require(_projects[projectId].exists, "Project not found");
        return _projects[projectId];
    }

    function getProjectsByOwner(address owner) external view returns (uint256[] memory) {
        return _ownerProjects[owner];
    }

    function verifyProject(uint256 projectId, bytes32 expectedHash) external view returns (bool) {
        require(_projects[projectId].exists, "Project not found");
        return _projects[projectId].contentHash == expectedHash;
    }
}
