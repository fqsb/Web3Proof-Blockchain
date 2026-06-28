import { Card, Col, Empty, List, Progress, Row, Space, Statistic, Tag, Typography } from "antd";
import { IdcardOutlined, SafetyCertificateOutlined, TrophyOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getPortfolio } from "../api/auth";

type PortfolioData = {
  user?: { wallet_address?: string; did?: string; nickname?: string; bio?: string };
  reputation?: { total_score?: number; grade?: string; project_score?: number; cert_score?: number; activity_score?: number };
  works?: any[];
  evidences?: any[];
  credentials?: any[];
};

export default function ResumePage() {
  const { address = "" } = useParams();
  const [data, setData] = useState<PortfolioData | null>(null);

  useEffect(() => {
    getPortfolio(address).then((payload) => setData(payload as PortfolioData));
  }, [address]);

  if (!data) return null;
  const score = Number(data.reputation?.total_score || 0);

  return (
    <div className="inner-page portfolio-page">
      <section className="portfolio-hero">
        <div>
          <div className="page-kicker">Trusted Portfolio</div>
          <Typography.Title level={2}>{data.user?.nickname || "公开可信档案"}</Typography.Title>
          <Typography.Paragraph>{data.user?.bio || "该用户公开展示的作品、链上存证记录和认证凭证。"}</Typography.Paragraph>
          <Space wrap>
            <Tag color="blue"><IdcardOutlined /> {data.user?.did || "DID 未登记"}</Tag>
            <Tag color="green"><SafetyCertificateOutlined /> {data.evidences?.length || 0} 条存证</Tag>
            <Tag color="gold"><TrophyOutlined /> {data.credentials?.length || 0} 个 SBT</Tag>
          </Space>
        </div>
        <Card className="score-card">
          <span>可信评分</span>
          <strong>{score}</strong>
          <Progress percent={Math.min(100, Math.round(score / 10))} showInfo={false} />
          <Tag color="gold">等级 {data.reputation?.grade || "D"}</Tag>
        </Card>
      </section>

      <Card className="section-card">
        <Row gutter={[16, 16]}>
          <Col xs={24} md={8}><Statistic title="公开作品" value={data.works?.length || 0} /></Col>
          <Col xs={24} md={8}><Statistic title="链上存证" value={data.evidences?.length || 0} /></Col>
          <Col xs={24} md={8}><Statistic title="认证凭证" value={data.credentials?.length || 0} /></Col>
        </Row>
        <div className="identity-box">
          <span>钱包地址</span>
          <strong>{data.user?.wallet_address}</strong>
        </div>
      </Card>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={12}>
          <Card title="公开作品" className="section-card">
            {data.works?.length ? (
              <List
                dataSource={data.works}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta title={item.title} description={item.description || "暂无说明"} />
                    <Tag>{item.status}</Tag>
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无公开作品" />
            )}
          </Card>
        </Col>
        <Col xs={24} lg={12}>
          <Card title="链上存证" className="section-card">
            {data.evidences?.length ? (
              <List
                dataSource={data.evidences}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta title={item.evidence_no} description={<code>{item.file_hash}</code>} />
                    <Tag color="green">{item.status}</Tag>
                  </List.Item>
                )}
              />
            ) : (
              <Empty description="暂无公开存证" />
            )}
          </Card>
        </Col>
      </Row>

      <Card title="认证 SBT 凭证" className="section-card">
        {data.credentials?.length ? (
          <List
            grid={{ gutter: 16, xs: 1, md: 2, lg: 3 }}
            dataSource={data.credentials}
            renderItem={(item) => (
              <List.Item>
                <Card className="credential-card">
                  <Tag color="green">{item.status}</Tag>
                  <Typography.Title level={4}>Token #{item.token_id}</Typography.Title>
                  <p>关联作品：#{item.work_id}</p>
                  <code>{item.tx_hash}</code>
                </Card>
              </List.Item>
            )}
          />
        ) : (
          <Empty description="暂无 SBT 认证凭证" />
        )}
      </Card>
    </div>
  );
}
