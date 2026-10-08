import { Navigate, Outlet } from "react-router-dom";
import { useAuthStore, type UserRole } from "@/features/auth/store";

interface ProtectedRouteProps {
  roles?: UserRole[];
}

export function ProtectedRoute({ roles }: ProtectedRouteProps) {
  const user = useAuthStore((s) => s.user);
  const token = useAuthStore((s) => s.accessToken);

  if (!token || !user) {
    return <Navigate to="/login" replace />;
  }

  if (roles && !roles.includes(user.role)) {
    const fallback =
      user.role === "parent"
        ? "/parent"
        : user.role === "admin"
          ? "/admin"
          : user.role === "tutor"
            ? "/tutor"
            : "/student";
    return <Navigate to={fallback} replace />;
  }

  return <Outlet />;
}
