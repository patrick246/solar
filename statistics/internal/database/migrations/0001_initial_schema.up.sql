CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE home_energy_counters
(
    id                            bigint generated always as identity,
    time                          timestamp with time zone default now(),
    device                        text not null,
    phase_a_total_energy          float8,
    phase_a_total_energy_returned float8,
    phase_b_total_energy          float8,
    phase_b_total_energy_returned float8,
    phase_c_total_energy          float8,
    phase_c_total_energy_returned float8,
    total_energy                  float8,
    total_energy_returned         float8,
    PRIMARY KEY (id, time)
);

CREATE TABLE home_energy_stats
(
    id                     bigint generated always as identity,
    time                   timestamp with time zone default now(),
    device                 text not null,
    phase_a_actual_power   float8,
    phase_a_apparent_power float8,
    phase_a_current        float8,
    phase_a_frequency      float8,
    phase_a_power_factor   float8,
    phase_a_voltage        float8,
    phase_b_actual_power   float8,
    phase_b_apparent_power float8,
    phase_b_current        float8,
    phase_b_frequency      float8,
    phase_b_power_factor   float8,
    phase_b_voltage        float8,
    phase_c_actual_power   float8,
    phase_c_apparent_power float8,
    phase_c_current        float8,
    phase_c_frequency      float8,
    phase_c_power_factor   float8,
    phase_c_voltage        float8,
    total_actual_power     float8,
    total_apparent_power   float8,
    total_current          float8,
    PRIMARY KEY (id, time)
);

SELECT create_hypertable('home_energy_stats', 'time');
SELECT create_hypertable('home_energy_counters', 'time');