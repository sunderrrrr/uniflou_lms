import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import { z } from "zod";
import { registerRequest } from "@/features/auth/api";
import type { UserRole } from "@/features/auth/store";
import { useAuthStore } from "@/features/auth/store";
import { Button } from "@/shared/ui/Button/Button";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import styles from "./AuthPages.module.css";

const schema = z.object({
  email: z.string().min(1, "Введите email").email("Введите корректный email"),
  password: z.string().min(8, "Минимум 8 символов"),
  role: z.enum(["student", "parent"]),
});

type FormValues = z.infer<typeof schema>;

export function RegisterPage() {
  const navigate = useNavigate();
  const setSession = useAuthStore((s) => s.setSession);

  const {
    control,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
  } = useForm<FormValues>({
    resolver: zodResolver(schema),
    defaultValues: { email: "", password: "", role: "student" },
  });

  async function onSubmit(values: FormValues) {
    try {
      const data = await registerRequest(
        values.email,
        values.password,
        values.role as UserRole,
      );
      setSession(data.tokens, data.user);
      navigate(values.role === "parent" ? "/parent" : "/student");
    } catch (e) {
      setError("root", {
        message: e instanceof Error ? e.message : "Ошибка регистрации",
      });
    }
  }

  return (
    <div className={styles.page}>
      <Card index={1} className={styles.card}>
        <img src="/logo.png" alt="" className={styles.logo} />
        <h1>Регистрация</h1>
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
                autoComplete="new-password"
                name={field.name}
                value={field.value}
                onChange={field.onChange}
                onBlur={field.onBlur}
                error={errors.password?.message}
              />
            )}
          />

          <label className={styles.selectWrap}>
            Роль
            <Controller
              name="role"
              control={control}
              render={({ field }) => (
                <select
                  className={styles.select}
                  name={field.name}
                  value={field.value}
                  onChange={field.onChange}
                  onBlur={field.onBlur}
                >
                  <option value="student">Ученик</option>
                  <option value="parent">Родитель</option>
                </select>
              )}
            />
          </label>

          {errors.root ? (
            <p className={styles.error}>{errors.root.message}</p>
          ) : null}

          <Button type="submit" disabled={isSubmitting}>
            Создать аккаунт
          </Button>
        </form>

        <p className={styles.meta}>
          Уже есть аккаунт? <Link to="/login">Войти</Link>
        </p>
      </Card>
    </div>
  );
}
