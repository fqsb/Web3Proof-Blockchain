import client, { ApiResponse } from "./client";

export type RoleCode = "creator" | "verifier" | "auditor" | "admin";

export interface User {
  id: number;
  wallet_address: string;
  did?: string;
  nickname?: string;
  avatar_url?: string;
  bio?: string;
  email?: string;
  github?: string;
  active_role: RoleCode;
  roles: RoleCode[];
  is_did_registered: boolean;
}

export async function getNonce(address: string) {
  const res = await client.get<ApiResponse<{ message: string }>>("/auth/nonce", { params: { address } });
  return res.data.data.message;
}

export async function login(address: string, signature: string, message: string) {
  const res = await client.post<ApiResponse<{ token: string; user: User }>>("/auth/login", { address, signature, message });
  return res.data.data;
}

export async function getMe() {
  const res = await client.get<ApiResponse<User>>("/auth/me");
  return res.data.data;
}

export async function updateProfile(data: Partial<User>) {
  const res = await client.put<ApiResponse<User>>("/users/profile", data);
  return res.data.data;
}

export async function switchRole(role_code: RoleCode) {
  const res = await client.put<ApiResponse<{ token: string; user: User }>>("/users/current-role", { role_code });
  return res.data.data;
}

export async function requestVerifierRole(data: { org_name: string; industry?: string; contact_email?: string; website?: string }) {
  const res = await client.post<ApiResponse<User>>("/users/verifier-profile", data);
  return res.data.data;
}

export async function getPortfolio(address: string) {
  const res = await client.get<ApiResponse<unknown>>(`/portfolio/${address}`);
  return res.data.data;
}
