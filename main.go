package main

import (
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"os"
)

const default_port = ":8080"

// Collector implements the Collector interface
type Collector struct {
	url string
}

// Entry is the json format returned from the smartgateways device
type Entry struct {
	//MacAddress              string    `json:"mac_address"`
	//GatewayModel            string    `json:"gateway_model"`
	//StartupTime             time.Time `json:"startup_time"`
	//FirmwareRunning         string    `json:"firmware_running"`
	//FirmwareAvailable       string    `json:"firmware_available"`
	//FirmwareUpdateAvailable string    `json:"firmware_update_available"`
	//WifiRssi                string    `json:"wifi_rssi"`
	//MqttConfigured          string    `json:"mqtt_configured"`
	HeatEnergy string `json:"heat_energy"`
	//Power                   string    `json:"power"`
	Temp1    string `json:"temp1"`
	Temp2    string `json:"temp2"`
	Tempdiff string `json:"tempdiff"`
	Flow     string `json:"flow"`
	Volume   string `json:"volume"`
	//MinflowM                string    `json:"minflow_m"`
	//MaxflowM                string    `json:"maxflow_m"`
	//MinflowdateM            string    `json:"minflowdate_m"`
	//MaxflowdateM            string    `json:"maxflowdate_m"`
	//MinpowerM               string    `json:"minpower_m"`
	//MaxpowerM               string    `json:"maxpower_m"`
	//Avgtemp2M               string    `json:"avgtemp2_m"`
	//MinpowerdateM           string    `json:"minpowerdate_m"`
	//MaxpowerdateM           string    `json:"maxpowerdate_m"`
	//MinflowY                string    `json:"minflow_y"`
	//MaxflowY                string    `json:"maxflow_y"`
	//MinflowdateY            string    `json:"minflowdate_y"`
	//MaxflowdateY            string    `json:"maxflowdate_y"`
	//MinpowerY               string    `json:"minpower_y"`
	//MaxpowerY               string    `json:"maxpower_y"`
	//Avgtemp1Y               string    `json:"avgtemp1_y"`
	//Avgtemp2Y               string    `json:"avgtemp2_y"`
	//MinpowerdateY           string    `json:"minpowerdate_y"`
	//MaxpowerdateY           string    `json:"maxpowerdate_y"`
	//Temp1Xm3                string    `json:"temp1xm3"`
	//Temp2Xm3                string    `json:"temp2xm3"`
	//Infoevent               string    `json:"infoevent"`
	//Hourcounter             string    `json:"hourcounter"`
}

// Descriptors used by the Collector.
// - kamstir_gj_total from "heat_energy"
// - kamstir_temp1_c_current from "temp1"
// - kamstir_temp2_c_current from "temp2"
// - kamstir_tempdiff_c_current from "tempdiff"
// - kamstir_flow_m3h_current from "flow"
// - kamstir_volume_m3_total from "volume"
var (
	kamstir_gj_total = prometheus.NewDesc(
		"kamstir_gj_total",
		"Total GJ consumed.",
		[]string{}, nil,
	)
	kamstir_temp1_c_current = prometheus.NewDesc(
		"kamstir_temp1_c_current",
		"Water temperature going in.",
		[]string{}, nil,
	)
	kamstir_temp2_c_current = prometheus.NewDesc(
		"kamstir_temp2_c_current",
		"Water temperature going out.",
		[]string{}, nil,
	)
	kamstir_tempdiff_c_current = prometheus.NewDesc(
		"kamstir_tempdiff_c_current",
		"Difference in temperature.",
		[]string{}, nil,
	)

	kamstir_flow_m3h_current = prometheus.NewDesc(
		"kamstir_flow_m3h_current",
		"Current water flow in m3.",
		[]string{}, nil,
	)

	kamstir_volume_m3_total = prometheus.NewDesc(
		"kamstir_volume_m3_total",
		"Total water consumed in m3.",
		[]string{}, nil,
	)
)

// Call module URL
func (kc *Collector) RemoteMetrics() Entry {
	var entry Entry
	resp, err := http.Get(kc.url)

	if err != nil {
		log.Printf("Failed calling youless: %v", err)
		return entry
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed reading response from youless: %v", err)
		return entry
	}

	if err := json.Unmarshal(body, &entry); err != nil {
		log.Printf("Failed unmarshalling response from youless: %v", err)
		return entry
	}

	return entry
}

// Describe is implemented with DescribeByCollect. That's possible because the
// Collect method will always return the same metrics with the same descriptors.
func (kc *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(kc, ch)
}

// Collect first triggers the collection of metrics at the youless URL.
func (kc *Collector) Collect(ch chan<- prometheus.Metric) {
	measurements := kc.RemoteMetrics()

	ch <- prometheus.MustNewConstMetric(kamstir_gj_total, prometheus.CounterValue, stringToFloat64(measurements.HeatEnergy))
	ch <- prometheus.MustNewConstMetric(kamstir_temp1_c_current, prometheus.GaugeValue, stringToFloat64(measurements.Temp1))
	ch <- prometheus.MustNewConstMetric(kamstir_temp2_c_current, prometheus.GaugeValue, stringToFloat64(measurements.Temp2))
	ch <- prometheus.MustNewConstMetric(kamstir_tempdiff_c_current, prometheus.GaugeValue, stringToFloat64(measurements.Tempdiff))
	ch <- prometheus.MustNewConstMetric(kamstir_flow_m3h_current, prometheus.GaugeValue, stringToFloat64(measurements.Flow))
	ch <- prometheus.MustNewConstMetric(kamstir_volume_m3_total, prometheus.CounterValue, stringToFloat64(measurements.Volume))
}

func stringToFloat64(value string) float64 {
	var f float64
	_, _ = fmt.Sscanf(value, "%f", &f)
	return f
}

func main() {
	reg := prometheus.NewRegistry()

	// Add the collector.
	reg.MustRegister(
		&Collector{url: "http://connectix_kamstir.local:82/kamst-ir/api/read"},
	)

	port := default_port
	envPort := os.Getenv("PORT")
	if envPort != "" {
		port = fmt.Sprintf(":%s", envPort)
	}

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	log.Fatal(http.ListenAndServe(port, nil))
}
