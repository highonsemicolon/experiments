FROM mcr.microsoft.com/dotnet/sdk:8.0 AS build
WORKDIR /repo

COPY Directory.Build.props Directory.Packages.props global.json ./
COPY apis/ apis/
COPY packages/platform/Platform.csproj            packages/platform/
COPY services/greetWorker/Greet.Worker.csproj     services/greetWorker/

RUN dotnet restore services/greetWorker/Greet.Worker.csproj

COPY packages/platform/    packages/platform/
COPY services/greetWorker/ services/greetWorker/

RUN dotnet publish services/greetWorker/Greet.Worker.csproj \
    -c Release -o /app

FROM mcr.microsoft.com/dotnet/aspnet:8.0
WORKDIR /app
COPY --from=build /app .

ENTRYPOINT ["dotnet", "Greet.Worker.dll"]