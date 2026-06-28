import { Card, Progress, Tag } from "antd";

interface Reputation {
  total_score: number;
  grade: string;
  project_score: number;
  cert_score: number;
  activity_score: number;
}

const gradeColor: Record<string, string> = { A: "green", B: "blue", C: "orange", D: "red" };

export default function ReputationBadge({ reputation }: { reputation?: Reputation }) {
  if (!reputation) return null;
  return (
    <Card title="档案可信评分" size="small">
      <div style={{ textAlign: "center", marginBottom: 16 }}>
        <div
          style={{
            fontSize: 56,
            fontWeight: 700,
            lineHeight: 1,
            margin: 0,
            background: "linear-gradient(135deg, #fff 0%, #64d2ff 100%)",
            WebkitBackgroundClip: "text",
            WebkitTextFillColor: "transparent",
          }}
        >
          {reputation.total_score}
        </div>
        <Tag color={gradeColor[reputation.grade] || "default"} style={{ fontSize: 16 }}>
          评级 {reputation.grade}
        </Tag>
      </div>
      <div>证明材料: {reputation.project_score}/500</div>
      <Progress percent={(reputation.project_score / 500) * 100} showInfo={false} />
      <div>学校背书: {reputation.cert_score}/300</div>
      <Progress percent={(reputation.cert_score / 300) * 100} showInfo={false} />
      <div>资料完整度: {reputation.activity_score}/200</div>
      <Progress percent={(reputation.activity_score / 200) * 100} showInfo={false} />
    </Card>
  );
}
