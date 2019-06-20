# Dockerfile References: https://docs.docker.com/engine/reference/builder/


# Start from golang v1.11 base image
FROM golang:1.11

LABEL maintainer="Percy Bolmér"

# BUILD ARG
ARG APP_NAME=dockerexample
ARG LOG_DIR=/${APP_NAME}/logs

# Create Log directory
RUN mkdir -p ${LOG_DIR}

# ENviroment Variables
ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log


# WORKDIR is the Current Working Directory inside the Container
WORKDIR $GOPATH/src/dockerexample
# Copy all files from Current WD into PWD 
COPY . .

# DOWNLOAD ALL DEP? 
RUN go get -d -v ./...

RUN go install -v ./...

RUN apt-get update && apt-get install -y ffmpeg


WORKDIR $GOPATH/src/dockerexample

# OPEN UP port 8080 
EXPOSE 8080

# Declare a Volume to mount to the docker
# ~/app-logs:/dockerexample/logs dockerexample
# the docker run specifies this to Mount
# docker run -v HOSTPATH:CONTAINERPATH
VOLUME ["/dockerexample/logs"]

# CMD is the commands to run in CMD mode, We want to start our server
CMD ["dockerexample"]