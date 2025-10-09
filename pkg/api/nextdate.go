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

	// Нормализуем время
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	// Обрабатываем правила повторения
	switch {
	case strings.HasPrefix(repeat, "d "):
		// Ежедневное повторение с интервалом
		parts := strings.Split(repeat, " ")
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

		// ВСЕГДА начинаем с startDate + интервал
		nextDate := startDate.AddDate(0, 0, days)

		// Если nextDate все еще не после now, продолжаем прибавлять
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil

	case repeat == "y":
		// Ежегодное повторение - ВСЕГДА прибавляем минимум 1 год
		nextDate := startDate

		// ВСЕГДА прибавляем минимум 1 год
		nextYear := nextDate.Year() + 1

		// Обработка 29 февраля
		if nextDate.Month() == 2 && nextDate.Day() == 29 {
			if isLeapYear(nextYear) {
				nextDate = time.Date(nextYear, 2, 29, 0, 0, 0, 0, nextDate.Location())
			} else {
				nextDate = time.Date(nextYear, 3, 1, 0, 0, 0, 0, nextDate.Location())
			}
		} else {
			nextDate = time.Date(nextYear, nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, nextDate.Location())
		}

		// Если nextDate все еще не после now, продолжаем прибавлять годы
		for nextDate.Before(now) || nextDate.Equal(now) {
			nextYear := nextDate.Year() + 1

			if nextDate.Month() == 2 && nextDate.Day() == 29 {
				if isLeapYear(nextYear) {
					nextDate = time.Date(nextYear, 2, 29, 0, 0, 0, 0, nextDate.Location())
				} else {
					nextDate = time.Date(nextYear, 3, 1, 0, 0, 0, 0, nextDate.Location())
				}
			} else {
				nextDate = time.Date(nextYear, nextDate.Month(), nextDate.Day(), 0, 0, 0, 0, nextDate.Location())
			}
		}

		return nextDate.Format("20060102"), nil

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", repeat)
	}
}

// isLeapYear проверяет, является ли год високосным
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// NextDateForTask вычисляет следующую дату для задач (для добавления и завершения задач)
func NextDateForTask(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не может быть пустым")
	}

	// Парсим исходную дату
	startDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("некорректный формат исходной даты: %v", err)
	}

	// Нормализуем время
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())

	// Обрабатываем правила повторения
	switch {
	case strings.HasPrefix(repeat, "d "):
		parts := strings.Split(repeat, " ")
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

		// Для задач: всегда прибавляем интервал от startDate
		nextDate := startDate.AddDate(0, 0, days)
		return nextDate.Format("20060102"), nil

	case repeat == "y":
		nextYear := startDate.Year() + 1

		// Обработка 29 февраля
		if startDate.Month() == 2 && startDate.Day() == 29 {
			if isLeapYear(nextYear) {
				nextDate := time.Date(nextYear, 2, 29, 0, 0, 0, 0, startDate.Location())
				return nextDate.Format("20060102"), nil
			} else {
				nextDate := time.Date(nextYear, 3, 1, 0, 0, 0, 0, startDate.Location())
				return nextDate.Format("20060102"), nil
			}
		} else {
			nextDate := time.Date(nextYear, startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
			return nextDate.Format("20060102"), nil
		}

	default:
		return "", fmt.Errorf("неподдерживаемый формат правила повторения: %s", repeat)
	}
}
