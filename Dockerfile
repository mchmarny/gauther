# build from latest go image
FROM golang:latest as build

WORKDIR /go/src/github.com/mchmarny/gauther/
COPY . /src/

# build gauther
WORKDIR /src/
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 go build -o /gauther



# run image
FROM scratch

# certs from build to avoid "certificate signed by unknown authority" error
# you can also build from golang:alpine ...
# and run "apk --no-cache add ca-certificates"
COPY --from=build /src/certs/ca-certificates.crt /etc/ssl/certs/

# copy app executable
COPY --from=build /gauther /app/

# copy static dependancies
COPY --from=build /src/templates /app/templates/
COPY --from=build /src/static /app/static/

# start server
WORKDIR /app
ENTRYPOINT ["./gauther"]