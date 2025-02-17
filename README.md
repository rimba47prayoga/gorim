# Gorim  
**A Golang Framework Inspired by Django and Django Rest Framework**  
**Built on top of Echo**

[![Go Reference](https://pkg.go.dev/badge/gorim.org/gorim.svg)](https://pkg.go.dev/gorim.org/gorim)  
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)  

Gorim is a backend framework built with Golang that adopts the architecture and design of Django and Django Rest Framework (DRF).  
It is built on top of the [Echo](https://echo.labstack.com/) web framework, providing a high-performance, minimalist web server with a powerful routing system.  
Gorim simplifies backend application development with class-based views, serializers, middlewares, and other features similar to Django.  

## ✨ Key Features  
- **Built on Echo**: Leverages Echo's high-performance routing and middleware system.  
- **Global Settings**: Configuration for database, server, and middleware.  
- **Class-Based Views**: `GenericViewSet` with features like `GetQuerySet`, `GetObject`, `Filter`, and `Pagination`.  
- **Mixins**: `CreateMixin`, `UpdateMixin`, `ListMixin`, `RetrieveMixin`.  
- **Serializer**: Supports field validation and default Create & Update.  
- **Permissions**: Permission system with `DEFAULT_PERMISSION_STRUCTS`.  
- **DefaultRouter**: Typed parameter support (`/users/<int:id>`) and `RegisterFunc` for adding routes.  
- **Middlewares**: Struct-based middleware with `RecoverMiddleware` and `LoggerMiddleware`.  
- **Error Handling**: `errors.Raise` to stop the current process and immediately return a response.  
- **Migrations**: Track migration versioning and apply data migrations with single execution.  
- **Command Line (CLI)**: Use commands like `runserver`, `migrate`, and warnings for unapplied migrations.  

## 🚀 Installation  
Make sure Go is installed, then run the following command:  

```sh
go get gorim.org/gorim
```