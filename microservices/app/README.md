# go-gin

[Gin Gonic](https://gin-gonic.github.io/gin/) (a full-featured web framework for [Golang](https://golang.org/)) microservice with [Dep](https://golang.github.io/dep/) for package management [codegangsta/gin](https://github.com/codegangsta/gin) for live reloading in local development.

## Managing dependencies

### Adding a new Golang package

If you need to install a package, say `github.com/gin-contrib/authz`
```bash
# with docker

$ cd microservices/app
$ docker build -t hello-golang-gin-app .
$ docker run --rm -it -v $(pwd):/go/src/app \
             hello-golang-gin-app \
             dep ensure -add github.com/gin-contrib/authz


# without docker

$ cd microservices/app
$ dep ensure -add github.com/gin-contrib/authz
```
This will update `Gopkg.toml` and `Gopkg.lock`.

### Adding a new system package

The base image used in this boilerplate is [golang:1.8.5-jessie](https://hub.docker.com/_/golang/). Hence, all debian packages are available for installation. You can add a package by mentioning it in the `Dockerfile` among the existing apt-get install packages.

```dockerfile
FROM golang:1.8.5-jessie

# install required debian packages
# add any package that is required after `build-essential`, end the line with \
RUN apt-get update && apt-get install -y \
    build-essential \
&& rm -rf /var/lib/apt/lists/*

...
```

## Local development

### With Docker

- Install [Docker CE](https://docs.docker.com/engine/installation/)

```bash
# go to app directory
$ cd microservices/app

# build the docker image
$ docker build -t hello-golang-gin-app .

# run the image using either 1 or 2

# 1) without live reloading
$ docker run --rm -it -p 8080:8080 hello-golang-gin-app

# 2) with live reloading
# any change you make to your source code will be immediately updated on the running app
$ docker run --rm -it -p 8080:8080 \
             -v $(pwd):/go/src/app \
             hello-golang-gin-app \
             bash
root@<id>:/go/src/app/src$ gin --bin main-bin --port 8080 run

# app will be available at http://localhost:8080
# press Ctrl+C to stop the server
```

### Without Docker

- Install [Golang](https://golang.org/doc/install)
- Move the `hello-golang-gin` directory to your `GOPATH` and cd into the directory

```bash
# change to app directory
$ cd mircoservices/app

# install dep for package management
$ go get github.com/golang/dep/cmd/dep
# install gin for live reloading
$ go get github.com/codegangsta/gin

# install dependencies
$ dep ensure -vendor-only

$ cd src
$ gin --bin main-bin --port 8080 run

# app will be available at http://localhost:8080
# any change you make to your source code will be immediately updated on the running app
# press Ctrl+C to stop the server
```
