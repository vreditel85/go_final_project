package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// nextDateHandler обрабатывает запросы для вычисления следующей даты
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Разрешаем CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	// Проверяем метод запроса
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры с помощью FormValue
	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	// Проверяем обязательные параметры
	if dateParam == "" {
		http.Error(w, "Параметр 'date' обязателен", http.StatusBadRequest)
		return
	}
	if repeatParam == "" {
		http.Error(w, "Параметр 'repeat' обязателен", http.StatusBadRequest)
		return
	}

	// Определяем текущее время (now)
	var nowTime time.Time
	if nowParam == "" {
		// Если now не указан, используем текущую дату
		nowTime = time.Now()
	} else {
		// Парсим переданную дату now
		parsedNow, err := time.Parse("20060102", nowParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("Некорректный формат параметра 'now': %v", err), http.StatusBadRequest)
			return
		}
		nowTime = parsedNow
	}

	// Вызываем функцию NextDate
	nextDate, err := NextDate(nowTime, dateParam, repeatParam)
	if err != nil {
		http.Error(w, fmt.Sprintf("Ошибка вычисления следующей даты: %v", err), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым")
	}

	// Парсим исходную дату
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат исходной даты: %v", err)
	}

	// Приводим время к началу дня для корректного сравнения
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	switch {
	case strings.HasPrefix(repeat, "d "):
		// Ежедневное повторение с интервалом
		parts := strings.Fields(repeat)
		if len(parts) != 2 {
			return "", fmt.Errorf("некорректный формат правила 'd': %s", repeat)
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("некорректное число дней: %s", parts[1])
		}

		if days < 1 || days > 400 {
			return "", fmt.Errorf("число дней должно быть от 1 до 400: %d", days)
		}

		// Вычисляем следующую дату
		nextDate := startDate
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil

	case strings.HasPrefix(repeat, "y"):
		// Ежегодное повторение
		nextDate := startDate
		for !nextDate.After(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", repeat)
	}
}
