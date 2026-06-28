import { Card, List, Tag, Typography } from "antd";
import { BLOCK_EXPLORER } from "../config";

interface Project {
  id: number;
  name: string;
  github_url?: string;
  ipfs_cid?: string;
  tx_hash?: string;
  status: string;
}

export default function ProjectTimeline({ projects }: { projects: Project[] }) {
  return (
    <Card title="证明材料记录" size="small">
      <List
        dataSource={projects}
        locale={{ emptyText: "暂无可核验证明材料" }}
        renderItem={(p) => (
          <List.Item>
            <List.Item.Meta
              title={p.name}
              description={
                <>
                  {p.github_url && (
                    <Typography.Link href={p.github_url} target="_blank">
                      证明链接
                    </Typography.Link>
                  )}
                  {p.ipfs_cid && (
                    <Typography.Text type="secondary" style={{ marginLeft: 8 }}>
                      CID: {p.ipfs_cid.slice(0, 16)}...
                    </Typography.Text>
                  )}
                </>
              }
            />
            <Tag color="success">{p.status}</Tag>
            {p.tx_hash && (
              <Typography.Link href={`${BLOCK_EXPLORER}/tx/${p.tx_hash}`} target="_blank">
                可信记录
              </Typography.Link>
            )}
          </List.Item>
        )}
      />
    </Card>
  );
}
