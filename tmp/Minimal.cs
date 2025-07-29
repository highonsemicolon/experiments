var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/users1/{id}", (int id) => 
{
    return Results.Ok(new { Id = id, Name = "Alice" });
});

app.MapPost("/users", (UserDto user) => 
{
    return Results.Created($"/users/{123}", user);
});

app.MapDelete("/users/{id}", (int id) =>
{
    return Results.NoContent();
});

app.Run();

record UserDto(string Name);
