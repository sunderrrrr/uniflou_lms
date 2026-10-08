# УНИФЛОУ Тренажёр (MVP)

Персональный веб‑тренажёр подготовки к ЕГЭ: цель → план → прогресс → задания → ИИ‑анализ. Роли: ученик, родитель, наставник, admin.

## Стек

- **Backend:** Go 1.22, Gin, GORM, PostgreSQL 16, JWT, zerolog
- **Frontend:** React 18, TypeScript, Vite, Zustand, TanStack Query, CSS Modules, Recharts, Framer Motion (зависимость подключена)

## Быстрый старт

1. Скопируйте переменные окружения:

```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env
```

2. Запустите инфраструктуру:

```bash
cd backend
docker compose up --build
```

3. Примените миграции (локально или в контейнере postgres):

```bash
migrate -path ./migrations -database "postgres://uniflou:PASSWORD@localhost:5432/uniflou?sslmode=disable" up
```

4. Frontend в dev‑режиме:

```bash
cd frontend
npm ci
npm run dev
```

API: `http://localhost:8080/api/v1` · UI: `http://localhost:5173`

## Роли и сценарии

| Роль | Маршруты |
|------|----------|
| student | `/student`, `/student/plan`, `/student/trainer`, `/student/progress` |
| parent | `/parent` (+ привязка ученика по email) |
| tutor / admin | `/tutor`, `/tutor/students/:id` |

## ИИ‑анализ

Кнопка «Проанализировать через ИИ» на дашборде, плане, прогрессе и после задания. Лимит: 10 запросов/день. Без `LLM_API_KEY` работает офлайн‑заглушка.

## Логотип

Файлы: `frontend/public/logo.png`, `frontend/public/logo.svg`, favicon из логотипа.

## Структура репозитория

- `backend/` — API, миграции, Docker
- `frontend/` — SPA (Feature‑Sliced Design)

## Health

`GET /health` → `{ "status": "ok" }`
