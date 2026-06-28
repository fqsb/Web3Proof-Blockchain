import { Button, Card, Col, Empty, Form, Input, List, Row, Space, Tag, Typography, Upload, message } from "antd";
import { CheckCircleOutlined, CloseCircleOutlined, FileSearchOutlined, InboxOutlined, SafetyCertificateOutlined, WalletOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { listReports, verifyCertificate, verifyEvidence, verifyFile, verifyWallet } from "../api/verify";

type VerifyResult = {
  passed?: boolean;
  query_type?: string;
  query_value?: string;
  reason?: string;
  file_hash?: string;
  verified_at?: string;
  evidence?: any;
  work?: any;
  certificate?: any;
  credential?: any;
  user?: any;
  works?: any[];
  evidences?: any[];
  credentials?: any[];
};

type ReportRow = {
  id: number;
  query_type: string;
  query_value: string;
  passed: boolean;
  report_json?: string;
  created_at?: string;
};

function parseReport(row: ReportRow): VerifyResult {
  try {
    return row.report_json ? JSON.parse(row.report_json) : {};
  } catch {
    return {};
  }
}

function formatDateTime(value?: string) {
  if (!value) return "-";
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) return value;
  return new Intl.DateTimeFormat("zh-CN", {
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  }).format(date);
}

function ResultCard({ result }: { result: VerifyResult }) {
  const passed = Boolean(result.passed);
  const Icon = passed ? CheckCircleOutlined : CloseCircleOutlined;
  return (
    <Card className={`result-card ${passed ? "result-card--pass" : "result-card--fail"}`}>
      <Space align="start" size={16}>
        <Icon className="result-card__icon" />
        <div>
          <Tag color={passed ? "green" : "red"}>{passed ? "核验通过" : "未通过"}</Tag>
          <Typography.Title level={3}>{passed ? "材料与可信记录匹配" : "未找到可确认的可信记录"}</Typography.Title>
          <Typography.Paragraph>
            {passed ? "系统已匹配数据库中的确认存证记录，可继续查看存证、证书与钱包身份信息。" : result.reason || "请确认编号、钱包地址或上传文件是否正确。"}
          </Typography.Paragraph>
        </div>
      </Space>

      <Row gutter={[12, 12]} className="result-facts">
        <Col xs={24} md={8}><div><span>查询类型</span><strong>{result.query_type || "-"}</strong></div></Col>
        <Col xs={24} md={8}><div><span>查询值</span><strong>{result.query_value || result.file_hash || "-"}</strong></div></Col>
        <Col xs={24} md={8}><div><span>核验时间</span><strong>{formatDateTime(result.verified_at)}</strong></div></Col>
      </Row>

      {passed && (
        <div className="verify-detail-grid">
          {result.work && <div><strong>作品信息</strong><p>{result.work.title || "未命名作品"}</p><span>{result.work.description || "暂无说明"}</span></div>}
          {result.evidence && <div><strong>存证记录</strong><p>{result.evidence.evidence_no}</p><code>{result.evidence.file_hash}</code></div>}
          {result.certificate?.certificate_no && <div><strong>证书</strong><p>{result.certificate.certificate_no}</p><span>证书可公开核验</span></div>}
          {result.user && <div><strong>钱包身份</strong><p>{result.user.wallet_address}</p><span>DID：{result.user.did || "未登记"}</span></div>}
        </div>
      )}
    </Card>
  );
}

export default function VerifyPage({ reportsMode = false }: { reportsMode?: boolean }) {
  const [result, setResult] = useState<VerifyResult | null>(null);
  const [reports, setReports] = useState<ReportRow[]>([]);

  useEffect(() => {
    if (!reportsMode) return;
    listReports()
      .then((items) => setReports(items as ReportRow[]))
      .catch((err) => message.error((err as Error).message || "加载核验报告失败，请稍后重试。"));
  }, [reportsMode]);

  async function runQuery(task: () => Promise<unknown>) {
    try {
      setResult(await task() as VerifyResult);
    } catch (err) {
      message.error((err as Error).message || "查询失败，请检查输入后重试。");
    }
  }

  async function upload(option: any) {
    try {
      const report = await verifyFile(option.file as File);
      setResult(report as VerifyResult);
      option.onSuccess?.({}, new XMLHttpRequest());
      message.success("文件哈希核验完成。");
    } catch (err) {
      option.onError?.(err as Error);
      message.error((err as Error).message);
    }
  }

  if (reportsMode) {
    return (
      <div className="inner-page">
        <div className="page-kicker">Verifier Workspace</div>
        <Typography.Title level={2}>核验报告</Typography.Title>
        <Typography.Paragraph>这里记录核验方近期发起的文件、证书、存证编号和钱包核验结果。</Typography.Paragraph>
        {reports.length ? (
          <List
            dataSource={reports}
            renderItem={(item) => {
              const detail = parseReport(item);
              return (
                <List.Item>
                  <Card className="report-card">
                    <Space direction="vertical" size={8}>
                      <Space wrap>
                        <Tag color={item.passed ? "green" : "red"}>{item.passed ? "通过" : "未通过"}</Tag>
                        <span>{item.query_type}</span>
                        <span>{formatDateTime(item.created_at)}</span>
                      </Space>
                      <Typography.Title level={4}>{item.query_value}</Typography.Title>
                      <p>{item.passed ? "已匹配可信存证记录" : detail.reason || "未找到匹配记录"}</p>
                    </Space>
                  </Card>
                </List.Item>
              );
            }}
          />
        ) : (
          <Empty description="暂无核验报告" />
        )}
      </div>
    );
  }

  return (
    <div className="inner-page">
      <div className="page-kicker">Public Verification</div>
      <Typography.Title level={2}>公开核验入口</Typography.Title>
      <Typography.Paragraph>
        上传原文件或输入证书编号、存证编号、钱包地址，系统会重新计算哈希并匹配链上存证记录。
      </Typography.Paragraph>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={10}>
          <Card title="上传文件核验" className="verify-console">
            <Upload.Dragger customRequest={upload} showUploadList={false} multiple={false}>
              <p className="ant-upload-drag-icon"><InboxOutlined /></p>
              <p className="ant-upload-text">拖拽或点击选择原文件</p>
              <p className="ant-upload-hint">文件不会上链，只计算 SHA-256 后进行匹配。</p>
            </Upload.Dragger>
          </Card>
        </Col>
        <Col xs={24} lg={14}>
          <Card title="编号与钱包核验" className="verify-console">
            <Form layout="vertical" onFinish={(v) => runQuery(() => verifyEvidence(v.evidence_no))}>
              <Form.Item name="evidence_no" label="存证编号" rules={[{ required: true, message: "请输入存证编号" }]}>
                <Input prefix={<FileSearchOutlined />} placeholder="例如 W3P-EV-..." />
              </Form.Item>
              <Button htmlType="submit">查询存证</Button>
            </Form>
            <Form layout="vertical" className="inline-form" onFinish={(v) => runQuery(() => verifyCertificate(v.certificate_no))}>
              <Form.Item name="certificate_no" label="证书编号" rules={[{ required: true, message: "请输入证书编号" }]}>
                <Input prefix={<SafetyCertificateOutlined />} placeholder="例如 W3P-CERT-..." />
              </Form.Item>
              <Button htmlType="submit">查询证书</Button>
            </Form>
            <Form layout="vertical" className="inline-form" onFinish={(v) => runQuery(() => verifyWallet(v.wallet))}>
              <Form.Item name="wallet" label="钱包地址" rules={[{ required: true, message: "请输入钱包地址" }]}>
                <Input prefix={<WalletOutlined />} placeholder="0x..." />
              </Form.Item>
              <Button htmlType="submit">查询钱包档案</Button>
            </Form>
          </Card>
        </Col>
      </Row>

      {result && <ResultCard result={result} />}
    </div>
  );
}
