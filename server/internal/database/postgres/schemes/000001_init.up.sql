CREATE TABLE devices
(
    id          SERIAL PRIMARY KEY,
    ssdp        TEXT NOT NULL UNIQUE,
    mdns        TEXT NOT NULL DEFAULT '',
    active      BOOLEAN NOT NULL DEFAULT FALSE,
    activated   BIGINT,
    started     BIGINT,
    updated     BIGINT,
    pins        TEXT NOT NULL DEFAULT '',
    action      TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_devices_mdns ON devices(mdns);
