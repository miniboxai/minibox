FROM golang:1.11.2-stretch

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN make go-build-cli
RUN make go-build

CMD ["app"]