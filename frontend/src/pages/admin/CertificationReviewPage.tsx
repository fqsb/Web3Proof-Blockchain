import { Button, Card, Empty, List, Space, Tag, Typography, message } from "antd";
import { CheckOutlined, CloseOutlined, SafetyCertificateOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { confirmCredentialMint, listAuditApplications, prepareCredentialMint, reviewApplication } from "../../api/applications";
import { mintCredentialOnChain } from "../../utils/sbtChain";

type AuditApplication = {
  id: number;
  user_id: number;
  work_id: number;
  evidence_id: number;
  materials_desc?: string;
  status: "pending" | "approved" | "rejected" | string;
  review_note?: string;
  created_at?: string;
  work?: { title?: string; description?: string };
  user?: { wallet_address?: string; nickname?: string };
};

function statusColor(status?: string) {
  if (status === "approved") return "green";
  if (status === "pending") return "gold";
  if (status === "rejected") return "red";
  return "blue";
}

export default function CertificationReviewPage() {
  const [items, setItems] = useState<AuditApplication[]>([]);
  const [busy, setBusy] = useState<number | null>(null);

  async function load() {
    try {
      setItems(await listAuditApplications() as AuditApplication[]);
    } catch (err) {
      message.error((err as Error).message || "加载认证审核列表失败，请稍后重试。");
    }
  }

  useEffect(() => { load(); }, []);

  async function review(id: number, status: "approved" | "rejected") {
    const note = status === "approved" ? "材料完整，准予认证。" : "材料不足，请补充证明后重新提交。";
    try {
      await reviewApplication(id, status, note);
      message.success("审核结果已保存。");
      await load();
    } catch (err) {
      message.error((err as Error).message || "保存审核结果失败，请稍后重试。");
    }
  }

  async function mint(id: number) {
    setBusy(id);
    try {
      const prepared = await prepareCredentialMint(id);
      const chain = await mintCredentialOnChain(prepared);
      await confirmCredentialMint(id, { tx_hash: chain.txHash, token_id: chain.tokenId, token_uri: prepared.token_uri });
      message.success("SBT 认证凭证已发放。");
      await load();
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setBusy(null);
    }
  }

  return (
    <div className="inner-page">
      <div className="page-kicker">Auditor Workspace</div>
      <Typography.Title level={2}>认证审核</Typography.Title>
      <Typography.Paragraph>审核创作者基于链上存证提交的认证申请，通过后可继续发放不可转让 SBT 凭证。</Typography.Paragraph>

      {items.length ? (
        <List
          grid={{ gutter: 16, xs: 1, lg: 2 }}
          dataSource={items}
          renderItem={(item) => (
            <List.Item>
              <Card className="review-card">
                <Space direction="vertical" size={10}>
                  <Space wrap>
                    <Tag color={statusColor(item.status)}>{item.status}</Tag>
                    <span>申请 #{item.id}</span>
                    <span>{item.created_at}</span>
                  </Space>
                  <Typography.Title level={4}>{item.work?.title || `作品 #${item.work_id}`}</Typography.Title>
                  <p>{item.materials_desc || "申请人未填写补充说明。"}</p>
                  <div className="identity-box">
                    <span>申请人钱包</span>
                    <strong>{item.user?.wallet_address || `用户 #${item.user_id}`}</strong>
                  </div>
                  <Space wrap>
                    <Tag>作品 #{item.work_id}</Tag>
                    <Tag>存证 #{item.evidence_id}</Tag>
                  </Space>
                  <Space wrap>
                    {item.status === "pending" && <Button type="primary" icon={<CheckOutlined />} onClick={() => review(item.id, "approved")}>通过</Button>}
                    {item.status === "pending" && <Button danger icon={<CloseOutlined />} onClick={() => review(item.id, "rejected")}>驳回</Button>}
                    {item.status === "approved" && <Button icon={<SafetyCertificateOutlined />} loading={busy === item.id} onClick={() => mint(item.id)}>发放 SBT</Button>}
                  </Space>
                </Space>
              </Card>
            </List.Item>
          )}
        />
      ) : (
        <Empty description="暂无待审核申请" />
      )}
    </div>
  );
}
