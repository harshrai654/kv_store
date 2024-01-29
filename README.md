# Key Value Store Based On MySQL

This is a key value store based on Relational database MySQL.

## Architecture

This is a project to implement the concept of soft delete and TTL in a key value store similar to Redis.
There is User facing API written in Go, which is used to interact with the mysql DB. It provides basic CRUD operations.

## Getting Started

- Install docker and docker-compose
- Create .env file
  - Example .env file
    ```env
    MYSQL_USER = kv_user    #user created as part of db init
    MYSQL_PASSWORD = 1234567890@Abcd
    MYSQL_DATABASE = kv_store
    DB_HOST = kv-mysql
    DB_PORT = 3306
    PORT = 5501
    MYSQL_ROOT_PASSWORD=1234567890@Abcd
    ```
- Run `docker-compose up --build -d` to start api and mysqll containers
- Run `docker-compose logs -f` to follow logs
- Run `docker-compose down` to stop containers
