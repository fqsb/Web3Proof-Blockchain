import client, { ApiResponse } from "./client";
import { RoleCode } from "./auth";

export async function getAdminData(mode: string) {
  const res = await client.get<ApiResponse<unknown>>(`/admin/${mode}`);
  return res.data.data;
}

export async function updateUserRoles(id: number, roles: RoleCode[], active_role: RoleCode) {
  const res = await client.put<ApiResponse<unknown>>(`/admin/users/${id}/roles`, { roles, active_role });
  return res.data.data;
}

export async function syncChainEvents(lookback = 5000) {
  const res = await client.post<ApiResponse<{ scanned: number; inserted: number }>>(`/admin/chain-events/sync`, null, {
    params: { lookback },
  });
  return res.data.data;
}
