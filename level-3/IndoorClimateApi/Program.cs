using IndoorClimateApi;

var builder = WebApplication.CreateBuilder(args);

builder.Logging
       .AddConsole()
       .AddDebug();

var app = builder.Build();

DependencyProvider.Initialize(app.Services.GetRequiredService<ILoggerFactory>(), app.Configuration);

app.UseMiddleware<LoggingMiddleware>();
app.UseMiddleware<ExceptionHandlingMiddleware>();

app.MapPost("/scd30", Scd30Endpoint.HandleScd30Data);

app.Run();
