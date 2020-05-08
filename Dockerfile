FROM golang:1.14-alpine AS build
RUN apk --no-cache update && apk --no-cache add ca-certificates git
WORKDIR /go/src/rollbar-drone
COPY . ./
RUN	GO111MODULE=on CGO_ENABLED=0 go build -o bin/rollbar-drone

FROM scratch AS rollbar-plugin
COPY --from=build /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=build /go/src/rollbar-drone/bin/rollbar-drone /rollbar-drone
ENTRYPOINT ["/rollbar-drone"]
