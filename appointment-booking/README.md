simply run it using `docker compose up`


OpenAPI Code gen:
```bash
oapi-codegen --config=api/oapi-server.cfg.yaml api/swagger.yaml
```

This repo has reused few already existing components: config, logger. Same can be found in respective commit history.

For gateway config nginx, istio, etc can be used to validate OAuth2 tokens.

nginx reference (self compiled with lua support): https://github.com/highonsemicolon/aura/tree/dev/nginx
