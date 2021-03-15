FROM golang:1.14.3-alpine 
WORKDIR /
COPY ./src /
RUN go build -o /app
ENTRYPOINT /app
