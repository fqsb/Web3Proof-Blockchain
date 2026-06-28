import { Button, Card, Col, Row, Space, Statistic, Tag, message } from "antd";
import { CheckCircleOutlined, CloudServerOutlined, FileProtectOutlined, SafetyCertificateOutlined, ScanOutlined, TeamOutlined } from "@ant-design/icons";
import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../hooks/useAuth";

const capabilityCards = [
  { icon: <FileProtectOutlined />, title: "作品哈希存证", text: "上传文件后计算 SHA-256，链上只保存摘要、作者地址和存证编号。" },
  { icon: <SafetyCertificateOutlined />, title: "SBT 认证凭证", text: "审核员通过材料后发放不可转让凭证，适合证书、奖项和成果认证。" },
  { icon: <ScanOutlined />, title: "多入口核验", text: "支持文件、存证编号、证书编号、钱包地址四种核验方式。" },
  { icon: <TeamOutlined />, title: "多角色工作台", text: "创作者、核验方、审核员、管理员可以在一个账号内切换使用。" },
];

const flows = ["上传文件", "生成哈希", "Sepolia 存证", "生成证书", "公开核验"];

export default function HomePage() {
  const { user, signIn } = useAuth();
  const navigate = useNavigate();
  const [signingIn, setSigningIn] = useState(false);

  async function handleSignIn() {
    setSigningIn(true);
    try {
      await signIn();
    } catch (err) {
      message.error((err as Error).message || "连接钱包失败，请重试。");
    } finally {
      setSigningIn(false);
    }
  }

  return (
    <div className="home">
      <section className="home-hero">
        <div className="home-hero__copy">
          <Space wrap className="home-hero__tags">
            <Tag color="cyan">Sepolia 测试网</Tag>
            <Tag color="blue">链上摘要</Tag>
            <Tag color="green">链下文件</Tag>
          </Space>
          <p className="eyebrow">Blockchain Evidence Platform</p>
          <h1 className="home-hero__title">
            <span>链上可信存证</span>
            <span>认证核验平台</span>
          </h1>
          <p className="lead">
            Web3Proof 面向课程成果、竞赛材料、原创作品、证书证明和代码项目，提供从文件哈希存证、审核认证、SBT 凭证到第三方核验报告的一站式流程。
          </p>
          <div className="actions">
            <Button type="primary" size="large" loading={signingIn} onClick={user ? () => navigate("/dashboard") : handleSignIn}>
              {user ? "进入工作台" : "连接钱包开始存证"}
            </Button>
            <Button size="large" onClick={() => navigate("/verify")}>立即核验材料</Button>
            <Button size="large" loading={signingIn} onClick={user ? () => navigate("/creator/works") : handleSignIn}>查看作品存证</Button>
          </div>
        </div>

        <div className="evidence-console">
          <div className="console-header">
            <span className="dot dot-green" />
            <span className="dot dot-yellow" />
            <span className="dot dot-red" />
            <strong>Evidence Live Preview</strong>
          </div>
          <div className="hash-card">
            <span>FILE HASH</span>
            <code>0x9f42...a7c8e12</code>
          </div>
          <div className="console-grid">
            <div><span>存证编号</span><strong>EV-20260626-000128</strong></div>
            <div><span>作者地址</span><strong>0x2d63...62f0</strong></div>
            <div><span>认证状态</span><strong>待审核</strong></div>
            <div><span>网络</span><strong>Sepolia</strong></div>
          </div>
          <div className="chain-line">
            {flows.map((item, index) => (
              <div key={item} className="chain-step">
                <span>{index + 1}</span>
                <strong>{item}</strong>
              </div>
            ))}
          </div>
        </div>
      </section>

      <Row gutter={[16, 16]} className="home-stats">
        <Col xs={12} md={6}><Card><Statistic title="存证对象" value="作品/成果" /></Card></Col>
        <Col xs={12} md={6}><Card><Statistic title="链上凭证" value="SBT" /></Card></Col>
        <Col xs={12} md={6}><Card><Statistic title="核验方式" value={4} suffix="种" /></Card></Col>
        <Col xs={12} md={6}><Card><Statistic title="部署环境" value="Linux" /></Card></Col>
      </Row>

      <section className="home-section">
        <div className="section-heading">
          <p className="eyebrow">Core Workflow</p>
          <h2>一条完整的可信存证链路</h2>
        </div>
        <div className="workflow-board">
          {flows.map((item, index) => (
            <Card key={item} className="workflow-card">
              <span>{String(index + 1).padStart(2, "0")}</span>
              <strong>{item}</strong>
              <p>{["文件保存在服务器或对象存储。", "系统计算不可逆文件摘要。", "调用智能合约写入摘要。", "生成带核验入口的证书。", "第三方可独立比对结果。"][index]}</p>
            </Card>
          ))}
        </div>
      </section>

      <section className="home-section">
        <div className="section-heading">
          <p className="eyebrow">Modules</p>
          <h2>不是只有上链，平台功能也完整</h2>
        </div>
        <Row gutter={[16, 16]}>
          {capabilityCards.map((item) => (
            <Col xs={24} md={12} lg={6} key={item.title}>
              <Card className="capability-card">
                <div className="capability-icon">{item.icon}</div>
                <h3>{item.title}</h3>
                <p>{item.text}</p>
              </Card>
            </Col>
          ))}
        </Row>
      </section>

      <section className="home-section role-section">
        <Card>
          <CloudServerOutlined className="role-section__icon" />
          <h2>适合部署到国内服务器的架构</h2>
          <p>
            前端、后端、MySQL、Redis 和文件存储部署在 Linux 服务器；Sepolia 作为测试网络承载链上写入验证。
            生产场景可切换为私有链或联盟链，系统不发行代币、不提供交易市场。
          </p>
        </Card>
        <Card>
          <CheckCircleOutlined className="role-section__icon" />
          <h2>可信流程清晰</h2>
          <p>
            连接钱包、上传作品、生成哈希、发起链上存证、下载证书、提交认证、审核发放 SBT、上传原文件核验，形成完整可追溯链路。
          </p>
        </Card>
      </section>
    </div>
  );
}
