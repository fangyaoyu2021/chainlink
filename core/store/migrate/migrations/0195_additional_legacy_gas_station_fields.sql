-- +goose Up
ALTER TABLE legacy_gasless_txs ADD tx_hash BYTEA;
ALTER TABLE legacy_gasless_txs ADD CONSTRAINT tx_hash_len_chk CHECK (
    octet_length(tx_hash) = 32
);

ALTER TABLE legacy_gas_station_sidecar_specs ADD status_update_url text NOT NULL;
ALTER TABLE legacy_gas_station_sidecar_specs ADD client_certificate text;
ALTER TABLE legacy_gas_station_sidecar_specs ADD client_key text;

-- +goose Down
ALTER TABLE legacy_gasless_txs DROP CONSTRAINT tx_hash_len_chk;
ALTER TABLE legacy_gasless_txs DROP tx_hash;

ALTER TABLE legacy_gas_station_sidecar_specs DROP status_update_url;
ALTER TABLE legacy_gas_station_sidecar_specs DROP client_certificate;
ALTER TABLE legacy_gas_station_sidecar_specs DROP client_key;