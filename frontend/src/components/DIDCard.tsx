import { Card, Descriptions, Tag, Typography } from "antd";
import { User } from "../api/auth";

export default function DIDCard({ user }: { user: User }) {
  return (
    <Card title="候选人信息" size="small">
      <Descriptions column={1} size="small">
        <Descriptions.Item label="公开身份">{user.did || "-"}</Descriptions.Item>
        <Descriptions.Item label="钱包地址">{user.wallet_address}</Descriptions.Item>
        <Descriptions.Item label="可信记录">
          <Tag color={user.is_did_registered ? "green" : "default"}>
            {user.is_did_registered ? "已生成" : "未生成"}
          </Tag>
        </Descriptions.Item>
        <Descriptions.Item label="作品主页">{user.github || "-"}</Descriptions.Item>
        <Descriptions.Item label="姓名/昵称">{user.nickname || "-"}</Descriptions.Item>
      </Descriptions>
      {user.bio && <Typography.Paragraph type="secondary">{user.bio}</Typography.Paragraph>}
    </Card>
  );
}
