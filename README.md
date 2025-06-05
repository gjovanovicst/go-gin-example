# Go Gin Example [![rcard](https://goreportcard.com/badge/github.com/EDDYCJY/go-gin-example)](https://goreportcard.com/report/github.com/EDDYCJY/go-gin-example) [![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/EDDYCJY/go-gin-example) [![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/EDDYCJY/go-gin-example/master/LICENSE)

An example of gin contains many useful features

[简体中文](https://github.com/EDDYCJY/go-gin-example/blob/master/README_ZH.md)

## Installation
```
$ go get github.com/EDDYCJY/go-gin-example
```

## How to run

### Required

- Mysql
- Redis

### Ready

Create a **blog database**. The database schema and seed data will be automatically created using the migration system when you start the application.

**Note**: The application now uses an automated migration system instead of manual SQL imports. See [Database Migrations](#database-migrations) section below for more details.

### Conf

You should modify `conf/app.ini`

```
[database]
Type = mysql
User = root
Password =
Host = 127.0.0.1:3306
Name = blog
TablePrefix = blog_

[redis]
Host = 127.0.0.1:6379
Password =
MaxIdle = 30
MaxActive = 30
IdleTimeout = 200
...
```

### Run

#### Development with Live Reload (Recommended)
For development with automatic reloading when files change (similar to nodemon):

```bash
# Install Air (Go live reload tool)
$ go install github.com/air-verse/air@latest

# Initialize Air configuration (creates .air.toml)
$ air init

# Start development server with live reload
$ air
```

The server will automatically restart when you make changes to any Go files.

#### Standard Run
```bash
$ cd $GOPATH/src/go-gin-example

$ go run main.go
```

Project information and existing API

```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /auth                     --> github.com/EDDYCJY/go-gin-example/routers/api.GetAuth (3 handlers)
[GIN-debug] GET    /swagger/*any             --> github.com/EDDYCJY/go-gin-example/vendor/github.com/swaggo/gin-swagger.WrapHandler.func1 (3 handlers)
[GIN-debug] GET    /api/v1/tags              --> github.com/EDDYCJY/go-gin-example/routers/api/v1.GetTags (4 handlers)
[GIN-debug] POST   /api/v1/tags              --> github.com/EDDYCJY/go-gin-example/routers/api/v1.AddTag (4 handlers)
[GIN-debug] PUT    /api/v1/tags/:id          --> github.com/EDDYCJY/go-gin-example/routers/api/v1.EditTag (4 handlers)
[GIN-debug] DELETE /api/v1/tags/:id          --> github.com/EDDYCJY/go-gin-example/routers/api/v1.DeleteTag (4 handlers)
[GIN-debug] GET    /api/v1/articles          --> github.com/EDDYCJY/go-gin-example/routers/api/v1.GetArticles (4 handlers)
[GIN-debug] GET    /api/v1/articles/:id      --> github.com/EDDYCJY/go-gin-example/routers/api/v1.GetArticle (4 handlers)
[GIN-debug] POST   /api/v1/articles          --> github.com/EDDYCJY/go-gin-example/routers/api/v1.AddArticle (4 handlers)
[GIN-debug] PUT    /api/v1/articles/:id      --> github.com/EDDYCJY/go-gin-example/routers/api/v1.EditArticle (4 handlers)
[GIN-debug] DELETE /api/v1/articles/:id      --> github.com/EDDYCJY/go-gin-example/routers/api/v1.DeleteArticle (4 handlers)

Listening port is 8000
Actual pid is 4393
```
Swagger doc

![image](https://i.imgur.com/bVRLTP4.jpg)

## Database Migrations

This project uses an automated database migration system that replaces the manual SQL setup. The migrations will run automatically when the application starts, ensuring your database schema is always up-to-date.

### Migration Features

- **Automatic migrations**: Runs on application startup
- **Version control**: Tracks migration versions
- **Rollback support**: Ability to rollback changes
- **Seed data**: Automated insertion of initial data
- **CLI management**: Manual migration control

### Migration Commands

```bash
# Run all pending migrations
make migrate-up

# Rollback the last migration
make migrate-down

# Check current migration version
make migrate-version

# Migrate to a specific version
make migrate-to VERSION=1

# Create new migration files
make new-migration NAME=add_user_email_column

# Reset database (WARNING: drops all data)
make migrate-reset
```

### Manual Migration Management

You can also use the migration CLI tool directly:

```bash
# Build the migration tool
go build -o bin/migrate cmd/migrate/main.go

# Run migrations
./bin/migrate -action=up

# Check version
./bin/migrate -action=version

# Rollback
./bin/migrate -action=down

# Migrate to specific version
./bin/migrate -action=migrate -version=1
```

For detailed migration documentation, see [MIGRATIONS.md](MIGRATIONS.md).

## Features

- RESTful API
- Gorm
- Swagger
- logging
- Jwt-go
- Gin
- Graceful restart or stop (fvbock/endless)
- App configurable
- Cron
- Redis
- Live reload development with Air (auto-restart on file changes)
- **Database migrations and seeding**