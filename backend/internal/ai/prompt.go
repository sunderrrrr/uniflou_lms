package ai

import "fmt"

func DashboardPrompt(goal int, weakTopics string) string {
	return fmt.Sprintf("Ты наставник ЕГЭ. Проанализируй прогресс ученика. Цель %%: %d. Слабые темы: %s. Дай 5 рекомендаций на русском.", goal, weakTopics)
}

func PlanPrompt(planSummary string) string {
	return fmt.Sprintf("Проанализируй план поступления: %s. Укажи слабые места и шаги на неделю.", planSummary)
}

func ProgressPrompt(summary string) string {
	return fmt.Sprintf("Сделай глубокий анализ прогресса и прогноз: %s", summary)
}

func TaskPrompt(question, answer string, correct bool) string {
	status := "верно"
	if !correct {
		status = "неверно"
	}
	return fmt.Sprintf("Разбери задание ЕГЭ. Вопрос: %s. Ответ ученика: %s (%s). Объясни ошибку и дай совет.", question, answer, status)
}
