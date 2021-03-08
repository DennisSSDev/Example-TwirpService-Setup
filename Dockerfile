FROM golang:latest AS build
WORKDIR /go/src/github.com/dennisssdev/Example-TwirpService-Setup
COPY . .
RUN make build

FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build /go/bin/Example-TwirpService-Setup /app/
EXPOSE 8000
ENTRYPOINT ["./Example-TwirpService-Setup"]