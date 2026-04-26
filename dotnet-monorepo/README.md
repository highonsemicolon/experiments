# Platform

## Commands

make restore
make build
make test
make run

## Structure

services/   -> deployable apps
packages/   -> shared libraries
apis/       -> contracts
tests/      -> tests
deploy/     -> infra
docs/       -> docs


# Minimal templates

dotnet new classlib -n Logging -o packages/logging
dotnet sln add packages/logging/Logging.csproj

dotnet new web -n Greeter.Service -o services/greeter

dotnet new xunit -n Greeter.Tests -o tests/greeter


docker compose --profile infra up