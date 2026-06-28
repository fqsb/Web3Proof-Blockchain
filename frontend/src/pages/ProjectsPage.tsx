import { Button, Card, Empty, Form, Input, List, Select, Space, Tag, Typography, message } from "antd";
import { PlusOutlined, SafetyCertificateOutlined } from "@ant-design/icons";
import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { createWork, listCategories, listWorks, Work } from "../api/works";

const statusOptions = [
  { value: "all", label: "全部作品" },
  { value: "draft", label: "草稿" },
  { value: "pending_chain", label: "待上链" },
  { value: "confirmed", label: "已存证" },
];

export default function ProjectsPage({ createMode = false }: { createMode?: boolean }) {
  const [works, setWorks] = useState<Work[]>([]);
  const [categories, setCategories] = useState<Array<{ id: number; name: string }>>([]);
  const [status, setStatus] = useState<string>("all");
  const navigate = useNavigate();

  async function load() {
    try {
      const [workList, categoryList] = await Promise.all([listWorks(), listCategories()]);
      setWorks(workList);
      setCategories(categoryList);
    } catch (err) {
      message.error((err as Error).message || "加载作品失败，请稍后重试。");
    }
  }

  useEffect(() => { load(); }, []);

  const filteredWorks = useMemo(() => status === "all" ? works : works.filter((item) => item.status === status), [works, status]);

  async function onFinish(values: Partial<Work>) {
    try {
      const work = await createWork(values);
      message.success("作品已创建，可以继续上传文件并发起链上存证。");
      navigate(`/creator/works/${work.id}`);
    } catch (err) {
      message.error((err as Error).message || "创建作品失败，请检查后重试。");
    }
  }

  return (
    <div className="inner-page">
      <div className="page-kicker">Creator Workspace</div>
      <Space className="page-title" align="center">
        <div>
          <Typography.Title level={2}>作品与成果存证</Typography.Title>
          <Typography.Paragraph>
            上传课程成果、证书、论文、代码压缩包或数字作品，生成可核验的哈希记录和链上存证编号。
          </Typography.Paragraph>
        </div>
        <Button type="primary" icon={<PlusOutlined />} onClick={() => navigate("/creator/works/create")}>
          新增作品
        </Button>
      </Space>

      {createMode && (
        <Card title="新增作品" className="section-card form-card">
          <Form layout="vertical" onFinish={onFinish}>
            <Form.Item label="作品标题" name="title" rules={[{ required: true, message: "请输入作品标题" }]}>
              <Input placeholder="例如：毕业设计系统源码、摄影作品、竞赛获奖证书" />
            </Form.Item>
            <Form.Item label="作品分类" name="category_id">
              <Select placeholder="选择分类" options={categories.map((item) => ({ value: item.id, label: item.name }))} />
            </Form.Item>
            <Form.Item label="作品说明" name="description">
              <Input.TextArea rows={4} placeholder="说明作品来源、用途、创作过程或证明材料背景" />
            </Form.Item>
            <Form.Item label="外部链接" name="external_url">
              <Input placeholder="GitHub、作品集、证书查询页等" />
            </Form.Item>
            <Button type="primary" htmlType="submit">创建并上传文件</Button>
          </Form>
        </Card>
      )}

      <Card className="section-card toolbar-card">
        <Space wrap>
          <span>状态筛选</span>
          <Select value={status} style={{ width: 160 }} onChange={setStatus} options={statusOptions} />
          <Tag color="blue">共 {filteredWorks.length} 条</Tag>
        </Space>
      </Card>

      {filteredWorks.length ? (
        <List
          grid={{ gutter: 16, xs: 1, md: 2, lg: 3 }}
          dataSource={filteredWorks}
          renderItem={(item) => (
            <List.Item>
              <Card className="work-card" title={item.title} extra={<Link to={`/creator/works/${item.id}`}>详情</Link>}>
                <p>{item.description || "暂无说明"}</p>
                <Space wrap>
                  <Tag color={item.status === "confirmed" ? "green" : "default"}>{item.status}</Tag>
                  <Tag>{item.visibility}</Tag>
                </Space>
                <div className="work-card__footer">
                  <SafetyCertificateOutlined />
                  <span>链上摘要存证</span>
                </div>
              </Card>
            </List.Item>
          )}
        />
      ) : (
        <Empty description="暂无作品" />
      )}
    </div>
  );
}
