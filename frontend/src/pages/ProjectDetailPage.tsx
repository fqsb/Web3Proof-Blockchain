import { Button, Card, Col, Descriptions, Empty, List, Row, Space, Statistic, Steps, Tag, Typography, Upload, message } from "antd";
import { CloudUploadOutlined, FileDoneOutlined, LinkOutlined, SafetyCertificateOutlined, SendOutlined } from "@ant-design/icons";
import { useEffect, useMemo, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Certificate, confirmEvidence, EvidenceRecord, getWork, prepareEvidence, uploadWorkFile, Work, WorkFile } from "../api/works";
import { submitEvidenceOnChain } from "../utils/projectChain";

function statusColor(status?: string) {
  if (status === "confirmed" || status === "active") return "green";
  if (status === "pending_chain" || status === "pending") return "gold";
  if (status === "revoked" || status === "failed") return "red";
  return "blue";
}

function formatSize(size?: number) {
  if (!size) return "0 KB";
  if (size < 1024 * 1024) return `${(size / 1024).toFixed(1)} KB`;
  return `${(size / 1024 / 1024).toFixed(2)} MB`;
}

export default function ProjectDetailPage() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const [work, setWork] = useState<Work | null>(null);
  const [files, setFiles] = useState<WorkFile[]>([]);
  const [evidences, setEvidences] = useState<EvidenceRecord[]>([]);
  const [certificates, setCertificates] = useState<Certificate[]>([]);
  const [busy, setBusy] = useState(false);

  async function load() {
    try {
      const data = await getWork(id);
      setWork(data.work);
      setFiles(data.files);
      setEvidences(data.evidences);
      setCertificates(data.certificates);
    } catch (err) {
      message.error((err as Error).message || "加载作品详情失败，请稍后重试。");
    }
  }

  useEffect(() => { load(); }, [id]);

  const confirmed = useMemo(() => evidences.filter((item) => item.status === "confirmed"), [evidences]);
  const stepIndex = files.length === 0 ? 0 : confirmed.length === 0 ? 1 : certificates.length === 0 ? 2 : 3;

  async function upload(option: any) {
    try {
      const file = option.file as File;
      await uploadWorkFile(id, file);
      option.onSuccess?.({}, new XMLHttpRequest());
      message.success("文件已上传，SHA-256 哈希已计算完成。");
      await load();
    } catch (err) {
      option.onError?.(err as Error);
      message.error((err as Error).message);
    }
  }

  async function submitChain() {
    setBusy(true);
    try {
      const prepared = await prepareEvidence(id);
      const chain = await submitEvidenceOnChain(prepared);
      await confirmEvidence(id, chain.txHash, chain.chainEvidenceId);
      message.success("Sepolia 链上存证已确认，电子证书已生成。");
      await load();
    } catch (err) {
      message.error((err as Error).message);
    } finally {
      setBusy(false);
    }
  }

  if (!work) return null;

  return (
    <div className="inner-page">
      <div className="page-kicker">Evidence Detail</div>
      <Space className="page-title" align="center">
        <div>
          <Typography.Title level={2}>{work.title}</Typography.Title>
          <Typography.Paragraph>{work.description || "暂无作品说明，可在后续版本中补充编辑能力。"}</Typography.Paragraph>
        </div>
        <Space wrap>
          <Button onClick={() => navigate("/creator/applications")}>认证中心</Button>
          <Button type="primary" icon={<SendOutlined />} loading={busy} disabled={!files.length} onClick={submitChain}>
            发起 Sepolia 存证
          </Button>
        </Space>
      </Space>

      <Card className="section-card">
        <Steps
          current={stepIndex}
          items={[
            { title: "创建作品", description: "登记作品信息" },
            { title: "上传文件", description: "生成文件哈希" },
            { title: "链上存证", description: "写入 Sepolia" },
            { title: "证书归档", description: "可公开核验" },
          ]}
        />
      </Card>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={16}>
          <Card title="作品信息" className="section-card">
            <Descriptions column={1}>
              <Descriptions.Item label="状态"><Tag color={statusColor(work.status)}>{work.status}</Tag></Descriptions.Item>
              <Descriptions.Item label="可见性"><Tag>{work.visibility}</Tag></Descriptions.Item>
              <Descriptions.Item label="外部链接">
                {work.external_url ? <a href={work.external_url} target="_blank" rel="noreferrer">{work.external_url}</a> : "未填写"}
              </Descriptions.Item>
              <Descriptions.Item label="创建时间">{work.created_at}</Descriptions.Item>
            </Descriptions>
          </Card>
        </Col>
        <Col xs={24} lg={8}>
          <Card className="metric-stack">
            <Statistic title="文件数" value={files.length} />
            <Statistic title="链上存证" value={confirmed.length} />
            <Statistic title="电子证书" value={certificates.length} />
          </Card>
        </Col>
      </Row>

      <Row gutter={[16, 16]} className="detail-card-row">
        <Col xs={24} lg={12}>
          <Card title="文件与哈希" className="section-card detail-card" extra={<Upload customRequest={upload} showUploadList={false}><Button icon={<CloudUploadOutlined />}>上传文件</Button></Upload>}>
            {files.length ? (
              <List
                dataSource={files}
                renderItem={(file) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={<FileDoneOutlined />}
                      title={file.file_name}
                      description={
                        <Space direction="vertical" size={4}>
                          <span>{formatSize(file.file_size)}</span>
                          <code>{file.sha256_hash}</code>
                        </Space>
                      }
                    />
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="还没有上传文件，上传后系统会自动计算 SHA-256 哈希。" />
            )}
          </Card>
        </Col>

        <Col xs={24} lg={12}>
          <Card title="链上存证记录" className="section-card detail-card">
            {evidences.length ? (
              <List
                dataSource={evidences}
                renderItem={(item) => (
                  <List.Item className="evidence-list-item">
                    <List.Item.Meta
                      avatar={<SafetyCertificateOutlined />}
                      title={<Space wrap><span>{item.evidence_no}</span><Tag color={statusColor(item.status)}>{item.status}</Tag></Space>}
                      description={
                        <Space direction="vertical" size={4}>
                          <span>链上 ID：{item.chain_evidence_id || "待确认"}</span>
                          <code>{item.file_hash}</code>
                          {item.tx_hash && <code>{item.tx_hash}</code>}
                        </Space>
                      }
                    />
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无存证记录。上传文件后可提交到 Sepolia。" />
            )}
          </Card>
        </Col>
      </Row>

      <Card title="电子证书" className="section-card">
        {certificates.length ? (
          <List
            grid={{ gutter: 16, xs: 1, md: 2 }}
            dataSource={certificates}
            renderItem={(item) => (
              <List.Item>
                <Card className="certificate-card">
                  <Space direction="vertical" size={8}>
                    <Tag color="green">已签发</Tag>
                    <Typography.Title level={4}>{item.certificate_no}</Typography.Title>
                    <span>关联存证 ID：{item.evidence_id}</span>
                    <a href={item.verify_url} target="_blank" rel="noreferrer"><LinkOutlined /> 核验证书</a>
                  </Space>
                </Card>
              </List.Item>
            )}
          />
        ) : (
          <Empty description="链上存证确认后会生成可核验电子证书。" />
        )}
      </Card>
    </div>
  );
}
