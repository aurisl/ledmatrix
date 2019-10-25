// +build hardware

package matrix

import (
	"github.com/aurisl/ledmatrix/error"
	"github.com/mcuadros/go-rpi-rgb-led-matrix"
	"github.com/deathowl/go-metrics-prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var Renderer = CreateHardwareMatrix()

func CreateHardwareMatrix() Matrix {
	matrixConfig := CreateMatrixDefaultConfiguration()

	matrix, err := rgbmatrix.NewRGBLedMatrix(matrixConfig)
	error.Fatal(err)

	prometheusRegistry := prometheus.NewRegistry()
	metricsRegistry := metrics.NewRegistry()
	pClient := NewPrometheusProvider(metricsRegistry, "led-matrix", "subsys", prometheusRegistry, 1*time.Second)
	go pClient.UpdatePrometheusMetrics()

	return matrix
}

func CreateMatrixDefaultConfiguration() *rgbmatrix.HardwareConfig {

	matrixConfig := &rgbmatrix.DefaultConfig
	matrixConfig.Rows = 32
	matrixConfig.Brightness = 50

	return matrixConfig
}
