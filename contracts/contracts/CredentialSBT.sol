// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract CredentialSBT is ERC721, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    struct CredentialMeta {
        uint256 evidenceId;
        string tokenURIValue;
        uint64 issuedAt;
        bool revoked;
    }

    uint256 private _nextTokenId;
    mapping(uint256 => CredentialMeta) public credentialMetadata;
    mapping(address => mapping(uint256 => uint256)) public credentialOfEvidence;

    event CredentialMinted(address indexed to, uint256 indexed tokenId, uint256 indexed evidenceId, string tokenURIValue);
    event CredentialRevoked(uint256 indexed tokenId, address indexed owner);

    constructor(address admin) ERC721("Web3Proof Credential SBT", "W3PCRED") {
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(MINTER_ROLE, admin);
    }

    function mintCredential(
        address to,
        uint256 evidenceId,
        string calldata tokenURIValue
    ) external onlyRole(MINTER_ROLE) returns (uint256 tokenId) {
        require(to != address(0), "Invalid recipient");
        require(evidenceId != 0, "Evidence required");
        require(credentialOfEvidence[to][evidenceId] == 0, "Credential already exists");

        tokenId = ++_nextTokenId;
        _safeMint(to, tokenId);

        credentialMetadata[tokenId] = CredentialMeta({
            evidenceId: evidenceId,
            tokenURIValue: tokenURIValue,
            issuedAt: uint64(block.timestamp),
            revoked: false
        });
        credentialOfEvidence[to][evidenceId] = tokenId;

        emit CredentialMinted(to, tokenId, evidenceId, tokenURIValue);
    }

    function tokenURI(uint256 tokenId) public view override returns (string memory) {
        _requireOwned(tokenId);
        return credentialMetadata[tokenId].tokenURIValue;
    }

    function revokeCredential(uint256 tokenId) external onlyRole(MINTER_ROLE) {
        address owner = ownerOf(tokenId);
        credentialMetadata[tokenId].revoked = true;
        emit CredentialRevoked(tokenId, owner);
    }

    function hasCredential(address owner, uint256 evidenceId) external view returns (bool, uint256 tokenId) {
        tokenId = credentialOfEvidence[owner][evidenceId];
        if (tokenId == 0 || credentialMetadata[tokenId].revoked) {
            return (false, 0);
        }
        return (true, tokenId);
    }

    function locked(uint256) external pure returns (bool) {
        return true;
    }

    function _update(address to, uint256 tokenId, address auth) internal override returns (address) {
        address from = _ownerOf(tokenId);
        if (from != address(0) && to != address(0)) {
            revert("SBT: non-transferable");
        }
        return super._update(to, tokenId, auth);
    }

    function supportsInterface(bytes4 interfaceId) public view override(ERC721, AccessControl) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}
