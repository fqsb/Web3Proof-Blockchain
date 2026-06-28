import client, { ApiResponse } from "./client";

export interface Work {
  id: number;
  title: string;
  description?: string;
  category_id?: number;
  external_url?: string;
  visibility: string;
  status: string;
  created_at: string;
}

export interface WorkFile {
  id: number;
  work_id: number;
  file_name: string;
  storage_url?: string;
  file_size: number;
  sha256_hash: string;
}

export interface EvidencePrepareResult {
  contract_address: string;
  evidence_no: string;
  evidence_no_hash: string;
  file_hash: string;
  metadata_uri: string;
  chain_id: number;
}

export interface EvidenceRecord {
  id: number;
  work_id: number;
  evidence_no: string;
  file_hash: string;
  owner_address: string;
  chain_evidence_id?: number;
  tx_hash?: string;
  status: string;
  created_at: string;
}

export interface Certificate {
  id: number;
  evidence_id: number;
  certificate_no: string;
  pdf_storage_key: string;
  verify_url: string;
  created_at: string;
}

export async function listWorks() {
  const res = await client.get<ApiResponse<Work[]>>("/works");
  return res.data.data;
}

export async function createWork(data: Partial<Work>) {
  const res = await client.post<ApiResponse<Work>>("/works", data);
  return res.data.data;
}

export async function getWork(id: string | number) {
  const res = await client.get<ApiResponse<{ work: Work; files: WorkFile[]; evidences: EvidenceRecord[]; certificates: Certificate[] }>>(`/works/${id}`);
  return res.data.data;
}

export async function listMyEvidence() {
  const res = await client.get<ApiResponse<EvidenceRecord[]>>("/evidence/my");
  return res.data.data;
}

export async function listMyCertificates() {
  const res = await client.get<ApiResponse<Certificate[]>>("/certificates/my");
  return res.data.data;
}

export async function uploadWorkFile(workId: string | number, file: File) {
  const form = new FormData();
  form.append("file", file);
  const res = await client.post<ApiResponse<WorkFile>>(`/works/${workId}/files`, form);
  return res.data.data;
}

export async function prepareEvidence(workId: string | number) {
  const res = await client.post<ApiResponse<EvidencePrepareResult>>(`/works/${workId}/evidence/prepare`);
  return res.data.data;
}

export async function confirmEvidence(workId: string | number, tx_hash: string, chain_evidence_id: number) {
  const res = await client.post<ApiResponse<unknown>>(`/works/${workId}/evidence/confirm`, { tx_hash, chain_evidence_id });
  return res.data.data;
}

export async function generateCertificate(evidence_id: number) {
  const res = await client.post<ApiResponse<unknown>>("/certificates/generate", { evidence_id });
  return res.data.data;
}

export async function listCategories() {
  const res = await client.get<ApiResponse<Array<{ id: number; name: string; code: string }>>>("/categories");
  return res.data.data;
}
