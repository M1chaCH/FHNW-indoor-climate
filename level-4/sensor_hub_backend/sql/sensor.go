package sql

import (
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/logs"
	"sync"
	"time"
)

type lockableDevices struct {
	mu      sync.Mutex
	devices []*DeviceEntity
}

var entitiesToUpdate *lockableDevices

func UpsertSensorReadingThrottled(entityToUpdate *DeviceEntity) {
	defer ensureUpdateLoopRunning()
	defer func() {
		if entitiesToUpdate != nil {
			entitiesToUpdate.mu.Unlock()
		}
	}()

	if entitiesToUpdate == nil {
		entitiesToUpdate = &lockableDevices{}
		entitiesToUpdate.mu.Lock()

		devices := make([]*DeviceEntity, 1)
		devices[0] = entityToUpdate

		entitiesToUpdate.devices = devices
		return
	}

	entitiesToUpdate.mu.Lock()

	for i, queuedEntity := range entitiesToUpdate.devices {
		if queuedEntity.DeviceId == entityToUpdate.DeviceId {
			entitiesToUpdate.devices[i] = entityToUpdate
			return
		}
	}

	entitiesToUpdate.devices = append(entitiesToUpdate.devices, entityToUpdate)
}

var ticker *time.Ticker

func ensureUpdateLoopRunning() {
	if ticker == nil {
		ticker = time.NewTicker(30 * time.Second)
		go throttledUpsertLoop()
	}
}

func throttledUpsertLoop() {
	logs.LogInfo("Starting throttled upsert loop for sensor devices")
	ctx := lifecycle.GetStopContext()

	for {
		select {
		case <-ticker.C:
			upsertQueuedDevices()
			continue
		case <-ctx.Done():
		}
		break
	}

	ticker.Stop()
	logs.LogInfo("Stopped throttled upsert loop for sensor devices")
}

func upsertQueuedDevices() {
	if entitiesToUpdate == nil {
		return
	}

	entitiesToUpdate.mu.Lock()
	if len(entitiesToUpdate.devices) == 0 {
		entitiesToUpdate.mu.Unlock()
		return
	}

	devices := make([]*DeviceEntity, len(entitiesToUpdate.devices))
	copy(devices, entitiesToUpdate.devices)
	entitiesToUpdate.devices = make([]*DeviceEntity, 0)
	entitiesToUpdate.mu.Unlock()

	tx, err := getDb().Beginx()
	if err != nil {
		logs.LogErr("Failed to start transaction", err)
		return
	}

	for _, device := range devices {
		statement, err := tx.PrepareNamed(`
		INSERT INTO devices(device_id, name, last_ip, last_reading, last_reading_time)
		VALUES (:device_id, :name, :last_ip, :last_reading, :last_reading_time)
		ON CONFLICT (device_id) 
		    DO UPDATE SET last_ip = :last_ip, 
		                  last_reading = :last_reading, 
		                  last_reading_time = :last_reading_time,
		                  name = :name
		RETURNING (XMAX = 0) as inserted
		`)

		if err != nil {
			logs.LogErrCustom("Failed to prepare upsert statement for device %s: %s\n", device.DeviceId, err)
			continue
		}

		var inserted bool
		if err := statement.Get(&inserted, device); err != nil {
			logs.LogErrCustom("Failed to upsert device %s: %s\n", device.DeviceId, err)
			continue
		}

		if inserted {
			deviceChangedObs.Emit(1)
		}

		statement.Close()
	}

	err = tx.Commit()
	if err != nil {
		logs.LogErr("Failed to commit transaction", err)
		// TODO, do I want to add the devices back into the queue?
	} else {
		logs.LogInfo("successfully upserted %d devices\n", len(devices))
	}
}
