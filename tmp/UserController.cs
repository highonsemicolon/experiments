using Microsoft.AspNetCore.Mvc;

namespace UserService.Controllers
{
    [ApiController]
    [Route("api/[controller]")]
    public class UsersController : ControllerBase
    {
        [HttpGet("{id}")]
        public IActionResult GetUser(int id)
        {
            // Simulate business logic
            return Ok(new { Id = id, Name = "Alice" });
        }

        [HttpPost]
        public IActionResult CreateUser([FromBody] UserDto user)
        {
            return CreatedAtAction(nameof(GetUser), new { id = 123 }, user);
        }

        [HttpDelete("{id}")]
        public IActionResult DeleteUser(int id)
        {
            return NoContent();
        }
    }

    public class UserDto
    {
        public string Name { get; set; }
    }
}
