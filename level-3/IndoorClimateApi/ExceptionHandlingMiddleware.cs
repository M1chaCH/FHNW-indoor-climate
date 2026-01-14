using System.Text.Json;

namespace IndoorClimateApi;

public class ExceptionHandlingMiddleware
{
    private readonly RequestDelegate _next;
    private readonly ILogger<ExceptionHandlingMiddleware> _logger;

    public ExceptionHandlingMiddleware(RequestDelegate next)
    {
        _next = next;
        _logger = DependencyProvider.Instance.GetLogger<ExceptionHandlingMiddleware>();
    }
    
    public async Task InvokeAsync(HttpContext context)
    {
        try
        {
            await _next(context);
        }
        catch (UnauthorizedAccessException)
        {
            _logger.LogWarning("Unauthorized access");
            context.Response.StatusCode = 401;
            context.Response.ContentType = "text/plain";
            await context.Response.WriteAsync("Unauthorized");
        }
        catch (ElasticsearchException)
        {
            context.Response.StatusCode = 500;
            context.Response.ContentType = "text/plain";
            await context.Response.WriteAsync("An internal error with elasticsearch occured");
        }
        catch (JsonException jsonException)
        {
            _logger.LogError(jsonException, "Invalid JSON");
            context.Response.StatusCode = 400;
            context.Response.ContentType = "text/plain";
            await context.Response.WriteAsync("Invalid Body: " + jsonException.Message);
        }
        catch (Exception ex)
        {
            _logger.LogError(ex, "Unhandled exception");
            context.Response.StatusCode = 500;
            context.Response.ContentType = "text/plain";
            await context.Response.WriteAsync("Internal Server Error");
        }
    }
}