package shelly3em

import (
	"context"
	"database/sql"
)

type TimescaleRepository struct {
	db *sql.DB
}

func NewTimescaleRepository(db *sql.DB) *TimescaleRepository {
	return &TimescaleRepository{
		db: db,
	}
}

func (t *TimescaleRepository) InsertMeasurements(ctx context.Context, measurement Measurement) error {
	_, err := t.db.ExecContext(ctx, `
		INSERT INTO home_energy_stats (
			time,
			device,
			phase_a_actual_power,
			phase_a_apparent_power,
			phase_a_current,
			phase_a_frequency,
			phase_a_power_factor,
			phase_a_voltage,
			phase_b_actual_power,
			phase_b_apparent_power,
			phase_b_current,
			phase_b_frequency,
			phase_b_power_factor,
			phase_b_voltage,
			phase_c_actual_power,
			phase_c_apparent_power,
			phase_c_current,
			phase_c_frequency,
			phase_c_power_factor,
			phase_c_voltage,
			total_actual_power,
			total_apparent_power,
			total_current
        ) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23        
		)`,
		measurement.Time,
		measurement.Device,
		measurement.PhaseA.ActualPower,
		measurement.PhaseA.ApparentPower,
		measurement.PhaseA.Current,
		measurement.PhaseA.Frequency,
		measurement.PhaseA.PowerFactor,
		measurement.PhaseA.Voltage,
		measurement.PhaseB.ActualPower,
		measurement.PhaseB.ApparentPower,
		measurement.PhaseB.Current,
		measurement.PhaseB.Frequency,
		measurement.PhaseB.PowerFactor,
		measurement.PhaseB.Voltage,
		measurement.PhaseC.ActualPower,
		measurement.PhaseC.ApparentPower,
		measurement.PhaseC.Current,
		measurement.PhaseC.Frequency,
		measurement.PhaseC.PowerFactor,
		measurement.PhaseC.Voltage,
		measurement.TotalActualPower,
		measurement.TotalApparentPower,
		measurement.TotalCurrent,
	)

	return err
}

func (t *TimescaleRepository) InsertCounters(ctx context.Context, counters Counters) error {
	_, err := t.db.ExecContext(ctx, `
		INSERT INTO home_energy_counters (
			time,
			device,
			phase_a_total_energy,
		    phase_a_total_energy_returned,
			phase_b_total_energy,
		    phase_b_total_energy_returned,
			phase_c_total_energy,
		    phase_c_total_energy_returned,
			total_energy,
		    total_energy_returned
        ) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)`,
		counters.Time,
		counters.Device,
		counters.PhaseA.Energy,
		counters.PhaseA.EnergyReturned,
		counters.PhaseB.Energy,
		counters.PhaseB.EnergyReturned,
		counters.PhaseC.Energy,
		counters.PhaseC.EnergyReturned,
		counters.Total.Energy,
		counters.Total.EnergyReturned,
	)

	return err
}
