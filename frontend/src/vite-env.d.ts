/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string;
  readonly VITE_CHAIN_ID: string;
  readonly VITE_DID_REGISTRY_ADDRESS: string;
  readonly VITE_PROJECT_REGISTRY_ADDRESS: string;
  readonly VITE_SKILL_SBT_ADDRESS: string;
  readonly VITE_REPUTATION_ADDRESS: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
