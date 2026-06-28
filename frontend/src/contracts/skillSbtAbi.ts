export const skillSbtAbi = [
  "event SkillMinted(address indexed to, uint256 indexed tokenId, uint256 skillType)",
  "function mint(address to, uint256 skillType, string metadataCID) returns (uint256)",
  "function verify(address owner, uint256 skillType) view returns (bool, uint256)",
  "function locked(uint256 tokenId) view returns (bool)",
] as const;
