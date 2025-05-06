# Blog Service

This repository contains Protocol Buffer definitions for a Blog service with gRPC and REST API support.

## Overview

The service provides the following functionality:
- Create, read, update, and delete blogs
- Add comments to blogs
- List blogs with pagination

## Protocol Buffers

The service is defined using Protocol Buffers (proto3) and includes:
- Message definitions for blogs and comments
- A UUID message type for unique identifiers
- Service definitions with both gRPC and REST endpoints
- Field validations using buf validate

## Dependencies

- [buf](https://buf.build/) - For managing Protocol Buffer dependencies and generation
- [gRPC](https://grpc.io/) - For RPC communication
- [gRPC-Gateway](https://github.com/grpc-ecosystem/grpc-gateway) - For REST API generation
- [protovalidate](https://github.com/bufbuild/protovalidate-go) - For field validation

## Getting Started

### Prerequisites

1. Install Go (version 1.24 or later)
2. Install buf CLI:
   ```
   go install github.com/bufbuild/buf/cmd/buf@latest
   ```
3. Install required protoc plugins:
   ```
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
   go install github.com/bufbuild/protovalidate-go/cmd/protoc-gen-validate-go@latest
   ```

### Generating Code

To generate Go code from the Protocol Buffer definitions:

```
buf generate
```

This will generate:
- Go structs for all messages
- gRPC server and client code
- gRPC-Gateway REST API code
- Validation code
- OpenAPI v2 (Swagger) documentation in the `docs` directory

## API Endpoints

### gRPC

The gRPC service is defined in `protos/blog.proto` and includes the following methods:
- `CreateBlog`
- `GetBlog`
- `UpdateBlog`
- `DeleteBlog`
- `ListBlogs`
- `AddComment`

### REST

The REST API is generated from the gRPC service using gRPC-Gateway annotations:

| HTTP Method | Endpoint                      | Description                |
|-------------|-------------------------------|----------------------------|
| POST        | /v1/posts                     | Create a new blog          |
| GET         | /v1/posts/{id}                | Get a blog by ID           |
| PATCH       | /v1/posts/{id}                | Update a blog              |
| DELETE      | /v1/posts/{id}                | Delete a blog              |
| GET         | /v1/posts                     | List blogs                 |
| POST        | /v1/posts/{post_id}/comments  | Add a comment to a blog    |

## API Documentation

OpenAPI v2 (Swagger) documentation is automatically generated in the `docs` directory when running `buf generate`. The documentation provides a detailed description of all API endpoints, request/response schemas, and available operations.

The generated documentation can be found at:
```
docs/protos/blog/v1/blog.swagger.json
```

This Swagger JSON file can be used with tools like [Swagger UI](https://swagger.io/tools/swagger-ui/) to visualize and interact with the API.

## Validation

Field validation is implemented using buf validate. The following validations are applied:
- UUID: Must follow the standard UUID format (e.g., 123e4567-e89b-12d3-a456-426614174000)
- Blog title: 1-100 characters, alphanumeric with basic punctuation
- Blog content: 1-10000 characters
- Comment content: 1-1000 characters
- Comment author: 1-50 characters
- Page size for listing: 1-100 items

## License

This project is licensed under the MIT License - see the LICENSE file for details.
