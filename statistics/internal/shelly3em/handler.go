package shelly3em

import (
	"context"
	"encoding/json"
	"math"
	"time"

	"github.com/eclipse/paho.golang/paho"
)

const NotifyStatus = "NotifyStatus"

type Measurement struct {
	Time               time.Time
	Device             string
	PhaseA             PhaseMeasurement
	PhaseB             PhaseMeasurement
	PhaseC             PhaseMeasurement
	TotalActualPower   float64
	TotalApparentPower float64
	TotalCurrent       float64
}

type PhaseMeasurement struct {
	ActualPower   float64
	ApparentPower float64
	Current       float64
	Frequency     float64
	PowerFactor   float64
	Voltage       float64
}

type Counters struct {
	Time   time.Time
	Device string
	PhaseA EnergyCounters
	PhaseB EnergyCounters
	PhaseC EnergyCounters
	Total  EnergyCounters
}

type EnergyCounters struct {
	Energy         float64
	EnergyReturned float64
}

type DevicePayload struct {
	Source      string          `json:"src"`
	Destination string          `json:"dst"`
	Method      string          `json:"method"`
	Params      json.RawMessage `json:"params"`
}

type UnixTimestamp time.Time

func (ts *UnixTimestamp) UnmarshalJSON(msg []byte) error {
	var val float64

	err := json.Unmarshal(msg, &val)
	if err != nil {
		return err
	}

	sec, dec := math.Modf(val)

	*ts = UnixTimestamp(time.Unix(int64(sec), int64(dec*1e9)))

	return nil
}

type NotifyStatusMessage struct {
	Timestamp UnixTimestamp        `json:"ts"`
	EM0       *NotifyStatusEm0     `json:"em:0"`
	EMData0   *NotifyStatusEMData0 `json:"emdata:0"`
}

type NotifyStatusEm0 struct {
	ID int32 `json:"id,omitempty"`

	AActualPower   float64 `json:"a_act_power,omitempty"`
	AApparentPower float64 `json:"a_aprt_power,omitempty"`
	ACurrent       float64 `json:"a_current,omitempty"`
	AFrequency     float64 `json:"a_freq,omitempty"`
	APowerFactor   float64 `json:"a_pf,omitempty"`
	AVoltage       float64 `json:"a_voltage,omitempty"`

	BActualPower   float64 `json:"b_act_power,omitempty"`
	BApparentPower float64 `json:"b_aprt_power,omitempty"`
	BCurrent       float64 `json:"b_current,omitempty"`
	BFrequency     float64 `json:"b_freq,omitempty"`
	BPowerFactor   float64 `json:"b_pf,omitempty"`
	BVoltage       float64 `json:"b_voltage,omitempty"`

	CActualPower   float64 `json:"c_act_power,omitempty"`
	CApparentPower float64 `json:"c_aprt_power,omitempty"`
	CCurrent       float64 `json:"c_current,omitempty"`
	CFrequency     float64 `json:"c_freq,omitempty"`
	CPowerFactor   float64 `json:"c_pf,omitempty"`
	CVoltage       float64 `json:"c_voltage,omitempty"`

	// nil if there is no N clamp
	NCurrent           *float64 `json:"n_current,omitempty"`
	TotalActualPower   float64  `json:"total_act_power,omitempty"`
	TotalApparentPower float64  `json:"total_aprt_power,omitempty"`
	TotalCurrent       float64  `json:"total_current,omitempty"`
}

type NotifyStatusEMData0 struct {
	ID                         int32   `json:"id,omitempty"`
	ATotalActualEnergy         float64 `json:"a_total_act_energy,omitempty"`
	ATotalActualReturnedEnergy float64 `json:"a_total_act_ret_energy,omitempty"`
	BTotalActualEnergy         float64 `json:"b_total_act_energy,omitempty"`
	BTotalActualReturnedEnergy float64 `json:"b_total_act_ret_energy,omitempty"`
	CTotalActualEnergy         float64 `json:"c_total_act_energy,omitempty"`
	CTotalActualReturnedEnergy float64 `json:"c_total_act_ret_energy,omitempty"`
	TotalActual                float64 `json:"total_act,omitempty"`
	TotalActualReturned        float64 `json:"total_act_ret,omitempty"`
}

type Repository interface {
	InsertMeasurements(ctx context.Context, measurement Measurement) error
	InsertCounters(ctx context.Context, counters Counters) error
}

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) HandleMessage(message *paho.Publish) {
	var devicePayload DevicePayload

	err := json.Unmarshal(message.Payload, &devicePayload)
	if err != nil {
		return
	}

	switch devicePayload.Method {
	case NotifyStatus:
		var nfm NotifyStatusMessage

		err := json.Unmarshal(devicePayload.Params, &nfm)
		if err != nil {
			return
		}

		if nfm.EM0 != nil {
			err = h.handleEnergyStats(devicePayload.Source, nfm)
			if err != nil {
				return
			}
		}

		if nfm.EMData0 != nil {
			err = h.handleEnergyCounters(devicePayload.Source, nfm)
			if err != nil {
				return
			}
		}

	default:
		// unknown message
	}
}

func (h *Handler) handleEnergyStats(device string, nfm NotifyStatusMessage) error {
	measurement := Measurement{
		Time:   time.Time(nfm.Timestamp),
		Device: device,
		PhaseA: PhaseMeasurement{
			ActualPower:   nfm.EM0.AActualPower,
			ApparentPower: nfm.EM0.AApparentPower,
			Current:       nfm.EM0.ACurrent,
			Frequency:     nfm.EM0.AFrequency,
			PowerFactor:   nfm.EM0.APowerFactor,
			Voltage:       nfm.EM0.AVoltage,
		},
		PhaseB: PhaseMeasurement{
			ActualPower:   nfm.EM0.BActualPower,
			ApparentPower: nfm.EM0.BApparentPower,
			Current:       nfm.EM0.BCurrent,
			Frequency:     nfm.EM0.BFrequency,
			PowerFactor:   nfm.EM0.BPowerFactor,
			Voltage:       nfm.EM0.BVoltage,
		},
		PhaseC: PhaseMeasurement{
			ActualPower:   nfm.EM0.CActualPower,
			ApparentPower: nfm.EM0.CApparentPower,
			Current:       nfm.EM0.CCurrent,
			Frequency:     nfm.EM0.CFrequency,
			PowerFactor:   nfm.EM0.CPowerFactor,
			Voltage:       nfm.EM0.CVoltage,
		},
		TotalActualPower:   nfm.EM0.TotalActualPower,
		TotalApparentPower: nfm.EM0.TotalApparentPower,
		TotalCurrent:       nfm.EM0.TotalCurrent,
	}

	return h.repo.InsertMeasurements(context.TODO(), measurement)
}

func (h *Handler) handleEnergyCounters(device string, nfm NotifyStatusMessage) error {
	counters := Counters{
		Time:   time.Time(nfm.Timestamp),
		Device: device,
		PhaseA: EnergyCounters{
			Energy:         nfm.EMData0.ATotalActualEnergy,
			EnergyReturned: nfm.EMData0.ATotalActualReturnedEnergy,
		},
		PhaseB: EnergyCounters{
			Energy:         nfm.EMData0.BTotalActualEnergy,
			EnergyReturned: nfm.EMData0.BTotalActualReturnedEnergy,
		},
		PhaseC: EnergyCounters{
			Energy:         nfm.EMData0.CTotalActualEnergy,
			EnergyReturned: nfm.EMData0.CTotalActualReturnedEnergy,
		},
		Total: EnergyCounters{
			Energy:         nfm.EMData0.TotalActual,
			EnergyReturned: nfm.EMData0.TotalActualReturned,
		},
	}

	return h.repo.InsertCounters(context.TODO(), counters)
}
