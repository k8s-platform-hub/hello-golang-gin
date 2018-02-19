FROM golang:1.8.5-jessie

# # uncomment to install required debian packages
# # add any package that is required after `build-essential`, end the line with \
# RUN apt-get update && apt-get install -y \
#     build-essential \
# && rm -rf /var/lib/apt/lists/*

# install dep and gin
RUN go get github.com/golang/dep/cmd/dep github.com/codegangsta/gin

# setup the working directory
WORKDIR /go/src/app
COPY Gopkg.toml Gopkg.lock ./

# install dependencies
RUN dep ensure -vendor-only

# add source code
COPY src src
WORKDIR src
# build the source
RUN go build -o main *.go

# command to be executed on running the container
CMD ["./main"]

