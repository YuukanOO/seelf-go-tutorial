# Containerizing and deploying a Go application with seelf

This simple repository is here to demonstrate how you can easily containerize a Go application and deploy its stack using [seelf](https://github.com/YuukanOO/seelf).

The repository is divided into branches representing tutorial steps:

1. Base application ([main](https://github.com/YuukanOO/seelf-go-tutorial/tree/main))
1. Containerizing the application ([containerizing](https://github.com/YuukanOO/seelf-go-tutorial/tree/containerizing))
1. Deploying it on a seelf instance ([deploying](https://github.com/YuukanOO/seelf-go-tutorial/tree/deploying))

## Application

This is a web application which output the number of times it was started. To demonstrate how it's easy to deploy a stack on seelf, it persists its data in a PostgreSQL database.

### Building & running

```sh
docker compose up -d --build --wait
```
