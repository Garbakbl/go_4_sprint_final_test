package daysteps

import (
	"errors"
	"fmt"
	"github.com/Yandex-Practicum/go1fl-4-sprint-final/internal/spentcalories"
	"strconv"
	"strings"
	"time"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

func parsePackage(data string) (int, time.Duration, error) {
	var (
		steps    int
		duration time.Duration
		err      error
	)

	arr := strings.Split(data, ",")
	if len(arr) == 2 {
		steps, err = strconv.Atoi(arr[0])
		if err != nil || steps <= 0 {
			return 0, 0, err
		}
		duration, err = time.ParseDuration(arr[1])
		if err != nil {
			return 0, 0, err
		}
	} else {
		return 0, 0, errors.New("Invalid input")
	}
	return steps, duration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	distance := (StepLength * float64(steps)) / 1000
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	return fmt.Sprint("Количество шагов: %d.\nДистанция составила %f км.\nВы сожгли %f каллорий.",
		steps, distance, calories)
}
