import { Navigate, useLocation } from "react-router-dom";
import { Spin } from "antd";
import { useEffect, useMemo, useState } from "react";
import { RoleCode } from "../api/auth";
import { useAuth } from "../hooks/useAuth";

interface ProtectedRouteProps {
  children: React.ReactNode;
  roles?: string[];
}

export default function ProtectedRoute({ children, roles }: ProtectedRouteProps) {
  const { user, loading, switchRole } = useAuth();
  const location = useLocation();
  const [switchingRole, setSwitchingRole] = useState(false);
  const [roleSwitchFailed, setRoleSwitchFailed] = useState(false);

  const targetRole = useMemo(() => {
    if (!user || !roles?.length || roles.includes(user.active_role)) return null;
    return roles.find((role) => user.roles?.includes(role as RoleCode)) as RoleCode | undefined;
  }, [roles, user]);

  useEffect(() => {
    setRoleSwitchFailed(false);
    if (!targetRole) return;
    let cancelled = false;
    setSwitchingRole(true);
    switchRole(targetRole)
      .catch(() => {
        if (!cancelled) {
          setRoleSwitchFailed(true);
          setSwitchingRole(false);
        }
      })
      .finally(() => {
        if (!cancelled) setSwitchingRole(false);
      });
    return () => {
      cancelled = true;
    };
  }, [switchRole, targetRole]);

  if (loading || switchingRole) {
    return (
      <div style={{ display: "flex", justifyContent: "center", padding: 80 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!user) {
    return <Navigate to="/" replace state={{ from: location.pathname }} />;
  }

  if (roleSwitchFailed) {
    return <Navigate to="/dashboard" replace />;
  }

  if (targetRole) {
    return (
      <div style={{ display: "flex", justifyContent: "center", padding: 80 }}>
        <Spin size="large" />
      </div>
    );
  }

  if (roles && !roles.includes(user.active_role)) {
    return <Navigate to="/dashboard" replace />;
  }

  return <>{children}</>;
}
