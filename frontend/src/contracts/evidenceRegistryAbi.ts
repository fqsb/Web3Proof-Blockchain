export const evidenceRegistryAbi = [
  {
    type: "function",
    name: "createEvidence",
    stateMutability: "nonpayable",
    inputs: [
      { name: "evidenceNoHash", type: "bytes32" },
      { name: "fileHash", type: "bytes32" },
      { name: "metadataURI", type: "string" },
    ],
    outputs: [{ name: "evidenceId", type: "uint256" }],
  },
  {
    type: "event",
    name: "EvidenceCreated",
    inputs: [
      { name: "evidenceId", type: "uint256", indexed: true },
      { name: "owner", type: "address", indexed: true },
      { name: "fileHash", type: "bytes32", indexed: true },
      { name: "evidenceNoHash", type: "bytes32", indexed: false },
      { name: "metadataURI", type: "string", indexed: false },
    ],
  },
] as const;
