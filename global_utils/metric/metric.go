package metric

import (
	"context"
	"log"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"

	"github.com/go-park-mail-ru/2024_2_EaglesDesigner/global_utils/logger"
)

func WriteRequestDuration(start time.Time, met *prometheus.HistogramVec, method string) {
	elapsed := time.Since(start).Seconds()
	met.WithLabelValues(method).Observe(elapsed)
}

// hardware metrics.
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
			Name: "memory_usage_mbytes",
			Help: "Current memory usage in bytes",
		},
		nil, // no labels for this metric
	)
	diskUsage = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "disk_usage_mbytes",
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

var hitCount = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "hitCount",
		Help: "countOfHits",
	},
	[]string{"callMethod"}, // no labels for this metric
)

func init() {
	// Регистрируем метрику в Prometheus
	prometheus.MustRegister(cpuUsage, memoryUsage, diskUsage, errorsCount, hitCount)
	log := logger.LoggerWithCtx(context.Background(), logger.Log)
	log.Info("Метрики ошибок зарегистрированы")
	log.Info("Метрики железа зарегистрировагы")
}

// RecordMetrics запускаем в основном сервисе.
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
				memoryUsage.WithLabelValues().Set(float64(virtualMem.Used) / 1024 / 1024) // Устанавливаем значение метрики
			}

			partitions, err := disk.Partitions(true)
			if err == nil {
				for _, partition := range partitions {
					// Получаем информацию о замере дискового пространства
					usageStat, err := disk.Usage(partition.Mountpoint)
					if err == nil {
						// Устанавливаем значение для использования диска
						diskUsage.WithLabelValues(partition.Mountpoint).Set(float64(usageStat.UsedPercent))
					}
				}
			}

			time.Sleep(2 * time.Second) // Записываем метрики каждые 2 секунды
		}
	}()

	CollectMetrics()
}

// CollectMetrics запускаем не в основном сервисе.
func CollectMetrics() {
	go func() {
		for {
			// ошибки
			errors.mu.Lock()
			for key, value := range errors.errorsMap {
				errorsCount.WithLabelValues(key.callMethod, key.statusCode).Set(float64(value))
				errors.errorsMap[key] = 0
			}
			errors.mu.Unlock()

			// все метрики
			metricStorage.mu.Lock()
			for key, value := range metricStorage.metrics {
				key.WithLabelValues().Set(float64(value))
				metricStorage.metrics[key] = 0
			}
			metricStorage.mu.Unlock()
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

func IncHit() {
	pc, _, _, _ := runtime.Caller(1)
	funcPath := runtime.FuncForPC(pc).Name()

	hitCount.WithLabelValues(funcPath).Inc()
}

type MetricStorage struct {
	metrics map[prometheus.GaugeVec]int
	mu      sync.Mutex
}

var metricStorage MetricStorage = MetricStorage{
	metrics: map[prometheus.GaugeVec]int{},
	mu:      sync.Mutex{},
}

func IncMetric(met prometheus.GaugeVec) {
	metricStorage.mu.Lock()
	defer metricStorage.mu.Unlock()
	if value, ok := metricStorage.metrics[met]; ok {
		metricStorage.metrics[met] = value + 1
	} else {
		metricStorage.metrics[met] = 1
	}

	log.Println(metricStorage.metrics)
}
