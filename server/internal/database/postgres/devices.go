package postgres

import (
	"encoding/json"
	"remoteesp/internal/model"
	"time"
)

func (r *Repository) DevicesUploaded(devices []model.Device) error {
	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Preparex(`
		INSERT INTO devices (
			ssdp,
			mdns,
			active,
			activated,
			started,
			updated
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT(ssdp) DO UPDATE SET
			mdns      = excluded.mdns,
			active    = excluded.active,
			activated = excluded.activated,
			started   = excluded.started,
			updated   = excluded.updated
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range devices {
		_, err := stmt.Exec(
			d.SSDP,
			d.MDNS,
			d.Active,
			d.Activated,
			d.Started,
			d.Updated,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *Repository) Devices(filter string) ([]model.Device, error) {
	rows, err := r.db.Queryx(`
		SELECT
			ssdp,
			mdns,
			active,
			activated,
			started,
			updated,
			pins,
			action
		FROM devices
		WHERE ssdp LIKE $1
		   OR mdns LIKE $2
		ORDER BY id DESC
	`, filter+"%", filter+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []model.Device

	for rows.Next() {
		var d model.Device
		var jsonStr string

		err := rows.Scan(
			&d.SSDP,
			&d.MDNS,
			&d.Active,
			&d.Activated,
			&d.Started,
			&d.Updated,
			&jsonStr,
			&d.Action,
		)
		if err != nil {
			return nil, err
		}

		if jsonStr != "" {
			var restoredMap map[string]any
			err = json.Unmarshal([]byte(jsonStr), &restoredMap)
			if err != nil {
				return nil, err
			}
			d.Pins = restoredMap
		}

		devices = append(devices, d)
	}

	return devices, rows.Err()
}

func (r *Repository) UpdateDevice(device model.Device) error {
	now := time.Now().UTC().Unix()

	jsonData, err := json.Marshal(device.Pins)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`
		INSERT INTO devices (
			ssdp,
			mdns,
			active,
			activated,
			started,
			updated,
			pins,
			action
		)
		VALUES ($1, $2, TRUE, $3, $4, $5, $6, $7)
		ON CONFLICT(ssdp) DO UPDATE SET
			mdns = excluded.mdns,
			active = TRUE,
			started = excluded.started,
			updated = excluded.updated,
			pins   = excluded.pins,
			action = excluded.action,
			activated = CASE
				WHEN devices.activated IS NULL OR devices.activated = 0
					THEN excluded.activated
				ELSE devices.activated
			END
	`,
		device.SSDP,
		device.MDNS,
		now,
		device.Started,
		now,
		string(jsonData),
		device.Action,
	)
	return err
}

func (r *Repository) UpdateDeviceAction(ssdp string, action string) error {
	now := time.Now().UTC().Unix()

	_, err := r.db.Exec(`
		UPDATE devices SET
			action = $1,
			updated = $2
		WHERE ssdp = $3;
	`,
		action,
		now,
		ssdp,
	)
	return err
}

func (r *Repository) DeleteDevice(ssdp string) error {
	_, err := r.db.Exec(`
		DELETE FROM devices
		WHERE ssdp = $1;`,
		ssdp,
	)
	return err
}
