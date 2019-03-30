FROM golang:1.12 AS build

COPY .  /the-dummies-go

WORKDIR /the-dummies-go

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/the-dummies-go main.go

FROM scratch
COPY --from=build /bin/the-dummies-go /bin/the-dummies-go
ENTRYPOINT ["/bin/the-dummies-go"]
CMD []
