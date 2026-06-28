import client, { ApiResponse } from "./client";

export async function applyCertification(data: { work_id: number; evidence_id: number; materials_desc?: string }) {
  const res = await client.post<ApiResponse<unknown>>("/applications", data);
  return res.data.data;
}

export async function listMyApplications() {
  const res = await client.get<ApiResponse<unknown[]>>("/applications/my");
  return res.data.data;
}

export async function listAuditApplications() {
  const res = await client.get<ApiResponse<unknown[]>>("/auditor/applications");
  return res.data.data;
}

export async function reviewApplication(id: number, status: "approved" | "rejected", review_note?: string) {
  const res = await client.put<ApiResponse<unknown>>(`/auditor/applications/${id}/review`, { status, review_note });
  return res.data.data;
}

export async function prepareCredentialMint(id: number) {
  const res = await client.post<ApiResponse<{ contract_address: string; to_address: string; evidence_id: number; token_uri: string; chain_id: number }>>(`/auditor/applications/${id}/sbt/prepare`);
  return res.data.data;
}

export async function confirmCredentialMint(id: number, data: { tx_hash: string; token_id: number; token_uri: string }) {
  const res = await client.post<ApiResponse<unknown>>(`/auditor/applications/${id}/sbt/mint`, data);
  return res.data.data;
}
