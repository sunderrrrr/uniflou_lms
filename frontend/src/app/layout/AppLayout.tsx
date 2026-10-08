import { motion } from "framer-motion";
import { useState } from "react";
import { Link, NavLink, Outlet, useNavigate } from "react-router-dom";
import { logoutRequest } from "@/features/auth/api";
import { useAuthStore } from "@/features/auth/store";
import { Button } from "@/shared/ui/Button/Button";
import styles from "./AppLayout.module.css";

const studentLinks = [
  { to: "/student", label: "Дашборд" },
  { to: "/student/plan", label: "Мой план" },
  { to: "/student/trainer", label: "Тренажёр" },
  { to: "/student/progress", label: "Прогресс" },
];

export function AppLayout() {
  const user = useAuthStore((s) => s.user);
  const clear = useAuthStore((s) => s.clear);
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);

  const links =
    user?.role === "parent"
      ? [{ to: "/parent", label: "Дашборд" }]
      : user?.role === "admin"
        ? [{ to: "/admin", label: "Администрирование" }]
        : user?.role === "tutor"
          ? [{ to: "/tutor", label: "Ученики" }]
          : studentLinks;

  async function onLogout() {
    try {
      await logoutRequest();
    } catch {
      /* ignore */
    }
    clear();
    navigate("/login");
  }

  return (
    <div className={styles.shell}>
      <header className={styles.header}>
        <div className={`container ${styles.headerInner}`}>
          <Link to="/" className={styles.brand}>
            <img src="/logo.png" alt="УНИФЛОУ" className={styles.logo} />
            <span>УНИФЛОУ Тренажёр</span>
          </Link>
          <button
            type="button"
            className={styles.burger}
            onClick={() => setOpen((v) => !v)}
            aria-label="Меню"
          >
            ☰
          </button>
          <nav className={`${styles.nav} ${open ? styles.navOpen : ""}`}>
            {links.map((l) => (
              <NavLink
                key={l.to}
                to={l.to}
                className={({ isActive }) => (isActive ? styles.active : "")}
                onClick={() => setOpen(false)}
              >
                {l.label}
              </NavLink>
            ))}
            <Button variant="ghost" size="sm" onClick={onLogout}>
              Выйти
            </Button>
          </nav>
        </div>
      </header>
      <motion.main
        className={`container ${styles.main}`}
        initial={{ opacity: 0, y: 12 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <Outlet />
      </motion.main>
    </div>
  );
}
