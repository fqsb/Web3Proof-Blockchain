import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { getMe, getNonce, login, RoleCode, switchRole as switchRoleApi, User } from "../api/auth";

interface AuthContextValue {
  user: User | null;
  token: string | null;
  loading: boolean;
  signIn: () => Promise<void>;
  signOut: () => void;
  refreshUser: () => Promise<void>;
  switchRole: (role: RoleCode) => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);
const validRoles: RoleCode[] = ["creator", "verifier", "auditor", "admin"];

function normalizeUser(raw: User): User {
  const roles = Array.isArray(raw.roles)
    ? raw.roles.filter((role): role is RoleCode => validRoles.includes(role as RoleCode))
    : [];
  const activeRole = validRoles.includes(raw.active_role) ? raw.active_role : roles[0] || "creator";

  return {
    ...raw,
    wallet_address: raw.wallet_address || "",
    active_role: activeRole,
    roles: roles.includes(activeRole) ? roles : [activeRole, ...roles],
  };
}

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(localStorage.getItem("token"));
  const [loading, setLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    if (!localStorage.getItem("token")) {
      setUser(null);
      return;
    }
    const me = await getMe();
    setUser(normalizeUser(me));
  }, []);

  useEffect(() => {
    (async () => {
      try {
        if (token) await refreshUser();
      } catch {
        localStorage.removeItem("token");
        setToken(null);
        setUser(null);
      } finally {
        setLoading(false);
      }
    })();
  }, [token, refreshUser]);

  useEffect(() => {
    const ethereum = window.ethereum;
    if (!ethereum?.on || !ethereum?.removeListener) return;
    const clearSession = () => {
      localStorage.removeItem("token");
      setToken(null);
      setUser(null);
    };
    ethereum.on("accountsChanged", clearSession);
    ethereum.on("chainChanged", clearSession);
    return () => {
      ethereum.removeListener?.("accountsChanged", clearSession);
      ethereum.removeListener?.("chainChanged", clearSession);
    };
  }, []);

  const signIn = useCallback(async () => {
    const { connectWallet } = await import("../utils/wallet");
    const { address, signer } = await connectWallet();
    const message = await getNonce(address);
    const signature = await signer.signMessage(message);
    const result = await login(address, signature, message);
    localStorage.setItem("token", result.token);
    setToken(result.token);
    setUser(normalizeUser(result.user));
  }, []);

  const signOut = useCallback(() => {
    localStorage.removeItem("token");
    setToken(null);
    setUser(null);
  }, []);

  const switchRole = useCallback(async (role: RoleCode) => {
    const result = await switchRoleApi(role);
    localStorage.setItem("token", result.token);
    setToken(result.token);
    setUser(normalizeUser(result.user));
  }, []);

  const value = useMemo(
    () => ({ user, token, loading, signIn, signOut, refreshUser, switchRole }),
    [user, token, loading, signIn, signOut, refreshUser, switchRole]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within AuthProvider");
  return ctx;
}
