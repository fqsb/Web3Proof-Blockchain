export const credentialSbtAbi = [
  {
    type: "function",
    name: "mintCredential",
    stateMutability: "nonpayable",
    inputs: [
      { name: "to", type: "address" },
      { name: "evidenceId", type: "uint256" },
      { name: "tokenURIValue", type: "string" },
    ],
    outputs: [{ name: "tokenId", type: "uint256" }],
  },
  {
    type: "event",
    name: "CredentialMinted",
    inputs: [
      { name: "to", type: "address", indexed: true },
      { name: "tokenId", type: "uint256", indexed: true },
      { name: "evidenceId", type: "uint256", indexed: true },
      { name: "tokenURIValue", type: "string", indexed: false },
    ],
  },
] as const;
