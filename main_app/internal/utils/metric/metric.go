package metric

import (
	"time"

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

func init() {
	// Регистрируем метрику в Prometheus
	prometheus.MustRegister(cpuUsage, memoryUsage, diskUsage)
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
}
