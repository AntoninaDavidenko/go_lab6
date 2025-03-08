package main

import (
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type ElectroDevice struct {
	Name                string
	Efficiency          float64
	PowerFactor         float64
	Voltage             float64
	DevicesCount        int
	Power               int
	UtilizationFactor   float64
	ReactivePowerFactor float64
	Multiplication      int
	CalculatedCurrent   float64
}

type CalculationResult struct {
	Devices                     []ElectroDevice
	GroupUtilizationRate        float64
	EffectiveAmountED           float64
	KP                          float64
	CalculatedActiveLoad        float64
	CalculatedReactiveLoad      float64
	FullPower                   float64
	EstimatedGroupCurrent       float64
	WorkshopUtilizationRates    float64
	EffectiveNumberEDWorkshop   float64
	KP2                         float64
	CalculatedActiveLoadBuses   float64
	CalculatedReactiveLoadBuses float64
	FullPowerBuses              float64
	EstimatedGroupCurrentBuses  float64
}

func calculate(devices []ElectroDevice) CalculationResult {
	for i := range devices {
		devices[i].Multiplication = devices[i].Power * devices[i].DevicesCount
		devices[i].CalculatedCurrent = float64(devices[i].Multiplication) /
			(math.Sqrt(3) * devices[i].Voltage * devices[i].PowerFactor * devices[i].Efficiency)
	}

	var nPnKv float64
	var nPn float64
	var nPn2 float64
	var kPkBkN float64

	for _, device := range devices {
		nPnKv += float64(device.Power*device.DevicesCount) * device.UtilizationFactor
		nPn += float64(device.Power * device.DevicesCount)
		nPn2 += float64(device.DevicesCount) * math.Pow(float64(device.Power), 2)
		kPkBkN += float64(device.DevicesCount) * float64(device.Power) *
			device.UtilizationFactor * device.ReactivePowerFactor
	}

	groupUtilizationRate := nPnKv / nPn
	effectiveAmountED := math.Pow(nPn, 2) / nPn2
	kP := 1.25
	calculatedActiveLoad := kP * nPnKv
	calculatedReactiveLoad := 1.0 * kPkBkN
	fullPower := math.Sqrt(math.Pow(calculatedActiveLoad, 2) + math.Pow(calculatedReactiveLoad, 2))
	estimatedGroupCurrent := calculatedActiveLoad / devices[0].Voltage

	workshopUtilizationRates := 752.0 / 2330.0
	effectiveNumberEDWorkshop := math.Pow(2330.0, 2) / 96388
	kP2 := 0.7
	calculatedActiveLoadBuses := kP2 * 752
	calculatedReactiveLoadBuses := kP2 * 657
	fullPowerBuses := math.Sqrt(math.Pow(calculatedActiveLoadBuses, 2) +
		math.Pow(calculatedReactiveLoadBuses, 2))
	estimatedGroupCurrentBuses := calculatedActiveLoadBuses / devices[0].Voltage

	return CalculationResult{
		Devices:                     devices,
		GroupUtilizationRate:        groupUtilizationRate,
		EffectiveAmountED:           effectiveAmountED,
		KP:                          kP,
		CalculatedActiveLoad:        calculatedActiveLoad,
		CalculatedReactiveLoad:      calculatedReactiveLoad,
		FullPower:                   fullPower,
		EstimatedGroupCurrent:       estimatedGroupCurrent,
		WorkshopUtilizationRates:    workshopUtilizationRates,
		EffectiveNumberEDWorkshop:   effectiveNumberEDWorkshop,
		KP2:                         kP2,
		CalculatedActiveLoadBuses:   calculatedActiveLoadBuses,
		CalculatedReactiveLoadBuses: calculatedReactiveLoadBuses,
		FullPowerBuses:              fullPowerBuses,
		EstimatedGroupCurrentBuses:  estimatedGroupCurrentBuses,
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("index.html"))

	if r.Method == http.MethodPost {
		r.ParseForm()

		device := ElectroDevice{
			Name:                r.FormValue("name"),
			Efficiency:          parseFloat(r.FormValue("efficiency")),
			PowerFactor:         parseFloat(r.FormValue("powerFactor")),
			Voltage:             parseFloat(r.FormValue("voltage")),
			DevicesCount:        parseInt(r.FormValue("devicesCount")),
			Power:               parseInt(r.FormValue("power")),
			UtilizationFactor:   parseFloat(r.FormValue("utilizationFactor")),
			ReactivePowerFactor: parseFloat(r.FormValue("reactivePowerFactor")),
		}

		// Зберігання ЕП
		devices = append(devices, device)

		if r.FormValue("calculate") != "" {
			result := calculate(devices)
			tmpl.Execute(w, result)
			return
		}

		tmpl.Execute(w, nil)
		return
	}

	tmpl.Execute(w, nil)
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

func parseInt(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

// Змінна для ЕП
var devices []ElectroDevice

func main() {
	http.HandleFunc("/", handler)
	println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
