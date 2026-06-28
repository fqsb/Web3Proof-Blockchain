import client, { ApiResponse } from "./client";

export async function getDashboardSummary() {
  const res = await client.get<ApiResponse<{
    counts: Record<string, number>;
    reputation?: { total_score: number; grade: string };
    recent_works: unknown[];
    recent_evidences: unknown[];
  }>>("/dashboard/summary");
  return res.data.data;
}
