# API Documentation

This directory contains auto-generated Swagger/OpenAPI documentation.

## Generate Documentation

To generate the Swagger documentation, run:

```bash
# Install swag CLI if not already installed
go install github.com/swaggo/swag/cmd/swag@v1.16.6

# Generate docs
swag init -g cmd/api/main.go -o docs
```

Or use the Makefile:

```bash
make swagger
```

## Access Documentation

Once the API is running, access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## Files

After generation, this directory will contain:

- `docs.go` - Go documentation file
- `swagger.json` - OpenAPI JSON specification
- `swagger.yaml` - OpenAPI YAML specification

These files are auto-generated and should not be edited manually.
