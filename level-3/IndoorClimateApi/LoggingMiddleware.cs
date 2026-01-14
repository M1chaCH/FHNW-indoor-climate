using System.Diagnostics;

namespace IndoorClimateApi;

public class LoggingMiddleware
{
    private readonly RequestDelegate _next;
    private readonly ILogger<LoggingMiddleware> _logger;

    public LoggingMiddleware(RequestDelegate next)
    {
        _next = next;
        _logger = DependencyProvider.Instance.GetLogger<LoggingMiddleware>();
    }

    public async Task InvokeAsync(HttpContext context)
    {
        var watch = Stopwatch.StartNew();
        _logger.LogInformation("{Now} {Method} {Path}", DateTime.Now, context.Request.Method, context.Request.Path);
        await _next(context);
        _logger.LogInformation("{Now} {Method} {Path}: {StatusCode} ({Duration}ms)", DateTime.Now, context.Request.Method, context.Request.Path, context.Response.StatusCode, watch.Elapsed.TotalMilliseconds);
    }
}