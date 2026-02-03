package elastic

import (
	"bytes"
	"encoding/json"
	"sensor_hub_backend/lifecycle"
	"sensor_hub_backend/logs"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/elastic/go-elasticsearch/v9/esutil"
)

type SensorDataDocument struct {
	DeviceId   string
	SensorType string
	Timestamp  time.Time
	DeviceName string
	Values     map[string]interface{}
}

type lockableTypedDocs struct {
	mu   sync.Mutex
	docs []*SensorDataDocument
}

var sensorDataToSend *lockableTypedDocs

func SendSensorDataToElasticDebounced(typedDocs []*SensorDataDocument) {
	defer ensureUpdateLoopRunning()
	defer func() {
		if sensorDataToSend != nil {
			sensorDataToSend.mu.Unlock()
		}
	}()

	if sensorDataToSend == nil {
		sensorDataToSend = &lockableTypedDocs{}
		sensorDataToSend.mu.Lock()

		sensorDataToSend.docs = typedDocs
		return
	}

	sensorDataToSend.mu.Lock()
	sensorDataToSend.docs = append(sensorDataToSend.docs, typedDocs...)
}

func SendSensorDataToElastic(typedDocs []*SensorDataDocument) {
	if len(typedDocs) > 10 {
		sendInBulk(typedDocs)
	}

	for _, typedDoc := range typedDocs {
		if typedDoc == nil {
			continue
		}

		body, err := createJsonBytes(typedDoc)
		if err != nil {
			logs.LogErr("Failed to marshal sensor data document", err)
			continue
		}

		req := esapi.IndexRequest{
			Index: "ipro-sensor-hub-data",
			Body:  bytes.NewReader(body),
		}

		res, err := req.Do(lifecycle.GetStopContext(), es)
		if err != nil {
			logs.LogErr("Failed to send sensor data document to elastic", err)
			continue
		}

		err = res.Body.Close()
		if err != nil {
			logs.LogErr("Failed to close response body", err)
			continue
		}
	}
}

func sendInBulk(typedDocs []*SensorDataDocument) {
	stopContext := lifecycle.GetStopContext()

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         "ipro-sensor-hub-data",
		Client:        es,
		NumWorkers:    8,
		FlushBytes:    10 * 1024 * 1024,
		FlushInterval: 10 * time.Second,
	})

	if err != nil {
		logs.LogErr("Failed to create bulk indexer", err)
		return
	}

	for _, typedDoc := range typedDocs {
		doc, err := createJsonBytes(typedDoc)
		if err != nil {
			logs.LogErr("Failed to marshal sensor data document", err)
			continue
		}

		err = bi.Add(stopContext, esutil.BulkIndexerItem{
			Action: "create",
			Body:   bytes.NewReader(doc),
		})

		if err != nil {
			logs.LogErr("Failed to add sensor data document to bulk indexer", err)
			continue
		}
	}

	if err := bi.Close(stopContext); err != nil {
		logs.LogErr("Failed to close bulk indexer", err)
	}

	stats := bi.Stats()
	logs.LogInfo("Indexed %d sensor data documents successfully (%d failed)\n", stats.NumFlushed, stats.NumFailed)
}

func createJsonBytes(typedDoc *SensorDataDocument) ([]byte, error) {
	doc := map[string]interface{}{
		"@timestamp":   time.Now().UTC().Format(time.RFC3339),
		"device_id":    typedDoc.DeviceId,
		"sensor_type":  typedDoc.SensorType,
		"device_name":  typedDoc.DeviceName,
		"reading_time": typedDoc.Timestamp.Format(time.RFC3339),
	}

	for k, v := range typedDoc.Values {
		doc[k] = v
	}

	return json.Marshal(doc)
}

var ticker *time.Ticker

func ensureUpdateLoopRunning() {
	if ticker == nil {
		ticker = time.NewTicker(30 * time.Second)
		go periodicSendSensorDataToElastic()
	}
}

func periodicSendSensorDataToElastic() {
	logs.LogInfo("ElasticSensorDataPublisher: Starting loop.")
	ctx := lifecycle.GetStopContext()

	for {
		select {
		case <-ticker.C:
			if sensorDataToSend == nil {
				logs.LogInfo("ElasticSensorDataPublisher: sensorDataToSend nil -> skipping")
				continue
			}

			sensorDataToSend.mu.Lock()
			if len(sensorDataToSend.docs) == 0 {
				sensorDataToSend.mu.Unlock()
				logs.LogInfo("ElasticSensorDataPublisher: no changes -> skipping")
				continue
			}

			docsToSend := make([]*SensorDataDocument, len(sensorDataToSend.docs))
			copy(docsToSend, sensorDataToSend.docs)
			sensorDataToSend.docs = make([]*SensorDataDocument, 0)
			sensorDataToSend.mu.Unlock()

			sendInBulk(docsToSend)

			continue
		case <-ctx.Done():
			logs.LogInfo("ElasticSensorDataPublisher: Stopping loop.")
		}
		break
	}

	logs.LogInfo("ElasticSensorDataPublisher: Stopped.")
}
