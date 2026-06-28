import { Card, List, Tag } from "antd";
import { BLOCK_EXPLORER } from "../config";

interface SBT {
  id: number;
  token_id: number;
  skill?: { name: string };
  tx_hash: string;
}

export default function SBTCertificateGrid({ sbts }: { sbts: SBT[] }) {
  return (
    <Card title="学校背书凭证" size="small">
      <List
        grid={{ gutter: 8, column: 2 }}
        dataSource={sbts}
        locale={{ emptyText: "暂无学校背书凭证" }}
        renderItem={(s) => (
          <List.Item>
            <Card size="small">
              <Tag color="blue">{s.skill?.name || "背书类别"}</Tag>
              <div>Token #{s.token_id}</div>
              <a href={`${BLOCK_EXPLORER}/tx/${s.tx_hash}`} target="_blank" rel="noreferrer">
                查看可信记录
              </a>
            </Card>
          </List.Item>
        )}
      />
    </Card>
  );
}
