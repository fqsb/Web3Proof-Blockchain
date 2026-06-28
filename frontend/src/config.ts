export const API_BASE = import.meta.env.VITE_API_BASE_URL || "/api/v1";
export const CHAIN_ID = Number(import.meta.env.VITE_CHAIN_ID || 11155111);

export const BLOCK_EXPLORER =
  CHAIN_ID === 11155111 ? "https://sepolia.etherscan.io" : "https://etherscan.io";

export const CONTRACTS = {
  didProfile: import.meta.env.VITE_DID_PROFILE_ADDRESS || import.meta.env.VITE_DID_REGISTRY_ADDRESS || "",
  evidenceRegistry: import.meta.env.VITE_EVIDENCE_REGISTRY_ADDRESS || import.meta.env.VITE_PROJECT_REGISTRY_ADDRESS || "",
  credentialSBT: import.meta.env.VITE_CREDENTIAL_SBT_ADDRESS || import.meta.env.VITE_SKILL_SBT_ADDRESS || "",
  reputation: import.meta.env.VITE_REPUTATION_ADDRESS || "",
};
