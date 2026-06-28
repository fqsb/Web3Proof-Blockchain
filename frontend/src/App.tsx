import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import AppLayout from "./components/AppLayout";
import PageErrorBoundary from "./components/PageErrorBoundary";
import ProtectedRoute from "./components/ProtectedRoute";
import { AuthProvider } from "./hooks/useAuth";
import AdminPage from "./pages/AdminPage";
import CertificationReviewPage from "./pages/admin/CertificationReviewPage";
import CertificationsPage from "./pages/CertificationsPage";
import DashboardPage from "./pages/DashboardPage";
import EnterpriseRegisterPage from "./pages/EnterpriseRegisterPage";
import HomePage from "./pages/HomePage";
import ProfilePage from "./pages/ProfilePage";
import ProjectDetailPage from "./pages/ProjectDetailPage";
import ProjectsPage from "./pages/ProjectsPage";
import ResumePage from "./pages/ResumePage";
import VerifyPage from "./pages/VerifyPage";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AppLayout>
          <Routes>
            <Route path="/" element={<HomePage />} />
            <Route path="/dashboard" element={<ProtectedRoute><DashboardPage /></ProtectedRoute>} />
            <Route path="/profile" element={<ProtectedRoute><ProfilePage /></ProtectedRoute>} />
            <Route path="/creator/works" element={<ProtectedRoute roles={["creator", "admin"]}><ProjectsPage /></ProtectedRoute>} />
            <Route path="/creator/works/create" element={<ProtectedRoute roles={["creator", "admin"]}><ProjectsPage createMode /></ProtectedRoute>} />
            <Route path="/creator/works/:id" element={<ProtectedRoute roles={["creator", "admin"]}><ProjectDetailPage /></ProtectedRoute>} />
            <Route path="/creator/applications" element={<ProtectedRoute roles={["creator", "admin"]}><CertificationsPage /></ProtectedRoute>} />
            <Route path="/creator/certificates" element={<ProtectedRoute roles={["creator", "admin"]}><CertificationsPage /></ProtectedRoute>} />
            <Route path="/verify" element={<VerifyPage />} />
            <Route path="/verifier/reports" element={<ProtectedRoute roles={["verifier", "admin"]}><VerifyPage reportsMode /></ProtectedRoute>} />
            <Route path="/auditor/applications" element={<ProtectedRoute roles={["auditor", "admin"]}><CertificationReviewPage /></ProtectedRoute>} />
            <Route path="/admin/users" element={<ProtectedRoute roles={["admin"]}><PageErrorBoundary><AdminPage mode="users" /></PageErrorBoundary></ProtectedRoute>} />
            <Route path="/admin/user" element={<Navigate to="/admin/users" replace />} />
            <Route path="/admin/chains" element={<ProtectedRoute roles={["admin"]}><PageErrorBoundary><AdminPage mode="chains" /></PageErrorBoundary></ProtectedRoute>} />
            <Route path="/admin/chain-events" element={<ProtectedRoute roles={["admin"]}><PageErrorBoundary><AdminPage mode="chain-events" /></PageErrorBoundary></ProtectedRoute>} />
            <Route path="/admin/statistics" element={<ProtectedRoute roles={["admin"]}><PageErrorBoundary><AdminPage mode="statistics" /></PageErrorBoundary></ProtectedRoute>} />
            <Route path="/admin/audit-logs" element={<ProtectedRoute roles={["admin"]}><PageErrorBoundary><AdminPage mode="audit-logs" /></PageErrorBoundary></ProtectedRoute>} />
            <Route path="/portfolio/:address" element={<ResumePage />} />
            <Route path="/enterprise/register" element={<ProtectedRoute><EnterpriseRegisterPage /></ProtectedRoute>} />
            <Route path="/projects" element={<Navigate to="/creator/works" replace />} />
            <Route path="/achievements" element={<Navigate to="/creator/works" replace />} />
            <Route path="/certifications" element={<Navigate to="/creator/applications" replace />} />
            <Route path="/admin" element={<Navigate to="/admin/statistics" replace />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </AppLayout>
      </BrowserRouter>
    </AuthProvider>
  );
}
