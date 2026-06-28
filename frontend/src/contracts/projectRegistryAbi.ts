export const projectRegistryAbi = [
  "event ProjectAdded(uint256 indexed projectId, address indexed owner, string ipfsCID, bytes32 contentHash)",
  "function addProject(string name, string ipfsCID, bytes32 contentHash, string githubUrl, address contractAddr) returns (uint256 projectId)",
  "function getProject(uint256 projectId) view returns (tuple(address owner, string name, string ipfsCID, bytes32 contentHash, string githubUrl, address contractAddr, uint64 createdAt, bool exists))",
  "function projectCount() view returns (uint256)",
] as const;
