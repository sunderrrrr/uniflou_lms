import { createBrowserRouter, Navigate } from "react-router-dom";
import { AppLayout } from "@/app/layout/AppLayout";
import { ProtectedRoute } from "@/app/ProtectedRoute";
import { LoginPage } from "@/pages/auth/LoginPage";
import { RegisterPage } from "@/pages/auth/RegisterPage";
import { ResetPasswordPage } from "@/pages/auth/ResetPasswordPage";
import { NotFoundPage } from "@/pages/NotFoundPage";
import { ParentDashboardPage } from "@/pages/parent/ParentDashboardPage";
import { DashboardPage } from "@/pages/student/DashboardPage";
import { PlanPage } from "@/pages/student/PlanPage";
import { ProgressPage } from "@/pages/student/ProgressPage";
import { TaskPage } from "@/pages/student/TaskPage";
import { TrainerPage } from "@/pages/student/TrainerPage";
import { TrainerSessionPage } from "@/pages/student/TrainerSessionPage";
import { StudentDetailPage } from "@/pages/tutor/StudentDetailPage";
import { StudentsPage } from "@/pages/tutor/StudentsPage";
import { AdminPage } from "@/pages/admin/AdminPage";
import { useAuthStore } from "@/features/auth/store";

function HomeRedirect() {
  const user = useAuthStore.getState().user;
  if (!user) return <Navigate to="/login" replace />;
  if (user.role === "parent") return <Navigate to="/parent" replace />;
  if (user.role === "admin") return <Navigate to="/admin" replace />;
  if (user.role === "tutor") return <Navigate to="/tutor" replace />;
  return <Navigate to="/student" replace />;
}

export const router = createBrowserRouter([
  { path: "/", element: <HomeRedirect /> },
  { path: "/login", element: <LoginPage /> },
  { path: "/register", element: <RegisterPage /> },
  { path: "/reset-password", element: <ResetPasswordPage /> },
  {
    element: <ProtectedRoute />,
    children: [
      {
        element: <AppLayout />,
        children: [
          {
            element: <ProtectedRoute roles={["student"]} />,
            children: [
              { path: "/student", element: <DashboardPage /> },
              { path: "/student/plan", element: <PlanPage /> },
              { path: "/student/trainer", element: <TrainerPage /> },
              {
                path: "/student/trainer/session",
                element: <TrainerSessionPage />,
              },
              { path: "/student/trainer/task/:id", element: <TaskPage /> },
              { path: "/student/progress", element: <ProgressPage /> },
            ],
          },
          {
            element: <ProtectedRoute roles={["parent"]} />,
            children: [{ path: "/parent", element: <ParentDashboardPage /> }],
          },
          {
            element: <ProtectedRoute roles={["admin"]} />,
            children: [{ path: "/admin", element: <AdminPage /> }],
          },
          {
            element: <ProtectedRoute roles={["tutor", "admin"]} />,
            children: [
              { path: "/tutor", element: <StudentsPage /> },
              { path: "/tutor/students/:id", element: <StudentDetailPage /> },
            ],
          },
        ],
      },
    ],
  },
  { path: "*", element: <NotFoundPage /> },
]);
