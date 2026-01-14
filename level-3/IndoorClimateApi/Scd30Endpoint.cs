using Microsoft.AspNetCore.Mvc;

namespace IndoorClimateApi;

public record Scd30Dto(double co2, double temp, double hum, string device, DateTime uptime)
{
    public string? SourceIp { get; set; }
}

public static class Scd30Endpoint
{
    public static async Task HandleScd30Data([FromBody] Scd30Dto dto,
                                             [FromHeader(Name = "X-Real-IP")] string? sourceIp,
                                             [FromHeader(Name = "X-Api-Key")] string? apiKey)
    {
        var apiKeyConfig = DependencyProvider.Instance.Configuration.GetApiKey();
        if (apiKey is null || apiKey != apiKeyConfig)
        {
            throw new UnauthorizedAccessException();
        }

        dto.SourceIp = sourceIp;
        await DependencyProvider.Instance.ElasticService.SendScd30Data(dto);
    }
}
