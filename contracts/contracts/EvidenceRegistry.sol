// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/access/AccessControl.sol";

contract EvidenceRegistry is AccessControl {
    bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");

    struct Evidence {
        uint256 id;
        bytes32 evidenceNoHash;
        bytes32 fileHash;
        address owner;
        string metadataURI;
        uint64 createdAt;
        bool revoked;
        bool exists;
    }

    uint256 public evidenceCount;
    mapping(uint256 => Evidence) private _evidences;
    mapping(bytes32 => uint256) private _evidenceIdByFileHash;
    mapping(address => uint256[]) private _ownerEvidenceIds;

    event EvidenceCreated(
        uint256 indexed evidenceId,
        address indexed owner,
        bytes32 indexed fileHash,
        bytes32 evidenceNoHash,
        string metadataURI
    );
    event EvidenceRevoked(uint256 indexed evidenceId, address indexed owner);

    constructor(address admin) {
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(OPERATOR_ROLE, admin);
    }

    function createEvidence(
        bytes32 evidenceNoHash,
        bytes32 fileHash,
        string calldata metadataURI
    ) external returns (uint256 evidenceId) {
        require(evidenceNoHash != bytes32(0), "Evidence no hash required");
        require(fileHash != bytes32(0), "File hash required");
        require(_evidenceIdByFileHash[fileHash] == 0, "File hash already exists");

        evidenceId = ++evidenceCount;
        _evidences[evidenceId] = Evidence({
            id: evidenceId,
            evidenceNoHash: evidenceNoHash,
            fileHash: fileHash,
            owner: msg.sender,
            metadataURI: metadataURI,
            createdAt: uint64(block.timestamp),
            revoked: false,
            exists: true
        });
        _evidenceIdByFileHash[fileHash] = evidenceId;
        _ownerEvidenceIds[msg.sender].push(evidenceId);

        emit EvidenceCreated(evidenceId, msg.sender, fileHash, evidenceNoHash, metadataURI);
    }

    function getEvidence(uint256 evidenceId) external view returns (Evidence memory) {
        require(_evidences[evidenceId].exists, "Evidence not found");
        return _evidences[evidenceId];
    }

    function getEvidenceByHash(bytes32 fileHash) external view returns (Evidence memory) {
        uint256 evidenceId = _evidenceIdByFileHash[fileHash];
        require(evidenceId != 0, "Evidence not found");
        return _evidences[evidenceId];
    }

    function getEvidenceIdsByOwner(address owner) external view returns (uint256[] memory) {
        return _ownerEvidenceIds[owner];
    }

    function verifyEvidence(bytes32 fileHash) external view returns (bool valid, uint256 evidenceId, address owner) {
        evidenceId = _evidenceIdByFileHash[fileHash];
        if (evidenceId == 0) {
            return (false, 0, address(0));
        }
        Evidence memory evidence = _evidences[evidenceId];
        return (!evidence.revoked, evidenceId, evidence.owner);
    }

    function revokeEvidence(uint256 evidenceId) external {
        Evidence storage evidence = _evidences[evidenceId];
        require(evidence.exists, "Evidence not found");
        require(
            msg.sender == evidence.owner || hasRole(OPERATOR_ROLE, msg.sender),
            "Not authorized"
        );
        evidence.revoked = true;
        emit EvidenceRevoked(evidenceId, evidence.owner);
    }
}
