import {
  AuditOutlined,
  BankOutlined,
  ControlOutlined,
  DashboardOutlined,
  DatabaseOutlined,
  DownOutlined,
  FileProtectOutlined,
  FileSearchOutlined,
  FolderOpenOutlined,
  HomeOutlined,
  HistoryOutlined,
  TeamOutlined,
} from "@ant-design/icons";
import { Button, Dropdown, Space, message } from "antd";
import { useState, type ReactNode } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { RoleCode } from "../api/auth";
import { useAuth } from "../hooks/useAuth";

const roleLabels: Record<RoleCode, string> = {
  creator: "创作者",
  verifier: "核验方",
  auditor: "审核员",
  admin: "管理员",
};

const navItems: Array<{ path: string; label: string; icon: ReactNode; roles?: RoleCode[] }> = [
  { path: "/", label: "首页", icon: <HomeOutlined /> },
  { path: "/dashboard", label: "工作台", icon: <DashboardOutlined /> },
  { path: "/creator/works", label: "作品存证", icon: <FolderOpenOutlined />, roles: ["creator", "admin"] },
  { path: "/creator/applications", label: "认证申请", icon: <FileProtectOutlined />, roles: ["creator", "admin"] },
  { path: "/verify", label: "公开核验", icon: <FileSearchOutlined /> },
  { path: "/enterprise/register", label: "核验方入驻", icon: <BankOutlined />, roles: ["creator", "admin"] },
  { path: "/verifier/reports", label: "核验报告", icon: <FileSearchOutlined />, roles: ["verifier", "admin"] },
  { path: "/auditor/applications", label: "认证审核", icon: <AuditOutlined />, roles: ["auditor", "admin"] },
  { path: "/admin/statistics", label: "平台统计", icon: <ControlOutlined />, roles: ["admin"] },
  { path: "/admin/users", label: "用户角色", icon: <TeamOutlined />, roles: ["admin"] },
  { path: "/admin/chains", label: "链网络", icon: <DatabaseOutlined />, roles: ["admin"] },
  { path: "/admin/chain-events", label: "链上事件", icon: <HistoryOutlined />, roles: ["admin"] },
  { path: "/admin/audit-logs", label: "审计日志", icon: <AuditOutlined />, roles: ["admin"] },
];

function isActive(pathname: string, path: string) {
  if (path === "/") return pathname === "/";
  return pathname === path || pathname.startsWith(`${path}/`);
}

export default function AppLayout({ children }: { children: ReactNode }) {
  const { user, signIn, signOut, switchRole } = useAuth();
  const [signingIn, setSigningIn] = useState(false);
  const location = useLocation();
  const navigate = useNavigate();

  const userRoles = user?.roles || [];
  const visibleNav = navItems.filter((item) => {
    if (!user) return item.path === "/";
    if (item.path === "/enterprise/register" && userRoles.includes("verifier")) return false;
    return !item.roles || item.roles.some((role) => userRoles.includes(role));
  });
  const walletAddress = user?.wallet_address || "";
  const displayName = user ? user.nickname || (walletAddress ? `${walletAddress.slice(0, 6)}...${walletAddress.slice(-4)}` : "已登录") : "";
  const activeRoleLabel = user ? roleLabels[user.active_role] || user.active_role : "";
  const roleMenuItems = userRoles.map((role) => ({
    key: role,
    label: roleLabels[role] || role,
  }));

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

  async function handleRoleSwitch(role: RoleCode) {
    if (!user || role === user.active_role) return;
    try {
      await switchRole(role);
    } catch (err) {
      message.error((err as Error).message || "身份切换失败，请重试。");
    }
  }

  function handleSignOut() {
    signOut();
    navigate("/", { replace: true });
  }

  if (!user && location.pathname === "/") {
    return (
      <div className="app-shell public-entry-shell">
        <main className="public-entry">
          <div className="public-entry__brand" aria-label="Web3Proof">
            <span className="public-entry__logo">W3P</span>
            <h1>Web3Proof</h1>
          </div>

          <section className="public-entry__intro" aria-label="平台简介">
            <p className="public-entry__subtitle">
              <span>数字作品可信存证</span>
              <span>项目成果认证</span>
              <span>公开核验平台</span>
            </p>
            <p className="public-entry__copy">
              Web3Proof 通过文件哈希、链上摘要、审核凭证和公开核验报告，帮助创作者证明作品归属，
              也让第三方能够独立验证材料是否真实、完整、可追溯。
            </p>
          </section>

          <div className="public-entry__features" aria-label="核心能力">
            <div className="public-entry__feature">
              <strong>哈希存证</strong>
              <span>上传文件后生成不可逆摘要，链上记录关键证据。</span>
            </div>
            <div className="public-entry__feature">
              <strong>SBT 凭证</strong>
              <span>审核通过后发放不可转让凭证，绑定钱包身份。</span>
            </div>
            <div className="public-entry__feature">
              <strong>公开核验</strong>
              <span>按文件、编号或钱包地址查询可信结果。</span>
            </div>
          </div>

          <div className="public-entry__actions">
            <Button type="primary" size="large" loading={signingIn} onClick={handleSignIn}>
              连接钱包
            </Button>
            <Button size="large" onClick={() => navigate("/verify")}>
              公开核验
            </Button>
          </div>
        </main>
      </div>
    );
  }

  if (!user) {
    return (
      <div className="app-shell">
        <main className="app-content">{children}</main>
      </div>
    );
  }

  return (
    <div className="app-shell app-shell--workspace">
      <aside className="workspace-sidebar">
        <Link to="/" className="workspace-brand" aria-label="Web3Proof 首页">
          <span className="workspace-brand__logo">W3P</span>
          <span className="workspace-brand__text">
            <strong>Web3Proof</strong>
            <span>可信成果认证平台</span>
          </span>
        </Link>

        <Dropdown
          menu={{
            items: roleMenuItems,
            selectable: true,
            selectedKeys: [user.active_role],
            onClick: ({ key }) => handleRoleSwitch(key as RoleCode),
          }}
          trigger={["click"]}
          placement="bottomLeft"
        >
          <button className="workspace-role" type="button" aria-label={`当前身份：${activeRoleLabel}`}>
            <span>
              <small>当前身份</small>
              <strong>{activeRoleLabel}</strong>
            </span>
            {userRoles.length > 1 && <DownOutlined />}
          </button>
        </Dropdown>

        <nav className="workspace-nav" aria-label="页面导航">
          {visibleNav.map((item) => {
            const active = isActive(location.pathname, item.path);
            return (
              <Link key={item.path} to={item.path} className={`workspace-nav__link ${active ? "workspace-nav__link--active" : ""}`}>
                <span className="workspace-nav__icon">{item.icon}</span>
                <span>{item.label}</span>
              </Link>
            );
          })}
        </nav>
      </aside>

      <div className="workspace-main">
        <header className="workspace-topbar">
          <div>
            <strong>{displayName}</strong>
            <span>{walletAddress || "未连接钱包"}</span>
          </div>
          <Space wrap size={8}>
            <Button onClick={() => walletAddress && navigate(`/portfolio/${walletAddress}`)}>公开档案</Button>
            <Button onClick={() => navigate("/profile")}>个人资料</Button>
            <Button onClick={handleSignOut}>退出</Button>
          </Space>
        </header>
        <main className="app-content">{children}</main>
      </div>
    </div>
  );
}
