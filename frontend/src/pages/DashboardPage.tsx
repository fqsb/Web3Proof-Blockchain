import { Button, Card, Col, List, Progress, Row, Space, Statistic, Tag, Typography, message } from "antd";
import { BankOutlined, CloudUploadOutlined, FileSearchOutlined, IdcardOutlined, SafetyCertificateOutlined, TeamOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getDashboardSummary } from "../api/dashboard";
import { useAuth } from "../hooks/useAuth";

const roleName: Record<string, string> = {
  creator: "创作者",
  verifier: "核验方",
  auditor: "审核员",
  admin: "管理员",
};

export default function DashboardPage() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [summary, setSummary] = useState<any>(null);

  useEffect(() => {
    getDashboardSummary()
      .then(setSummary)
      .catch((err) => message.error((err as Error).message || "加载工作台数据失败，请稍后重试。"));
  }, []);

  const counts = summary?.counts || {};
  const score = Number(summary?.reputation?.total_score || 0);
  const roles = user?.roles || [];

  return (
    <div className="inner-page">
      <section className="dashboard-hero">
        <div>
          <div className="page-kicker">Web3Proof Workspace</div>
          <Typography.Title level={2}>可信成果工作台</Typography.Title>
          <Typography.Paragraph>
            从作品上传、SHA-256 哈希计算、Sepolia 链上存证，到认证审核、SBT 凭证和公开核验，所有关键证据集中管理。
          </Typography.Paragraph>
          <Space wrap>
            <Button type="primary" icon={<CloudUploadOutlined />} onClick={() => navigate("/creator/works/create")}>
              上传新作品
            </Button>
            <Button icon={<FileSearchOutlined />} onClick={() => navigate("/verify")}>
              核验材料
            </Button>
            <Button icon={<IdcardOutlined />} onClick={() => user && navigate(`/portfolio/${user.wallet_address}`)}>
              公开档案
            </Button>
          </Space>
        </div>

        <Card className="score-card">
          <span>当前身份</span>
          <Tag color="blue">{roleName[user?.active_role || "creator"]}</Tag>
          <strong>{score}</strong>
          <Progress percent={Math.min(100, Math.round(score / 10))} showInfo={false} />
          <span>可信评分 · 等级 {summary?.reputation?.grade || "D"}</span>
        </Card>
      </section>

      <Row gutter={[16, 16]} className="section-card">
        <Col xs={12} md={4}><Card className="metric-card"><Statistic title="作品" value={counts.works || 0} /></Card></Col>
        <Col xs={12} md={4}><Card className="metric-card"><Statistic title="存证" value={counts.evidences || 0} /></Card></Col>
        <Col xs={12} md={4}><Card className="metric-card"><Statistic title="证书" value={counts.certificates || 0} /></Card></Col>
        <Col xs={12} md={4}><Card className="metric-card"><Statistic title="申请" value={counts.applications || 0} /></Card></Col>
        <Col xs={12} md={4}><Card className="metric-card"><Statistic title="SBT" value={counts.credentials || 0} /></Card></Col>
        <Col xs={12} md={4}><Card className="metric-card"><Statistic title="角色" value={user?.roles?.length || 1} /></Card></Col>
      </Row>

      <Card title="角色与权限" className="section-card role-guide-card">
        <div className="role-guide">
          <div>
            <strong>普通钱包默认是创作者</strong>
            <p>连接钱包后即可上传作品、提交链上存证、发起认证申请。</p>
          </div>
          <div>
            <strong>核验方可自助入驻</strong>
            <p>填写机构信息后，当前钱包会增加核验方身份，可保存核验报告。</p>
            {!roles.includes("verifier") && (
              <Button icon={<BankOutlined />} onClick={() => navigate("/enterprise/register")}>
                申请核验方
              </Button>
            )}
          </div>
          <div>
            <strong>审核员和管理员由管理员授权</strong>
            <p>为避免普通用户自提权限，需要管理员在用户角色页分配 auditor 或 admin。</p>
            {roles.includes("admin") && (
              <Button icon={<TeamOutlined />} onClick={() => navigate("/admin/users")}>
                管理用户角色
              </Button>
            )}
          </div>
        </div>
      </Card>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <Card title="近期作品" className="panel-card">
            <List
              dataSource={summary?.recent_works || []}
              locale={{ emptyText: "还没有作品，先上传一份成果材料吧。" }}
              renderItem={(item: any) => (
                <List.Item actions={[<Button type="link" onClick={() => navigate(`/creator/works/${item.id}`)}>查看</Button>]}>
                  <List.Item.Meta title={item.title} description={`状态：${item.status} · 可见性：${item.visibility}`} />
                </List.Item>
              )}
            />
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="近期存证" className="panel-card">
            <List
              dataSource={summary?.recent_evidences || []}
              locale={{ emptyText: "暂无链上存证记录。" }}
              renderItem={(item: any) => (
                <List.Item>
                  <List.Item.Meta
                    avatar={<SafetyCertificateOutlined />}
                    title={item.evidence_no}
                    description={<code>{item.file_hash}</code>}
                  />
                  <Tag color={item.status === "confirmed" ? "green" : "default"}>{item.status}</Tag>
                </List.Item>
              )}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
}
