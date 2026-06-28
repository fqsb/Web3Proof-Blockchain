import { Button, Card, Empty, Form, Input, List, Select, Space, Tabs, Tag, Typography, message } from "antd";
import { AuditOutlined, SafetyCertificateOutlined, SendOutlined } from "@ant-design/icons";
import { useEffect, useMemo, useState } from "react";
import { applyCertification, listMyApplications } from "../api/applications";
import { EvidenceRecord, listMyCertificates, listMyEvidence, Certificate } from "../api/works";

type ApplicationItem = {
  id: number;
  work_id: number;
  evidence_id: number;
  materials_desc?: string;
  status: string;
  review_note?: string;
  created_at?: string;
  reviewed_at?: string;
  work?: { title?: string; description?: string };
};

function statusColor(status?: string) {
  if (status === "approved" || status === "active") return "green";
  if (status === "pending") return "gold";
  if (status === "rejected") return "red";
  return "blue";
}

export default function CertificationsPage() {
  const [applications, setApplications] = useState<ApplicationItem[]>([]);
  const [evidences, setEvidences] = useState<EvidenceRecord[]>([]);
  const [certificates, setCertificates] = useState<Certificate[]>([]);

  async function load() {
    try {
      const [apps, evs, certs] = await Promise.all([listMyApplications(), listMyEvidence(), listMyCertificates()]);
      setApplications(apps as ApplicationItem[]);
      setEvidences(evs);
      setCertificates(certs);
    } catch (err) {
      message.error((err as Error).message || "加载认证数据失败，请稍后重试。");
    }
  }

  useEffect(() => { load(); }, []);

  const confirmedEvidence = useMemo(() => evidences.filter((item) => item.status === "confirmed"), [evidences]);

  async function onFinish(values: { evidence_id: number; materials_desc?: string }) {
    const evidence = evidences.find((item) => item.id === values.evidence_id);
    if (!evidence) {
      message.error("请选择有效的存证记录。");
      return;
    }
    try {
      await applyCertification({ work_id: evidence.work_id, evidence_id: evidence.id, materials_desc: values.materials_desc });
      message.success("认证申请已提交，等待审核员处理。");
      await load();
    } catch (err) {
      message.error((err as Error).message || "提交认证申请失败，请稍后重试。");
    }
  }

  return (
    <div className="inner-page">
      <div className="page-kicker">Certification Center</div>
      <div className="page-title">
        <div>
          <Typography.Title level={2}>认证与证书</Typography.Title>
          <Typography.Paragraph>
            基于已确认的链上存证提交认证申请，审核通过后可获得不可转让 SBT 凭证，并展示在公开可信档案中。
          </Typography.Paragraph>
        </div>
      </div>

      <Tabs
        className="web3-tabs"
        items={[
          {
            key: "apply",
            label: "提交认证",
            children: (
              <Card title="基于链上存证申请认证" className="form-card">
                <Form layout="vertical" onFinish={onFinish}>
                  <Form.Item label="已确认存证" name="evidence_id" rules={[{ required: true, message: "请选择存证记录" }]}>
                    <Select
                      placeholder="选择一条 Sepolia 已确认的存证"
                      options={confirmedEvidence.map((item) => ({
                        value: item.id,
                        label: `${item.evidence_no} · 作品 #${item.work_id}`,
                      }))}
                    />
                  </Form.Item>
                  <Form.Item label="补充说明" name="materials_desc">
                    <Input.TextArea rows={5} placeholder="说明作品来源、获奖信息、课程背景、审核依据或其他证明材料。" />
                  </Form.Item>
                  <Button type="primary" htmlType="submit" icon={<SendOutlined />} disabled={!confirmedEvidence.length}>
                    提交认证申请
                  </Button>
                </Form>
                {!confirmedEvidence.length && <div className="form-hint">当前没有已确认的链上存证，请先到作品详情页完成 Sepolia 存证。</div>}
              </Card>
            ),
          },
          {
            key: "applications",
            label: "我的申请",
            children: applications.length ? (
              <List
                grid={{ gutter: 16, xs: 1, md: 2 }}
                dataSource={applications}
                renderItem={(item) => (
                  <List.Item>
                    <Card className="application-card">
                      <Space direction="vertical" size={8}>
                        <Space wrap>
                          <Tag color={statusColor(item.status)}>{item.status}</Tag>
                          <span>申请 #{item.id}</span>
                        </Space>
                        <Typography.Title level={4}>{item.work?.title || `作品 #${item.work_id}`}</Typography.Title>
                        <span>关联存证：#{item.evidence_id}</span>
                        <p>{item.materials_desc || "未填写补充说明"}</p>
                        {item.review_note && <div className="review-note"><AuditOutlined /> {item.review_note}</div>}
                      </Space>
                    </Card>
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无认证申请" />
            ),
          },
          {
            key: "certificates",
            label: "我的证书",
            children: certificates.length ? (
              <List
                grid={{ gutter: 16, xs: 1, md: 2, lg: 3 }}
                dataSource={certificates}
                renderItem={(item) => (
                  <List.Item>
                    <Card className="certificate-card">
                      <SafetyCertificateOutlined className="certificate-card__icon" />
                      <Typography.Title level={4}>{item.certificate_no}</Typography.Title>
                      <p>关联存证 ID：{item.evidence_id}</p>
                      <a href={item.verify_url} target="_blank" rel="noreferrer">打开核验链接</a>
                    </Card>
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无证书，完成链上存证后会自动生成。" />
            ),
          },
        ]}
      />
    </div>
  );
}
