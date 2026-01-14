namespace IndoorClimateApi;

/// <summary>
/// Dependency provider for the application.
/// This is probably a bit of an anti-pattern, but for a service this small I did not want to set up full dependency injection.
/// </summary>
public class DependencyProvider(ILoggerFactory? _loggerFactory, IConfiguration? _configuration)
{
    public static DependencyProvider Instance { get; private set; } = null!;
    public IConfiguration Configuration => _configuration ?? throw new InvalidOperationException("Configuration not initialized");
    public ILogger<T> GetLogger<T>() => _loggerFactory?.CreateLogger<T>() ?? throw new InvalidOperationException("Logger not initialized");
    public ElasticService ElasticService { get; private set; } = null!;

    public static void Initialize(ILoggerFactory? loggerFactory, IConfiguration? configuration)
    {
        Instance = new DependencyProvider(loggerFactory, configuration);
        Instance.ElasticService = new ElasticService(Instance.GetLogger<ElasticService>());
    }
}   

public static class ConfigurationExtensions
{
    public static string GetApiKey (this IConfiguration configuration) => configuration["SensorApiKey"] ?? throw new InvalidOperationException("Sensor API key not found");
    
    public static string GetElasticUrl (this IConfiguration configuration) => configuration.GetSection("Elastic")["Url"] ?? throw new InvalidOperationException("Elastic URL not found");
    public static string GetElasticApiKey (this IConfiguration configuration) => configuration.GetSection("Elastic")["ApiKey"] ?? throw new InvalidOperationException("Elastic URL not found");
}