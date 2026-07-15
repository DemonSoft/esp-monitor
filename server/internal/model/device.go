package model

import (
	"encoding/csv"
	"fmt"
	"html/template"
	"mime/multipart"
	"remoteesp/internal/domain/utils"
	"sort"
	"strings"
	"time"
)

const timeout = 60 // sec.

type Device struct {
	SSDP      string
	MDNS      string
	Active    bool
	Activated int64
	Started   int64
	Updated   int64
	Pins      map[string]any
	Action    string
}

func (d Device) ActivatedHuman() string {
	if d.Activated == 0 {
		return ""
	}
	return time.Unix(d.Activated, 0).UTC().Format("2006-01-02 15:04:05")
}

func (d Device) StartedHuman() string {
	if d.Started == 0 {
		return ""
	}
	return time.Unix(d.Started, 0).UTC().Format("2006-01-02 15:04:05")
}

func (d Device) UpdatedHuman() string {
	if d.Updated == 0 {
		return ""
	}
	return time.Unix(d.Updated, 0).UTC().Format("2006-01-02 15:04:05")
}

func (d Device) UptimeHuman() string {
	if d.Started == 0 {
		return ""
	}
	totalSeconds := time.Now().UTC().Unix() - d.Started

	const (
		SecondsInDay   = 24 * 3600
		SecondsInMonth = 30 * SecondsInDay  // about 30 days
		SecondsInYear  = 365 * SecondsInDay // about 1 year
	)

	years := totalSeconds / SecondsInYear
	totalSeconds %= SecondsInYear

	months := totalSeconds / SecondsInMonth
	totalSeconds %= SecondsInMonth

	days := totalSeconds / SecondsInDay
	totalSeconds %= SecondsInDay

	hours := totalSeconds / 3600
	totalSeconds %= 3600

	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	// 4. Динамически собираем строку
	var result string

	if years > 0 {
		result += fmt.Sprintf("%d y. ", years)
	}
	if months > 0 {
		result += fmt.Sprintf("%d m. ", months)
	}
	if days > 0 {
		result += fmt.Sprintf("%d d. ", days)
	}

	result += fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)

	return result
}

func (d Device) PrettyAction() template.HTML {
	return utils.PrettyColoredJson(d.Action)
}

func ParseDevices(file multipart.File) ([]Device, error) {

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // all fields
	allRecords, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	list := make([]Device, 0, len(allRecords))
	for _, r := range allRecords {
		if len(r) < 2 {
			continue
		}

		device := Device{SSDP: r[0], MDNS: r[1], Updated: time.Now().UTC().Unix()}
		list = append(list, device)
	}

	return list, nil
}

func (d Device) PinsString() string {
	result := ""
	delimiter := " "

	keys := make([]string, 0, len(d.Pins))
	for k := range d.Pins {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		val := fmt.Sprintf("%v", d.Pins[k])
		pin := k + ":" + val + delimiter
		result = result + pin
	}
	return strings.TrimSuffix(result, delimiter)
}

func (d Device) TimeOut() bool {
	now := time.Now().UTC().Unix()
	updated := time.Unix(d.Updated, 0).UTC().Unix()
	limit := updated + timeout
	return int(limit) < int(now)
}
