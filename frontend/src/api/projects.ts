import client, { ApiResponse } from "./client";

export interface Project {
  id: number;
  user_id: number;
  name: string;
  description?: string;
  github_url?: string;
  contract_address?: string;
  ipfs_cid?: string;
  ipfs_url?: string;
  content_hash?: string;
  chain_project_id?: number;
  tx_hash?: string;
  status: "draft" | "pending_chain" | "confirmed" | "failed";
  created_at: string;
  updated_at: string;
}

export interface PrepareChainResult {
  contract_address: string;
  name: string;
  ipfs_cid: string;
  content_hash: string;
  github_url: string;
  contract_addr: string;
  chain_id: number;
}

export async function createProject(data: {
  name: string;
  description?: string;
  github_url?: string;
  contract_address?: string;
}) {
  const res = await client.post<ApiResponse<Project>>("/projects", data);
  return res.data.data;
}

export async function listProjects() {
  const res = await client.get<ApiResponse<Project[]>>("/projects");
  return res.data.data;
}

export async function prepareChain(projectId: number) {
  const res = await client.post<ApiResponse<PrepareChainResult>>(`/projects/${projectId}/prepare-chain`);
  return res.data.data;
}

export async function confirmChain(projectId: number, txHash: string, chainProjectId: number) {
  const res = await client.post<ApiResponse<Project>>(`/projects/${projectId}/confirm-chain`, {
    tx_hash: txHash,
    chain_project_id: chainProjectId,
  });
  return res.data.data;
}
