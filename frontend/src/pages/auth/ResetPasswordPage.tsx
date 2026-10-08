import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Link, useSearchParams } from "react-router-dom";
import { useEffect, useState } from "react";
import { z } from "zod";
import {
  resetPasswordConfirm,
  resetPasswordRequest,
} from "@/features/auth/api";
import { Button } from "@/shared/ui/Button/Button";
import { Card } from "@/shared/ui/Card/Card";
import { Input } from "@/shared/ui/Input/Input";
import styles from "./AuthPages.module.css";

const requestSchema = z.object({ email: z.string().email() });
const confirmSchema = z.object({
  token: z.string().min(1),
  newPassword: z
    .string()
    .min(8)
    .regex(/[A-Za-zА-Яа-яЁё]/)
    .regex(/[0-9]/),
});

export function ResetPasswordPage() {
  const [searchParams] = useSearchParams();
  const [notice, setNotice] = useState("");
  const requestForm = useForm<z.infer<typeof requestSchema>>({
    resolver: zodResolver(requestSchema),
  });
  const confirmForm = useForm<z.infer<typeof confirmSchema>>({
    resolver: zodResolver(confirmSchema),
  });
  const tokenFromLink = searchParams.get("token");

  useEffect(() => {
    if (tokenFromLink) confirmForm.setValue("token", tokenFromLink);
  }, [confirmForm, tokenFromLink]);

  return (
    <div className={styles.page}>
      <Card index={2} className={styles.card}>
        <img src="/logo.png" alt="" className={styles.logo} />
        <h1>Восстановление пароля</h1>
        <form
          className={styles.form}
          onSubmit={requestForm.handleSubmit(async (v) => {
            try {
              await resetPasswordRequest(v.email);
              setNotice(
                "Если адрес зарегистрирован, на него отправлена ссылка для сброса пароля.",
              );
              requestForm.reset({ email: v.email });
            } catch (error) {
              requestForm.setError("root", {
                message:
                  error instanceof Error
                    ? error.message
                    : "Не удалось отправить запрос",
              });
            }
          })}
        >
          <Input
            label="Email"
            type="email"
            {...requestForm.register("email")}
          />
          {requestForm.formState.errors.root ? (
            <p className={styles.error}>
              {requestForm.formState.errors.root.message}
            </p>
          ) : null}
          <Button type="submit" disabled={requestForm.formState.isSubmitting}>
            Отправить ссылку
          </Button>
        </form>
        <form
          className={styles.form}
          onSubmit={confirmForm.handleSubmit(async (v) => {
            try {
              await resetPasswordConfirm(v.token, v.newPassword);
              confirmForm.reset({ token: "", newPassword: "" });
              setNotice("Пароль обновлён. Теперь можно войти с новым паролем.");
            } catch (error) {
              confirmForm.setError("root", {
                message:
                  error instanceof Error
                    ? error.message
                    : "Не удалось обновить пароль",
              });
            }
          })}
        >
          <Input
            label="Токен из письма"
            {...confirmForm.register("token")}
            readOnly={Boolean(tokenFromLink)}
          />
          <Input
            label="Новый пароль"
            type="password"
            autoComplete="new-password"
            {...confirmForm.register("newPassword")}
          />
          {confirmForm.formState.errors.root ? (
            <p className={styles.error}>
              {confirmForm.formState.errors.root.message}
            </p>
          ) : null}
          <Button
            type="submit"
            variant="ghost"
            disabled={confirmForm.formState.isSubmitting}
          >
            Подтвердить
          </Button>
        </form>
        {notice ? (
          <p className={styles.meta} role="status">
            {notice}
          </p>
        ) : null}
        <p className={styles.meta}>
          <Link to="/login">Назад ко входу</Link>
        </p>
      </Card>
    </div>
  );
}
