using System.Text.Json.Serialization;
using Elastic.Clients.Elasticsearch;
using Elastic.Transport;

namespace IndoorClimateApi;

public class ElasticService
{
    private readonly Lazy<ElasticsearchClient> _client;
    private readonly ILogger<ElasticService> _logger;

    public ElasticService(ILogger<ElasticService> logger)
    {
        _logger = logger;
        _client = new Lazy<ElasticsearchClient>(CreateClient);
    }
    
    public async Task SendScd30Data(Scd30Dto data)
    {
        _logger.LogDebug("sending scd30 data: {data}", data);
        var response = await _client.Value.IndexAsync(new ElasticScd30Dto
                                                      {
                                                          co2 = data.co2,
                                                          temperature = data.temp,
                                                          humidity = data.hum,
                                                          device_id = data.device,
                                                          sensor_type = "SCD30",
                                                          sensor_uptime = data.uptime,
                                                          sensor_ip = data.SourceIp,
                                                          timestamp = DateTime.UtcNow
                                                      },
                                                      idx => idx.Index("ipro-sensor-data-stream"));
        if (!response.IsValidResponse)
        {
            throw new ElasticsearchException($"Failed to send data to elastic: {response.DebugInformation}");
        }
    }

    private ElasticsearchClient CreateClient()
    {
        var configuration = DependencyProvider.Instance.Configuration;
        var settings = new ElasticsearchClientSettings(new Uri(configuration.GetElasticUrl()))
            .Authentication(new ApiKey(configuration.GetElasticApiKey()))
            .OnRequestCompleted(apiCallDetails =>
                                {
                                    if (apiCallDetails.HasSuccessfulStatusCode)
                                    {
                                        _logger.LogDebug("Call to elastic succeeded: {Method} {Uri}", apiCallDetails.HttpMethod, apiCallDetails.Uri);
                                    }
                                    else
                                    {
                                        _logger.LogWarning("Call to elastic failed: {DebugInformation}", apiCallDetails.DebugInformation);
                                    }
                                });
        
        return new ElasticsearchClient(settings);
    }

    // ReSharper disable InconsistentNaming
    private class ElasticScd30Dto
    {
        public double co2 { get; init; }
        public double temperature { get; init; }
        public double humidity { get; init; }
        public string device_id { get; init; } = string.Empty;
        public string sensor_type { get; init; } = string.Empty;
        public DateTime sensor_uptime { get; init; }
        public string sensor_ip { get; init; } = string.Empty;
        [JsonPropertyName("@timestamp")]
        public DateTime timestamp { get; init; }
    }
}

public class ElasticsearchException(string _message) : Exception(_message);