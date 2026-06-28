export const didRegistryAbi = [
  "event DIDRegistered(address indexed owner, string did, string metadataCID)",
  "function registerDID(string did, string github, string metadataCID)",
  "function updateProfile(string github, string metadataCID)",
  "function getProfile(address owner) view returns (tuple(string did, string github, string metadataCID, uint64 updatedAt, bool exists))",
  "function hasProfile(address owner) view returns (bool)",
] as const;
