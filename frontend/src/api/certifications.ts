import client, { ApiResponse } from "./client";

export interface Skill {
  id: number;
  code: string;
  name: string;
}

export interface CertificationApplication {
  id: number;
  skill_id: number;
  materials_cid: string;
  materials_desc?: string;
  status: "pending" | "approved" | "rejected" | "minting" | "minted" | "revoked";
  review_note?: string;
  skill?: Skill;
  user?: { wallet_address: string; nickname?: string };
}

export async function applyCertification(data: {
  skill_id: number;
  materials_desc?: string;
  links?: string[];
}) {
  const res = await client.post<ApiResponse<CertificationApplication>>("/certifications/apply", data);
  return res.data.data;
}

export async function listMyCertifications() {
  const res = await client.get<ApiResponse<CertificationApplication[]>>("/certifications/my");
  return res.data.data;
}

export async function listPendingCertifications() {
  const res = await client.get<ApiResponse<CertificationApplication[]>>("/admin/certifications");
  return res.data.data;
}

export async function reviewCertification(id: number, status: "approved" | "rejected", review_note?: string) {
  const res = await client.put<ApiResponse<CertificationApplication>>(`/admin/certifications/${id}`, {
    status,
    review_note,
  });
  return res.data.data;
}

export async function prepareMint(id: number) {
  const res = await client.post<ApiResponse<unknown>>(`/admin/certifications/${id}/mint/prepare`);
  return res.data.data as {
    contract_address: string;
    to_address: string;
    skill_type: number;
    metadata_cid: string;
    chain_id: number;
  };
}

export async function confirmMint(id: number, txHash: string, tokenId: number, metadataCID: string) {
  const res = await client.post<ApiResponse<unknown>>(`/admin/certifications/${id}/mint/confirm`, {
    tx_hash: txHash,
    token_id: tokenId,
    metadata_cid: metadataCID,
  });
  return res.data.data;
}

export async function listSkills() {
  const res = await client.get<ApiResponse<Skill[]>>("/skills");
  return res.data.data;
}
