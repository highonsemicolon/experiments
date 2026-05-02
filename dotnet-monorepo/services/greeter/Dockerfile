FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /repo

COPY Directory.Build.props Directory.Packages.props global.json ./
COPY apis/ apis/
COPY packages/platform/Platform.csproj           packages/platform/
COPY services/greeter/Greeter.Service.csproj     services/greeter/

RUN dotnet restore services/greeter/Greeter.Service.csproj

COPY packages/platform/  packages/platform/
COPY services/greeter/   services/greeter/

RUN dotnet publish services/greeter/Greeter.Service.csproj \
    -c Release -o /app

FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=build /app .

ENTRYPOINT ["dotnet", "Greeter.Service.dll"]