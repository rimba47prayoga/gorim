# gorim
Golang Framework inspired by Django

Gorim features

- Global Settings
    - Database
    - Server
    - Middlewares

- GenericViewSet
    - GetQuerySet
    - GetObject
    - PKField: unique identifier for retrieving single object (default: id)
    - Filter
    - Pagination

- Mixins
    - CreateMixin
    - UpdateMixin
    - ListMixin
    - RetrieveMixin

- Serializer
    - support single field validation
    - default Create & Update

- Permissions
    - DEFAULT_PERMISSION_STRUCTS as global settings (default: permissions.IsAuthenticated)

- DefaultRouter
    - Type Parameter with validation (/users/<int:id>)
    - DefaultRouter.RegisterFunc: for adding new route from extra handler function

- Middlewares
    - Change to struct instead of function
    - RecoverMiddleware: for seamless error handling
    - LoggerMiddleware: request logger in terminal with time ellapsed

- Errors
    - errors.Raise: to stop current process and immediately Response to client

- utils
    - HasAttr: checks if a struct has a field or method with a given name
    - GetObjectOr404: Generic function to get object or return 404

- Migrations
    - Generate tables for tracking migration
    - Track migration with versioning hash
    - Datamigration / Query execution with only once applied

- cmd
    - runserver: only runserver without migration
    - warning message if there's unapplied migration
    - migrate: run migrations
