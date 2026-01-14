FROM mcr.microsoft.com/dotnet/aspnet:9.0 AS base
USER $APP_UID
ARG API_PORT=8080
WORKDIR /app
EXPOSE $API_PORT

FROM mcr.microsoft.com/dotnet/sdk:9.0 AS build
ARG BUILD_CONFIGURATION=Release
WORKDIR /src
COPY ["IndoorClimateApi.csproj", "./"]
RUN dotnet restore "IndoorClimateApi.csproj"
COPY . .
WORKDIR "/src/"
RUN dotnet build "./IndoorClimateApi.csproj" -c $BUILD_CONFIGURATION -o /app/build

FROM build AS publish
ARG BUILD_CONFIGURATION=Release
RUN dotnet publish "./IndoorClimateApi.csproj" -c $BUILD_CONFIGURATION -o /app/publish /p:UseAppHost=false

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
ENTRYPOINT ["dotnet", "IndoorClimateApi.dll"]
