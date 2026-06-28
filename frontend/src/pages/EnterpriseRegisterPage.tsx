import { Button, Card, Col, Form, Input, Row, Typography, message } from "antd";
import { AuditOutlined, BankOutlined } from "@ant-design/icons";
import { requestVerifierRole } from "../api/auth";
import { useAuth } from "../hooks/useAuth";

export default function EnterpriseRegisterPage() {
  const { refreshUser } = useAuth();

  async function onFinish(values: any) {
    try {
      await requestVerifierRole(values);
      await refreshUser();
      message.success("核验方身份已开通，可在顶部切换角色。");
    } catch (err) {
      message.error((err as Error).message || "开通核验方失败，请稍后重试。");
    }
  }

  return (
    <div className="inner-page">
      <div className="page-kicker">Verifier Onboarding</div>
      <Typography.Title level={2}>核验方入驻</Typography.Title>
      <Typography.Paragraph>为学校、实验室、企业或第三方机构开通核验方工作台，用于保存核验报告和查询可信档案。</Typography.Paragraph>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={9}>
          <Card className="profile-card">
            <BankOutlined className="profile-card__avatar" />
            <Typography.Title level={4}>机构核验身份</Typography.Title>
            <p>入驻后会为当前钱包增加 verifier 角色，可在导航栏切换到核验方工作台。</p>
            <p><AuditOutlined /> 支持文件哈希核验、证书核验、钱包可信档案查询和报告留存。</p>
          </Card>
        </Col>
        <Col xs={24} lg={15}>
          <Card title="机构信息" className="form-card">
            <Form layout="vertical" onFinish={onFinish}>
              <Form.Item label="机构名称" name="org_name" rules={[{ required: true, message: "请输入机构名称" }]}><Input /></Form.Item>
              <Form.Item label="行业/领域" name="industry"><Input placeholder="教育、软件开发、版权服务、招聘等" /></Form.Item>
              <Form.Item label="联系邮箱" name="contact_email"><Input /></Form.Item>
              <Form.Item label="官网" name="website"><Input placeholder="https://..." /></Form.Item>
              <Button type="primary" htmlType="submit">开通核验方</Button>
            </Form>
          </Card>
        </Col>
      </Row>
    </div>
  );
}
