// SPDX-License-Identifier: MIT
pragma solidity ^0.8.26;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/AccessControl.sol";

contract SkillSBT is ERC721, AccessControl {
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    struct SkillMeta {
        uint256 skillType;
        string metadataCID;
        uint64 issuedAt;
    }

    uint256 private _nextTokenId;
    mapping(uint256 => SkillMeta) public skillMetadata;
    mapping(address => mapping(uint256 => bool)) public hasSkillType;
    mapping(address => mapping(uint256 => uint256)) public skillTypeTokenId;

    event SkillMinted(address indexed to, uint256 indexed tokenId, uint256 skillType);

    constructor(address admin) ERC721("Web3Proof Skill SBT", "W3PSBT") {
        _grantRole(DEFAULT_ADMIN_ROLE, admin);
        _grantRole(MINTER_ROLE, admin);
    }

    function mint(
        address to,
        uint256 skillType,
        string calldata metadataCID
    ) external onlyRole(MINTER_ROLE) returns (uint256 tokenId) {
        require(!hasSkillType[to][skillType], "Skill type already minted");

        tokenId = ++_nextTokenId;
        _safeMint(to, tokenId);

        skillMetadata[tokenId] = SkillMeta({
            skillType: skillType,
            metadataCID: metadataCID,
            issuedAt: uint64(block.timestamp)
        });
        hasSkillType[to][skillType] = true;
        skillTypeTokenId[to][skillType] = tokenId;

        emit SkillMinted(to, tokenId, skillType);
    }

    function burn(uint256 tokenId) external onlyRole(MINTER_ROLE) {
        address owner = ownerOf(tokenId);
        uint256 skillType = skillMetadata[tokenId].skillType;
        hasSkillType[owner][skillType] = false;
        delete skillTypeTokenId[owner][skillType];
        delete skillMetadata[tokenId];
        _burn(tokenId);
    }

    function locked(uint256 /* tokenId */) external pure returns (bool) {
        return true;
    }

    function verify(address owner, uint256 skillType) external view returns (bool, uint256 tokenId) {
        if (!hasSkillType[owner][skillType]) {
            return (false, 0);
        }
        return (true, skillTypeTokenId[owner][skillType]);
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
