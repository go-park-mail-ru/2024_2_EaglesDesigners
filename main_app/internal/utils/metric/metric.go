package metric

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

func WriteRequestDuration(start time.Time, met *prometheus.HistogramVec, method string) {
	elapsed := time.Since(start).Seconds()
	met.WithLabelValues(method).Observe(elapsed)
}

var (
	cpuUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cpu_usage_percent",
			Help: "Current CPU usage percentage",
		},
		nil, // no labels for this metric
	)
	memoryUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "memory_usage_bytes",
			Help: "Current memory usage in bytes",
		},
		nil, // no labels for this metric
	)
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_usage_bytes",
			Help: "Current disk usage in bytes",
		},
		[]string{"mountpoint"}, // Метка для точки монтирования
	)
)

var errorsCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "errorsCount",
		Help: "count of errors",
	},
	[]string{"callMethod", "statusCode"}, // no labels for this metric
)

func init() {
	// Регистрируем метрику в Prometheus
	prometheus.MustRegister(cpuUsage, memoryUsage, diskUsage, errorsCount)
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	log.Info("Метрики ошибок зарегистрированы")
	log.Info("Метрики железа зарегистрировагы")
}

func RecordMetrics() {
	go func() {
		for {
			// Получаем использование CPU
			percent, err := cpu.Percent(0, false)
			if err == nil && len(percent) > 0 {
				cpuUsage.WithLabelValues().Set(percent[0]) // Устанавливаем значение метрики
			}

			// Получаем информацию о памяти
			virtualMem, err := mem.VirtualMemory()
			if err == nil {
				memoryUsage.WithLabelValues().Set(float64(virtualMem.Used)) // Устанавливаем значение метрики
			}

			partitions, err := disk.Partitions(true)
			if err == nil {
				for _, partition := range partitions {
					// Получаем информацию о замере дискового пространства
					usageStat, err := disk.Usage(partition.Mountpoint)
					if err == nil {
						// Устанавливаем значение для использования диска
						diskUsage.WithLabelValues(partition.Mountpoint).Set(float64(usageStat.Used))
					}
				}
			}

			time.Sleep(2 * time.Second) // Записываем метрики каждые 2 секунды
		}
	}()

	go func() {
		for {
			errors.mu.Lock()
			for key, value := range errors.errorsMap {
				errorsCount.WithLabelValues(key.callMethod, key.statusCode).Set(float64(value))
			}
			errors.errorsMap = map[ErrorLabels]int{}

			errors.mu.Unlock()
			time.Sleep(1 * time.Minute)
		}
	}()
}

type ErrorLabels struct {
	callMethod string
	statusCode string
}
type ErrorsStorage struct {
	errorsMap map[ErrorLabels]int
	mu        sync.Mutex
}

var errors ErrorsStorage = ErrorsStorage{
	errorsMap: map[ErrorLabels]int{},
	mu:        sync.Mutex{},
}

func PushError(callMethod string, statusCode int) {

	errorLabel := ErrorLabels{
		callMethod: callMethod,
		statusCode: strconv.Itoa(statusCode),
	}

	errors.mu.Lock()
	defer errors.mu.Unlock()

	if value, ok := errors.errorsMap[errorLabel]; ok {
		errors.errorsMap[errorLabel] = value + 1
	} else {
		errors.errorsMap[errorLabel] = 1
	}
}
