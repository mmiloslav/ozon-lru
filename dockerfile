FROM golang:1.22

WORKDIR /go

COPY ./api/ /go/api
COPY ./cmd/ /go/cmd
COPY ./deployments/ /go/deployments
COPY ./internal/ /go/internal
COPY ./pkg/ /go/pkg

COPY ./go.mod /go/
COPY ./go.sum /go/

RUN go build -o main ./cmd/lruapp/main.go
EXPOSE 8080

ENTRYPOINT ["./main"]