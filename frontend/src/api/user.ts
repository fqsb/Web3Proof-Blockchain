import client, { ApiResponse } from "./client";
import { User } from "./auth";

export async function prepareDID() {
  const res = await client.post<ApiResponse<unknown>>("/users/did/prepare");
  return res.data.data as {
    contract_address: string;
    did: string;
    github: string;
    metadata_cid: string;
    chain_id: number;
  };
}

export async function confirmDID(txHash: string) {
  const res = await client.post<ApiResponse<User>>("/users/did/confirm", { tx_hash: txHash });
  return res.data.data;
}

export async function updateProfile(data: Partial<User>) {
  const res = await client.put<ApiResponse<User>>("/users/profile", data);
  return res.data.data;
}

export async function registerEnterprise(data: {
  company_name: string;
  industry?: string;
  contact_email?: string;
  website?: string;
}) {
  const res = await client.put<ApiResponse<User>>("/users/role/enterprise", data);
  return res.data.data;
}

export async function getVerifyReport(address: string) {
  const res = await client.get<ApiResponse<unknown>>(`/verify/${address}`);
  return res.data.data as {
    report_id: number;
    identity_verified: boolean;
    did?: string;
    target_user_id: number;
    target_address: string;
    viewer_id: number;
    project_count: number;
    sbt_count: number;
    reputation_score: number;
    grade: string;
    skills: string[];
    verified_at: string;
  };
}

export async function listVerifyReports() {
  const res = await client.get<ApiResponse<unknown[]>>("/verify/reports");
  return res.data.data as {
    id: number;
    target_user_id: number;
    viewer_id: number;
    report: Awaited<ReturnType<typeof getVerifyReport>>;
    created_at: string;
  }[];
}

export interface ResumeData {
  user: User;
  projects: unknown[];
  sbts: { id: number; token_id: number; skill: { name: string }; tx_hash: string }[];
  reputation: { total_score: number; grade: string; project_score: number; cert_score: number; activity_score: number };
}

export async function getResume(address: string) {
  const res = await client.get<ApiResponse<ResumeData>>(`/resume/${address}`);
  return res.data.data;
}
