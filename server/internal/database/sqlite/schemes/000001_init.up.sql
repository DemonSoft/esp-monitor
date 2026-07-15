CREATE TABLE devices
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    ssdp        TEXT NOT NULL UNIQUE,
    mdns        TEXT NOT NULL DEFAULT "",
    active      INTEGER NOT NULL DEFAULT 0,
    activated   INTEGER,
    started     INTEGER,
    updated     INTEGER,
    pins        TEXT NOT NULL DEFAULT "",
    action      TEXT NOT NULL DEFAULT ""
);

CREATE INDEX idx_devices_ssdp ON devices(ssdp);
CREATE INDEX idx_devices_mdns ON devices(mdns);
