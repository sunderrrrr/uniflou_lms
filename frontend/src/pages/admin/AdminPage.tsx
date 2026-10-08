import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { useAuthStore } from "@/features/auth/store";
import {
  createAdminSubject,
  createAdminTopic,
  createTutor,
  deleteAdminTask,
  deleteUser,
  fetchAdminSubjects,
  fetchAdminTasks,
  fetchAdminUsers,
  saveAdminTask,
  updateUserRole,
  type AdminTask,
  type AdminTaskInput,
  type AdminUser,
} from "@/features/admin/api";
import styles from "./AdminPage.module.css";

type AdminTab = "users" | "content";

const emptyTask: AdminTaskInput = {
  subject_id: "",
  topic_id: "",
  task_number: 1,
  type: "short_answer",
  question: "",
  options: [],
  correct_answer: "",
  explanation: "",
  source: "",
};

function getErrorMessage(error: unknown) {
  return error instanceof Error
    ? error.message
    : "Не удалось выполнить операцию";
}

export function AdminPage() {
  const [tab, setTab] = useState<AdminTab>("users");
  const [tutorEmail, setTutorEmail] = useState("");
  const [tutorPassword, setTutorPassword] = useState("");
  const [subjectName, setSubjectName] = useState("");
  const [subjectCode, setSubjectCode] = useState("");
  const [topicName, setTopicName] = useState("");
  const [topicSubjectID, setTopicSubjectID] = useState("");
  const [taskInput, setTaskInput] = useState<AdminTaskInput>(emptyTask);
  const [optionsText, setOptionsText] = useState("");
  const [editingTaskID, setEditingTaskID] = useState<string | null>(null);
  const [notice, setNotice] = useState("");
  const currentUserID = useAuthStore((state) => state.user?.id);
  const queryClient = useQueryClient();
  const usersQuery = useQuery({
    queryKey: ["admin-users"],
    queryFn: fetchAdminUsers,
    enabled: tab === "users",
  });
  const subjectsQuery = useQuery({
    queryKey: ["admin-subjects"],
    queryFn: fetchAdminSubjects,
  });
  const tasksQuery = useQuery({
    queryKey: ["admin-tasks"],
    queryFn: fetchAdminTasks,
    enabled: tab === "content",
  });

  const invalidateUsers = () =>
    queryClient.invalidateQueries({ queryKey: ["admin-users"] });
  const invalidateTasks = () =>
    queryClient.invalidateQueries({ queryKey: ["admin-tasks"] });
  const tutorMutation = useMutation({
    mutationFn: () => createTutor(tutorEmail.trim(), tutorPassword),
    onSuccess: async () => {
      setTutorEmail("");
      setTutorPassword("");
      setNotice("Наставник создан");
      await invalidateUsers();
    },
  });
  const roleMutation = useMutation({
    mutationFn: ({ id, role }: { id: string; role: AdminUser["role"] }) =>
      updateUserRole(id, role),
    onSuccess: invalidateUsers,
  });
  const userDeleteMutation = useMutation({
    mutationFn: deleteUser,
    onSuccess: invalidateUsers,
  });
  const subjectMutation = useMutation({
    mutationFn: () =>
      createAdminSubject({ name: subjectName, code: subjectCode }),
    onSuccess: async () => {
      setSubjectName("");
      setSubjectCode("");
      setNotice("Предмет создан");
      await queryClient.invalidateQueries({ queryKey: ["admin-subjects"] });
    },
  });
  const topicMutation = useMutation({
    mutationFn: () =>
      createAdminTopic({ subject_id: topicSubjectID, name: topicName }),
    onSuccess: async () => {
      setTopicName("");
      setNotice("Тема создана");
      await queryClient.invalidateQueries({ queryKey: ["admin-subjects"] });
    },
  });
  const taskMutation = useMutation({
    mutationFn: () =>
      saveAdminTask(
        {
          ...taskInput,
          options:
            taskInput.type === "single_choice"
              ? optionsText
                  .split("\n")
                  .map((option) => option.trim())
                  .filter(Boolean)
              : [],
        },
        editingTaskID ?? undefined,
      ),
    onSuccess: async () => {
      setTaskInput(emptyTask);
      setOptionsText("");
      setEditingTaskID(null);
      setNotice("Задание сохранено");
      await invalidateTasks();
    },
  });
  const taskDeleteMutation = useMutation({
    mutationFn: deleteAdminTask,
    onSuccess: invalidateTasks,
  });

  function editTask(task: AdminTask) {
    setEditingTaskID(task.id);
    setTaskInput({
      subject_id: task.subject_id,
      topic_id: task.topic_id,
      task_number: task.task_number,
      type: task.type,
      question: task.question,
      options: task.options ?? [],
      correct_answer: task.correct_answer,
      explanation: task.explanation ?? "",
      source: task.source ?? "",
    });
    setOptionsText((task.options ?? []).join("\n"));
    setNotice("");
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  function confirmDeleteUser(user: AdminUser) {
    if (
      window.confirm(
        `Удалить пользователя ${user.email}? Связанные данные также будут удалены.`,
      )
    ) {
      userDeleteMutation.mutate(user.id);
    }
  }

  function confirmDeleteTask(task: AdminTask) {
    if (
      window.confirm("Удалить задание? Попытки по нему также будут удалены.")
    ) {
      taskDeleteMutation.mutate(task.id);
    }
  }

  const selectedSubject = subjectsQuery.data?.find(
    (subject) => subject.id === taskInput.subject_id,
  );
  const error =
    tutorMutation.error ??
    roleMutation.error ??
    userDeleteMutation.error ??
    subjectMutation.error ??
    topicMutation.error ??
    taskMutation.error ??
    taskDeleteMutation.error;

  return (
    <div className={styles.page}>
      <header className={styles.heading}>
        <div>
          <p className={styles.eyebrow}>Управление платформой</p>
          <h1>Администрирование</h1>
        </div>
        <span className={styles.account}>Администратор</span>
      </header>

      <div
        className={styles.tabs}
        role="tablist"
        aria-label="Разделы администрирования"
      >
        <button
          type="button"
          role="tab"
          aria-selected={tab === "users"}
          className={tab === "users" ? styles.activeTab : ""}
          onClick={() => {
            setTab("users");
            setNotice("");
          }}
        >
          Пользователи
        </button>
        <button
          type="button"
          role="tab"
          aria-selected={tab === "content"}
          className={tab === "content" ? styles.activeTab : ""}
          onClick={() => {
            setTab("content");
            setNotice("");
          }}
        >
          Контент
        </button>
      </div>

      {error ? (
        <p className={styles.error} role="alert">
          {getErrorMessage(error)}
        </p>
      ) : null}
      {notice ? (
        <p className={styles.notice} role="status">
          {notice}
        </p>
      ) : null}

      {tab === "users" ? (
        <section className={styles.section}>
          <div className={styles.sectionHeading}>
            <div>
              <h2>Пользователи</h2>
              <p>Роли и доступ к рабочим разделам</p>
            </div>
            <span>{usersQuery.data?.total ?? 0} всего</span>
          </div>
          <form
            className={styles.createTutor}
            onSubmit={(event) => {
              event.preventDefault();
              setNotice("");
              tutorMutation.mutate();
            }}
          >
            <h3>Добавить наставника</h3>
            <label>
              Email
              <input
                type="email"
                autoComplete="email"
                required
                value={tutorEmail}
                onChange={(event) => setTutorEmail(event.target.value)}
              />
            </label>
            <label>
              Временный пароль
              <input
                type="password"
                autoComplete="new-password"
                minLength={8}
                required
                value={tutorPassword}
                onChange={(event) => setTutorPassword(event.target.value)}
              />
            </label>
            <button type="submit" disabled={tutorMutation.isPending}>
              Создать наставника
            </button>
          </form>
          {usersQuery.isLoading ? (
            <p className={styles.muted}>Загружаем пользователей…</p>
          ) : usersQuery.isError ? (
            <p className={styles.error}>{getErrorMessage(usersQuery.error)}</p>
          ) : (
            <div className={styles.tableWrap}>
              <table>
                <thead>
                  <tr>
                    <th>Email</th>
                    <th>Роль</th>
                    <th>Почта</th>
                    <th>Создан</th>
                    <th>
                      <span className={styles.srOnly}>Действия</span>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {(usersQuery.data?.items ?? []).map((user) => (
                    <tr key={user.id}>
                      <td>
                        {user.email}
                        {user.id === currentUserID ? (
                          <span className={styles.you}>Вы</span>
                        ) : null}
                      </td>
                      <td>
                        <select
                          aria-label={`Роль ${user.email}`}
                          value={user.role}
                          disabled={
                            user.id === currentUserID || roleMutation.isPending
                          }
                          onChange={(event) =>
                            roleMutation.mutate({
                              id: user.id,
                              role: event.target.value as AdminUser["role"],
                            })
                          }
                        >
                          <option value="student">Ученик</option>
                          <option value="parent">Родитель</option>
                          <option value="tutor">Наставник</option>
                          <option value="admin">Администратор</option>
                        </select>
                      </td>
                      <td>
                        {user.email_verified
                          ? "Подтверждена"
                          : "Не подтверждена"}
                      </td>
                      <td>
                        {new Date(user.created_at).toLocaleDateString("ru-RU")}
                      </td>
                      <td>
                        <button
                          className={styles.dangerButton}
                          type="button"
                          disabled={
                            user.id === currentUserID ||
                            userDeleteMutation.isPending
                          }
                          onClick={() => confirmDeleteUser(user)}
                        >
                          Удалить
                        </button>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
              {(usersQuery.data?.items ?? []).length === 0 ? (
                <p className={styles.muted}>Пользователей пока нет.</p>
              ) : null}
            </div>
          )}
        </section>
      ) : (
        <section className={styles.section}>
          <div className={styles.sectionHeading}>
            <div>
              <h2>Банк заданий</h2>
              <p>Предметы, типы заданий и правильные ответы</p>
            </div>
            <span>{tasksQuery.data?.total ?? 0} всего</span>
          </div>
          <div className={styles.catalogForms}>
            <form
              className={styles.compactForm}
              onSubmit={(event) => {
                event.preventDefault();
                setNotice("");
                subjectMutation.mutate();
              }}
            >
              <h3>Новый предмет</h3>
              <label>
                Название
                <input
                  required
                  maxLength={100}
                  value={subjectName}
                  onChange={(event) => setSubjectName(event.target.value)}
                />
              </label>
              <label>
                Код
                <input
                  required
                  maxLength={20}
                  pattern="[A-Za-z0-9_-]+"
                  value={subjectCode}
                  onChange={(event) => setSubjectCode(event.target.value)}
                />
              </label>
              <button type="submit" disabled={subjectMutation.isPending}>
                Добавить предмет
              </button>
            </form>
            <form
              className={styles.compactForm}
              onSubmit={(event) => {
                event.preventDefault();
                setNotice("");
                topicMutation.mutate();
              }}
            >
              <h3>Новая тема</h3>
              <label>
                Предмет
                <select
                  required
                  value={topicSubjectID}
                  onChange={(event) => setTopicSubjectID(event.target.value)}
                >
                  <option value="">Выберите предмет</option>
                  {subjectsQuery.data?.map((subject) => (
                    <option key={subject.id} value={subject.id}>
                      {subject.name}
                    </option>
                  ))}
                </select>
              </label>
              <label>
                Название темы
                <input
                  required
                  maxLength={255}
                  value={topicName}
                  onChange={(event) => setTopicName(event.target.value)}
                />
              </label>
              <button
                type="submit"
                disabled={
                  topicMutation.isPending || !subjectsQuery.data?.length
                }
              >
                Добавить тему
              </button>
            </form>
          </div>
          <form
            className={styles.taskForm}
            onSubmit={(event) => {
              event.preventDefault();
              setNotice("");
              taskMutation.mutate();
            }}
          >
            <div className={styles.formTitle}>
              <h3>
                {editingTaskID ? "Редактировать задание" : "Новое задание"}
              </h3>
              {editingTaskID ? (
                <button
                  type="button"
                  className={styles.textButton}
                  onClick={() => {
                    setEditingTaskID(null);
                    setTaskInput(emptyTask);
                    setOptionsText("");
                  }}
                >
                  Отменить
                </button>
              ) : null}
            </div>
            <div className={styles.formGrid}>
              <label>
                Предмет
                <select
                  required
                  value={taskInput.subject_id}
                  onChange={(event) => {
                    const subject = subjectsQuery.data?.find(
                      (item) => item.id === event.target.value,
                    );
                    setTaskInput((current) => ({
                      ...current,
                      subject_id: event.target.value,
                      topic_id: subject?.topics[0]?.id ?? "",
                    }));
                  }}
                >
                  <option value="">Выберите предмет</option>
                  {subjectsQuery.data?.map((subject) => (
                    <option key={subject.id} value={subject.id}>
                      {subject.name}
                    </option>
                  ))}
                </select>
              </label>
              <label>
                Тема
                <select
                  required
                  value={taskInput.topic_id}
                  onChange={(event) =>
                    setTaskInput((current) => ({
                      ...current,
                      topic_id: event.target.value,
                    }))
                  }
                >
                  <option value="">Выберите тему</option>
                  {selectedSubject?.topics.map((topic) => (
                    <option key={topic.id} value={topic.id}>
                      {topic.name}
                    </option>
                  ))}
                </select>
              </label>
              <label>
                Номер в тренажёре
                <input
                  type="number"
                  min={1}
                  required
                  value={taskInput.task_number}
                  onChange={(event) =>
                    setTaskInput((current) => ({
                      ...current,
                      task_number: Number(event.target.value),
                    }))
                  }
                />
              </label>
              <label>
                Тип
                <select
                  value={taskInput.type}
                  onChange={(event) =>
                    setTaskInput((current) => ({
                      ...current,
                      type: event.target.value as AdminTaskInput["type"],
                    }))
                  }
                >
                  <option value="short_answer">Краткий ответ</option>
                  <option value="single_choice">Один вариант</option>
                </select>
              </label>
            </div>
            <label>
              Условие
              <textarea
                required
                rows={3}
                value={taskInput.question}
                onChange={(event) =>
                  setTaskInput((current) => ({
                    ...current,
                    question: event.target.value,
                  }))
                }
              />
            </label>
            {taskInput.type === "single_choice" ? (
              <label>
                Варианты ответа, каждый с новой строки
                <textarea
                  required
                  rows={4}
                  value={optionsText}
                  onChange={(event) => setOptionsText(event.target.value)}
                />
              </label>
            ) : null}
            <div className={styles.formGrid}>
              <label>
                Правильный ответ
                <input
                  required
                  value={taskInput.correct_answer}
                  onChange={(event) =>
                    setTaskInput((current) => ({
                      ...current,
                      correct_answer: event.target.value,
                    }))
                  }
                />
              </label>
              <label>
                Источник
                <input
                  value={taskInput.source}
                  onChange={(event) =>
                    setTaskInput((current) => ({
                      ...current,
                      source: event.target.value,
                    }))
                  }
                />
              </label>
            </div>
            <label>
              Пояснение
              <textarea
                rows={2}
                value={taskInput.explanation}
                onChange={(event) =>
                  setTaskInput((current) => ({
                    ...current,
                    explanation: event.target.value,
                  }))
                }
              />
            </label>
            <button
              type="submit"
              disabled={taskMutation.isPending || subjectsQuery.isLoading}
            >
              {taskMutation.isPending
                ? "Сохраняем…"
                : editingTaskID
                  ? "Сохранить изменения"
                  : "Добавить задание"}
            </button>
          </form>
          {tasksQuery.isLoading ? (
            <p className={styles.muted}>Загружаем задания…</p>
          ) : tasksQuery.isError ? (
            <p className={styles.error}>{getErrorMessage(tasksQuery.error)}</p>
          ) : (
            <div className={styles.tableWrap}>
              <table>
                <thead>
                  <tr>
                    <th>Предмет / тема</th>
                    <th>№</th>
                    <th>Условие</th>
                    <th>Ответ</th>
                    <th>
                      <span className={styles.srOnly}>Действия</span>
                    </th>
                  </tr>
                </thead>
                <tbody>
                  {(tasksQuery.data?.items ?? []).map((task) => {
                    const subject = subjectsQuery.data?.find(
                      (item) => item.id === task.subject_id,
                    );
                    const topic = subject?.topics.find(
                      (item) => item.id === task.topic_id,
                    );
                    return (
                      <tr key={task.id}>
                        <td>
                          {subject?.name ?? "Предмет"}
                          <span className={styles.subCell}>
                            {topic?.name ?? "Тема"}
                          </span>
                        </td>
                        <td>{task.task_number}</td>
                        <td className={styles.questionCell}>{task.question}</td>
                        <td>{task.correct_answer}</td>
                        <td className={styles.actions}>
                          <button
                            className={styles.textButton}
                            type="button"
                            onClick={() => editTask(task)}
                          >
                            Изменить
                          </button>
                          <button
                            className={styles.dangerButton}
                            type="button"
                            disabled={taskDeleteMutation.isPending}
                            onClick={() => confirmDeleteTask(task)}
                          >
                            Удалить
                          </button>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
              {(tasksQuery.data?.items ?? []).length === 0 ? (
                <p className={styles.muted}>Заданий пока нет.</p>
              ) : null}
            </div>
          )}
        </section>
      )}
    </div>
  );
}
