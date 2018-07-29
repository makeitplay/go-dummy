FROM golang:1.10 AS build

RUN go get github.com/golang/dep/cmd/dep

COPY .  /go/src/github.com/makeitplay/the-dummies-go

WORKDIR /go/src/github.com/makeitplay/the-dummies-go

RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/the-dummies-go

FROM scratch
COPY --from=build /bin/the-dummies-go /bin/the-dummies-go
ENTRYPOINT ["/bin/the-dummies-go"]
CMD []
