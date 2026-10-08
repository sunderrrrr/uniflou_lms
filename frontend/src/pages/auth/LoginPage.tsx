import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { z } from "zod";
import { loginRequest } from "@/features/auth/api";
import { useAuthStore } from "@/features/auth/store";
import { Button } from "@/shared/ui/Button/Button";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import styles from "./AuthPages.module.css";

const schema = z.object({
  email: z.string().min(1, "Введите email").email("Введите корректный email"),
  password: z.string().min(8, "Минимум 8 символов"),
});

type FormValues = z.infer<typeof schema>;

export function LoginPage() {
  const navigate = useNavigate();
  const setSession = useAuthStore((s) => s.setSession);

  const {
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: "", password: "" },
  });

  async function onSubmit(values: FormValues) {
    try {
      const data = await loginRequest(values.email, values.password);
      setSession(data.tokens, data.user);
      const dest =
        data.user.role === "parent"
          ? "/parent"
          : data.user.role === "admin"
            ? "/admin"
            : data.user.role === "tutor"
              ? "/tutor"
              : "/student";
      navigate(dest);
    } catch (e) {
      setError("root", {
        message: e instanceof Error ? e.message : "Ошибка входа",
      });
    }
  }

  return (
    <div className={styles.page}>
      <Card index={0} className={styles.card}>
        <img src="/logo.png" alt="" className={styles.logo} />
        <h1>Вход</h1>
        <form className={styles.form} onSubmit={handleSubmit(onSubmit)}>
          <Controller
            name="email"
            control={control}
            render={({ field }) => (
              <Input
                label="Email"
                type="email"
                autoComplete="email"
                name={field.name}
                value={field.value}
                onChange={field.onChange}
                onBlur={field.onBlur}
                error={errors.email?.message}
              />
            )}
          />

          <Controller
            name="password"
            control={control}
            render={({ field }) => (
              <Input
                label="Пароль"
                type="password"
                autoComplete="current-password"
                name={field.name}
                value={field.value}
                onChange={field.onChange}
                onBlur={field.onBlur}
                error={errors.password?.message}
              />
            )}
          />

          {errors.root ? (
            <p className={styles.error}>{errors.root.message}</p>
          ) : null}

          <Button type="submit" disabled={isSubmitting}>
            Войти
          </Button>
        </form>

        <p className={styles.meta}>
          Нет аккаунта? <Link to="/register">Регистрация</Link> ·{" "}
          <Link to="/reset-password">Забыли пароль?</Link>
        </p>
      </Card>
    </div>
  );
}
