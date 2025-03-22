package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep = 0.65 // средняя длина шага.
	mInKm   = 1000 // количество метров в километре.
	minInH  = 60   // количество минут в часе.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	var (
		steps        int
		trainingType string
		duration     time.Duration
		err          error
	)

	arr := strings.Split(data, ",")
	if len(arr) == 3 {
		steps, err = strconv.Atoi(arr[0])
		if err != nil || steps <= 0 {
			return 0, "", 0, err
		}
		trainingType = arr[1]
		duration, err = time.ParseDuration(arr[2])
		if err != nil {
			return 0, "", 0, err
		}
	} else {
		return 0, "", 0, errors.New("Invalid input")
	}
	return steps, trainingType, duration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return (float64(steps) * lenStep) / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= 0 || steps <= 0 {
		return 0
	}
	return distance(steps) / float64(duration.Hours())
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, trainType, duration, err := parseTraining(data)
	if err != nil {
		errors.New("Input error")
	}
	distance := distance(steps)
	meanSpeed := meanSpeed(steps, duration)
	var calories float64
	switch trainType {
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		calories = RunningSpentCalories(steps, weight, duration)
	default:
		return "неизвестный тип тренировки\n"
	}
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч. Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\nСожгли калорий: %.2f\n", trainType, duration.Hours(), distance, meanSpeed, calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	return ((runningCaloriesMeanSpeedMultiplier * meanSpeed(steps, duration)) -
		runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	return ((walkingCaloriesWeightMultiplier * weight) + (meanSpeed(steps, duration)*meanSpeed(steps, duration)/height)*
		walkingSpeedHeightMultiplier) * duration.Hours() * minInH
}
