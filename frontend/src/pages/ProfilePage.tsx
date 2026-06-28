import { Button, Card, Col, Form, Input, Row, Space, Tag, Typography, message } from "antd";
import { IdcardOutlined, SaveOutlined, UserOutlined } from "@ant-design/icons";
import { updateProfile } from "../api/auth";
import { useAuth } from "../hooks/useAuth";

export default function ProfilePage() {
  const { user, refreshUser } = useAuth();

  async function onFinish(values: any) {
    try {
      await updateProfile(values);
      await refreshUser();
      message.success("个人资料已更新。");
    } catch (err) {
      message.error((err as Error).message || "保存个人资料失败，请稍后重试。");
    }
  }

  return (
    <div className="inner-page">
      <div className="page-kicker">Account Profile</div>
      <Typography.Title level={2}>个人资料</Typography.Title>
      <Typography.Paragraph>完善昵称、邮箱、头像和简介，让公开可信档案更完整。</Typography.Paragraph>

      <Row gutter={[16, 16]}>
        <Col xs={24} lg={8}>
          <Card className="profile-card">
            <UserOutlined className="profile-card__avatar" />
            <Typography.Title level={4}>{user?.nickname || "未设置昵称"}</Typography.Title>
            <Space direction="vertical">
              <Tag color="blue"><IdcardOutlined /> {user?.active_role}</Tag>
              <code>{user?.wallet_address}</code>
              <span>DID：{user?.did || "未登记"}</span>
            </Space>
          </Card>
        </Col>
        <Col xs={24} lg={16}>
          <Card title="资料编辑" className="form-card">
            <Form layout="vertical" initialValues={user || {}} onFinish={onFinish}>
              <Form.Item label="昵称" name="nickname"><Input placeholder="展示在公开档案中的名称" /></Form.Item>
              <Form.Item label="邮箱" name="email"><Input placeholder="用于联系和认证沟通" /></Form.Item>
              <Form.Item label="头像 URL" name="avatar_url"><Input placeholder="https://..." /></Form.Item>
              <Form.Item label="简介" name="bio"><Input.TextArea rows={5} placeholder="介绍你的创作方向、课程成果、项目经历或认证背景。" /></Form.Item>
              <Button type="primary" icon={<SaveOutlined />} htmlType="submit">保存资料</Button>
            </Form>
          </Card>
        </Col>
      </Row>
    </div>
  );
}
