FROM registry.fedoraproject.org/fedora-minimal:33 as builder
RUN microdnf install golang

ADD go.mod .
ADD go.sum .
RUN GOPROXY=https://proxy.golang.org go mod download

ADD . .

RUN go build main.go

FROM registry.fedoraproject.org/fedora-minimal:33
COPY --from=builder main .

CMD ["./main"]
