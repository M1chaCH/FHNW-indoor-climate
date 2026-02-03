package sql

import (
	"sensor_hub_backend/logs"
)

type DeviceConfigEntity struct {
	DeviceId   string `db:"device_id"`
	ConfigJson string `db:"config_json"`
}

func SelectDeviceConfigJson(deviceId string) (string, error) {
	var configJson string
	db := getDb()
	err := db.Get(&configJson, "SELECT config_json FROM device_configs WHERE device_id = $1", deviceId)
	return configJson, err
}

func UpsertConfigJson(entity *DeviceConfigEntity) {
	db := getDb()
	_, err := db.Exec(`
		INSERT INTO device_configs (device_id, config_json) 
		VALUES ($1, $2) 
		ON CONFLICT (device_id) 
		DO UPDATE SET config_json = $2`,
		entity.DeviceId,
		entity.ConfigJson)

	if err != nil {
		logs.LogErr("Failed to upsert device config", err)
	} else {
		logs.LogInfo("Successfully upserted device config: %s\n", entity.DeviceId)
	}
}
