import { Button, Card, Col, Empty, List, Row, Select, Space, Statistic, Table, Tag, Typography, message } from "antd";
import { ApiOutlined, DatabaseOutlined, HistoryOutlined, ReloadOutlined, TeamOutlined, UserSwitchOutlined } from "@ant-design/icons";
import { useEffect, useState } from "react";
import { getAdminData, syncChainEvents, updateUserRoles } from "../api/admin";
import { RoleCode } from "../api/auth";

const roleOptions = [
  { value: "creator", label: "创作者" },
  { value: "verifier", label: "核验方" },
  { value: "auditor", label: "审核员" },
  { value: "admin", label: "管理员" },
] as const;

const roleLabel: Record<RoleCode, string> = {
  creator: "创作者",
  verifier: "核验方",
  auditor: "审核员",
  admin: "管理员",
};

type AdminMode = "users" | "chains" | "chain-events" | "statistics" | "audit-logs";

export default function AdminPage({ mode = "statistics" }: { mode?: AdminMode }) {
  const [data, setData] = useState<any>(null);
  const [roleDrafts, setRoleDrafts] = useState<Record<number, RoleCode[]>>({});
  const [loading, setLoading] = useState(false);
  const [loadError, setLoadError] = useState("");

  async function load() {
    setLoading(true);
    setLoadError("");
    try {
      setData(await getAdminData(mode));
    } catch (err) {
      const text = (err as Error).message || "加载管理数据失败，请稍后重试。";
      setLoadError(text);
      message.error(text);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => { load(); }, [mode]);

  async function saveRoles(userId: number) {
    const users = Array.isArray(data) ? data : [];
    const target = users.find((item: any) => item.id === userId);
    const current = (target?.roles || []).filter((r: any) => r.enabled).map((r: any) => r.role_code || r.RoleCode);
    const roles = roleDrafts[userId] || current;
    if (!roles?.length) {
      message.error("至少保留一个角色。");
      return;
    }
    try {
      await updateUserRoles(userId, roles, roles[0]);
      message.success("用户角色已更新。");
      await load();
    } catch (err) {
      message.error((err as Error).message || "保存角色失败，请稍后重试。");
    }
  }

  async function handleSyncChainEvents() {
    setLoading(true);
    try {
      const result = await syncChainEvents();
      message.success(`同步完成：扫描 ${result.scanned} 条，新增 ${result.inserted} 条。`);
      await load();
    } catch (err) {
      message.error((err as Error).message || "同步链上事件失败，请检查 RPC 与合约配置。");
    } finally {
      setLoading(false);
    }
  }

  if (mode === "users") {
    const users = Array.isArray(data) ? data : [];
    return (
      <div className="inner-page admin-shell">
        <div className="page-kicker">Admin Console</div>
        <Typography.Title level={2}>用户与角色</Typography.Title>
        <Typography.Paragraph>管理钱包用户的多角色权限，并控制当前默认工作台身份。</Typography.Paragraph>
        {loadError ? (
          <Card>
            <Empty description={loadError} />
          </Card>
        ) : (
          <Table
            rowKey="id"
            loading={loading}
            dataSource={users}
            locale={{ emptyText: loading ? "加载中" : "暂无用户数据，请确认当前钱包拥有管理员权限。" }}
            pagination={{ pageSize: 10 }}
            scroll={{ x: 780 }}
            columns={[
              { title: "ID", dataIndex: "id", width: 80 },
              {
                title: "钱包地址",
                dataIndex: "wallet_address",
                render: (value: string) => <code>{value}</code>,
              },
              {
                title: "当前角色",
                dataIndex: "active_role",
                width: 120,
                render: (value: RoleCode) => <Tag color="blue">{roleLabel[value] || value}</Tag>,
              },
              {
                title: "角色分配",
                render: (_, record: any) => {
                  const current = (record.roles || [])
                    .filter((r: any) => r.enabled)
                    .map((r: any) => r.role_code || r.RoleCode)
                    .filter(Boolean);
                  return (
                    <Space wrap>
                      <Select
                        mode="multiple"
                        style={{ width: 280 }}
                        value={roleDrafts[record.id] || current}
                        options={[...roleOptions]}
                        onChange={(value) => setRoleDrafts((prev) => ({ ...prev, [record.id]: value as RoleCode[] }))}
                      />
                      <Button icon={<UserSwitchOutlined />} onClick={() => saveRoles(record.id)}>保存</Button>
                    </Space>
                  );
                },
              },
            ]}
          />
        )}
      </div>
    );
  }

  if (mode === "chains") {
    const chains = data?.chains || [];
    const contracts = data?.contracts || [];
    return (
      <div className="inner-page admin-shell">
        <div className="page-kicker">Chain Configuration</div>
        <Typography.Title level={2}>链网络与合约配置</Typography.Title>
        <Row gutter={[16, 16]}>
          <Col xs={24} lg={12}>
            <Card title="链网络">
              {chains.length ? chains.map((item: any) => (
                <div className="admin-record" key={item.id}>
                  <Space wrap><ApiOutlined /><strong>{item.name}</strong><Tag color={item.is_active ? "green" : "default"}>{item.code}</Tag></Space>
                  <p>Chain ID：{item.chain_id}</p>
                  <code>{item.rpc_url || "未配置 RPC"}</code>
                </div>
              )) : <Empty description="暂无链网络配置" />}
            </Card>
          </Col>
          <Col xs={24} lg={12}>
            <Card title="合约地址">
              {contracts.length ? contracts.map((item: any) => (
                <div className="admin-record" key={item.id}>
                  <Space wrap><DatabaseOutlined /><strong>{item.contract_name}</strong><Tag>{item.network_code}</Tag></Space>
                  <code>{item.contract_address}</code>
                </div>
              )) : <Empty description="暂无合约配置" />}
            </Card>
          </Col>
        </Row>
      </div>
    );
  }

  if (mode === "audit-logs") {
    return (
      <div className="inner-page admin-shell">
        <div className="page-kicker">Audit Trail</div>
        <Typography.Title level={2}>审计日志</Typography.Title>
        {data?.length ? (
          <List
            dataSource={data}
            renderItem={(item: any) => (
              <List.Item>
                <Card className="audit-card">
                  <Space direction="vertical" size={6}>
                    <Space wrap><Tag color="blue">{item.action}</Tag><span>{item.resource}</span><span>{item.created_at}</span></Space>
                    <p>{item.detail || "无详情"}</p>
                    <span>操作用户：{item.user_id || "系统"}</span>
                  </Space>
                </Card>
              </List.Item>
            )}
          />
        ) : (
          <Empty description="暂无审计日志" />
        )}
      </div>
    );
  }

  if (mode === "chain-events") {
    const events = Array.isArray(data) ? data : [];
    return (
      <div className="inner-page admin-shell">
        <div className="page-kicker">Chain Events</div>
        <div className="page-title">
          <div>
            <Typography.Title level={2}>链上事件</Typography.Title>
            <Typography.Paragraph>查看平台合约近期事件，支持从链上同步最新日志。</Typography.Paragraph>
          </div>
          <Button type="primary" icon={<ReloadOutlined />} loading={loading} onClick={handleSyncChainEvents}>
            同步链上事件
          </Button>
        </div>
        <Table
          rowKey={(record: any) => `${record.tx_hash}-${record.log_index}`}
          loading={loading}
          dataSource={events}
          locale={{ emptyText: "暂无链上事件，可点击同步链上事件拉取最近日志。" }}
          pagination={{ pageSize: 10 }}
          scroll={{ x: 980 }}
          columns={[
            { title: "合约", dataIndex: "contract_name", width: 150, render: (value: string) => <Tag color="blue">{value || "未知合约"}</Tag> },
            { title: "事件", dataIndex: "event_name", width: 160, render: (value: string) => <Space><HistoryOutlined />{value || "未知事件"}</Space> },
            { title: "区块", dataIndex: "block_number", width: 120 },
            { title: "日志序号", dataIndex: "log_index", width: 100 },
            { title: "交易哈希", dataIndex: "tx_hash", render: (value: string) => <code>{value}</code> },
            { title: "状态", dataIndex: "processed", width: 100, render: (value: boolean) => <Tag color={value ? "green" : "gold"}>{value ? "已处理" : "未处理"}</Tag> },
          ]}
        />
      </div>
    );
  }

  return (
    <div className="inner-page admin-shell">
      <div className="page-kicker">Admin Console</div>
      <Typography.Title level={2}>平台统计</Typography.Title>
      <Typography.Paragraph>查看 Web3Proof 当前用户、作品、存证与认证凭证的核心数据。</Typography.Paragraph>
      <Row gutter={[16, 16]}>
        <Col xs={12} md={6}><Card className="metric-card"><Statistic title="用户" value={data?.users || 0} prefix={<TeamOutlined />} /></Card></Col>
        <Col xs={12} md={6}><Card className="metric-card"><Statistic title="作品" value={data?.works || 0} /></Card></Col>
        <Col xs={12} md={6}><Card className="metric-card"><Statistic title="已确认存证" value={data?.confirmed_evidences || 0} /></Card></Col>
        <Col xs={12} md={6}><Card className="metric-card"><Statistic title="SBT 凭证" value={data?.credentials || 0} /></Card></Col>
      </Row>
    </div>
  );
}
