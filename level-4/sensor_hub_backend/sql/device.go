package sql

import "time"

type DeviceEntity struct {
	DeviceId        string    `db:"device_id" json:"device_id"`
	Name            string    `db:"name" json:"name"`
	LastIp          string    `db:"last_ip" json:"last_ip"`
	LastReading     string    `db:"last_reading" json:"last_reading"`
	LastReadingTime time.Time `db:"last_reading_time" json:"last_reading_time"`
	Authorized      bool      `db:"authorized" json:"authorized"`
}

var authorizedDevicesCache = make(map[string]bool)

func SelectDevices() ([]DeviceEntity, error) {
	devices := make([]DeviceEntity, 0)
	db := getDb()
	err := db.Select(&devices, "SELECT * FROM devices")
	return devices, err
}

func ToggleDeviceAuthorization(deviceId string) (bool, error) {
	db := getDb()
	_, err := db.Exec("UPDATE devices SET authorized = NOT authorized WHERE device_id = $1", deviceId)

	if err == nil {
		v, ok := authorizedDevicesCache[deviceId]
		if ok {
			authorizedDevicesCache[deviceId] = !v
			return !v, nil
		}
	}
	return false, err
}

func IsDeviceAuthorizedCached(deviceId string) (bool, error) {
	cacheHit, ok := authorizedDevicesCache[deviceId]
	if ok {
		return cacheHit, nil
	}

	db := getDb()
	var authorized bool
	err := db.Get(&authorized, "SELECT authorized FROM devices WHERE device_id = $1", deviceId)
	if err != nil {
		return false, err
	}

	authorizedDevicesCache[deviceId] = authorized
	return authorized, nil
}
