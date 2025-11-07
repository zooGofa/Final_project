package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату для задачи в соответствии с правилом повторения
// Принимает:
//   - now: время, от которого ищется ближайшая дата
//   - dstart: исходное время в формате 20060102, от которого начинается отсчёт повторений
//   - repeat: правило повторения
//
// Поддерживаемые правила:
//   - "d <число>": повторение через указанное количество дней (1-400)
//   - "y": ежегодное повторение
//   - "w <дни недели>": повторение в указанные дни недели (1-7, где 1-понедельник, 7-воскресенье)
//   - "m <дни месяца> [месяцы]": повторение в указанные дни месяца (1-31, -1, -2)
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Если правило не указано, возвращаем пустую строку
	if strings.TrimSpace(repeat) == "" {
		return "", nil
	}

	// Парсим исходную дату
	date, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата dstart: %w", err)
	}

	// Разбираем строку с правилом повторения в слайс
	rep := strings.Split(repeat, " ")

	// Если только один символ и это не год, то плохое правило
	if rep[0] != "y" && len(rep) < 2 {
		return "", fmt.Errorf("неизвестное правило повторения: %s", rep[0])
	}

	// Повторяем год
	if rep[0] == "y" {
		for {
			date = date.AddDate(1, 0, 0)
			if date.After(now) {
				break
			}
		}
		return date.Format("20060102"), nil
	}

	// Повторяем дни
	if rep[0] == "d" {
		interval, err := strconv.Atoi(rep[1])
		if err != nil {
			return "", fmt.Errorf("невалидное число дней: %s", rep[1])
		}
		if interval > 400 || interval < 1 {
			return "", fmt.Errorf("интервал должен быть от 1 до 400, получено: %d", interval)
		}
		for {
			date = date.AddDate(0, 0, interval)
			if date.After(now) {
				break
			}
		}
		return date.Format("20060102"), nil
	}

	// Повторяем дни недели
	if rep[0] == "w" {
		weekDays := strings.Split(rep[1], ",")

		for {
			// Приращаем дату
			date = date.AddDate(0, 0, 1)
			weekDayMatch := false
			for _, v := range weekDays {
				weekDayNum, err := strconv.Atoi(strings.TrimSpace(v))
				if err != nil {
					return "", fmt.Errorf("невалидный день недели: %s", v)
				}
				// Если левый номер дня
				if weekDayNum > 7 || weekDayNum < 0 {
					return "", fmt.Errorf("день недели должен быть от 1 до 7, получено: %d", weekDayNum)
				}
				// Воскресенье в буржуйский формат
				if weekDayNum == 7 {
					weekDayNum = 0
				}
				// Если нашли
				if date.Weekday() == time.Weekday(weekDayNum) {
					weekDayMatch = true
					break
				}
			}
			// Все совпало - прекращаем поиск
			if date.After(now) && weekDayMatch {
				break
			}
		}
		return date.Format("20060102"), nil
	}

	// Повторяем дни месяца
	if rep[0] == "m" {
		monthDays := strings.Split(rep[1], ",")
		// Парсим дни месяца (сохраняем оригинальные строки для проверки формата)
		dayNums := make([]int, 0, len(monthDays))
		for _, d := range monthDays {
			d = strings.TrimSpace(d)
			day, err := strconv.Atoi(d)
			if err != nil {
				return "", fmt.Errorf("невалидный день месяца: %s", d)
			}
			if day > 31 || day < -2 || day == 0 {
				return "", fmt.Errorf("день месяца должен быть от 1 до 31, -1 или -2, получено: %d", day)
			}
			dayNums = append(dayNums, day)
		}

		// Парсим месяцы, если указаны
		monthNums := make(map[int]bool)
		if len(rep) > 2 {
			monthStrings := strings.Split(rep[2], ",")
			for _, m := range monthStrings {
				month, err := strconv.Atoi(strings.TrimSpace(m))
				if err != nil {
					return "", fmt.Errorf("невалидный месяц: %s", m)
				}
				if month < 1 || month > 12 {
					return "", fmt.Errorf("месяц должен быть от 1 до 12, получено: %d", month)
				}
				monthNums[month] = true
			}
		}

		// Нормализуем даты (убираем время)
		dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
		nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		// Начинаем с исходной даты или текущей (что больше)
		searchDate := dateOnly
		if nowOnly.After(dateOnly) || nowOnly.Equal(dateOnly) {
			searchDate = nowOnly.AddDate(0, 0, 1) // Начинаем со следующего дня после now
		}

		// Ищем ближайшую подходящую дату
		maxIterations := 800 // Ограничение на количество итераций
		for i := 0; i < maxIterations; i++ {
			// Проверяем, подходит ли текущая дата
			dayMatch := false
			for _, targetDay := range dayNums {
				// Определяем целевой день
				var actualDay int
				if targetDay == -1 {
					// Последний день месяца
					lastDay := time.Date(searchDate.Year(), searchDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
					actualDay = lastDay
				} else if targetDay == -2 {
					// Предпоследний день месяца
					lastDay := time.Date(searchDate.Year(), searchDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
					actualDay = lastDay - 1
				} else {
					actualDay = targetDay
				}

				// Проверяем, что этот день существует в месяце
				maxDayInMonth := time.Date(searchDate.Year(), searchDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
				if actualDay <= maxDayInMonth && searchDate.Day() == actualDay {
					dayMatch = true
					break
				}
			}

			// Проверяем месяц, если задана фильтрация по месяцам
			monthMatch := true
			if len(monthNums) > 0 {
				currentMonth := int(searchDate.Month())
				monthMatch = monthNums[currentMonth]
			}

			// Если всё совпало и дата после now - возвращаем
			if dayMatch && monthMatch && searchDate.After(nowOnly) {
				return searchDate.Format("20060102"), nil
			}

			// Переходим к следующему дню
			searchDate = searchDate.AddDate(0, 0, 1)
		}

		return "", fmt.Errorf("не удалось найти следующую дату")
	}

	return "", nil
}
