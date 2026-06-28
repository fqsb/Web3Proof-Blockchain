import client, { ApiResponse } from "./client";

export async function verifyFile(file: File) {
  const form = new FormData();
  form.append("file", file);
  const res = await client.post<ApiResponse<unknown>>("/verify/file", form);
  return res.data.data;
}

export async function verifyEvidence(no: string) {
  const res = await client.get<ApiResponse<unknown>>(`/verify/evidence/${encodeURIComponent(no)}`);
  return res.data.data;
}

export async function verifyCertificate(no: string) {
  const res = await client.get<ApiResponse<unknown>>(`/verify/certificate/${encodeURIComponent(no)}`);
  return res.data.data;
}

export async function verifyWallet(address: string) {
  const normalized = address.trim().toLowerCase();
  const res = await client.get<ApiResponse<unknown>>(`/verify/wallet/${encodeURIComponent(normalized)}`);
  return res.data.data;
}

export async function listReports() {
  const res = await client.get<ApiResponse<unknown[]>>("/verifier/reports");
  return res.data.data;
}
