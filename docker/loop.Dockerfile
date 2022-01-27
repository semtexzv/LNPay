FROM registry.fedoraproject.org/fedora-minimal:33 as builder
RUN microdnf install golang

RUN git clone https://github.com/lightninglabs/loop.git
RUN cd loop/cmd && GOBIN=/built go install ./...


FROM registry.fedoraproject.org/fedora-minimal:33
COPY --from=builder /built/ /usr/bin/

ENTRYPOINT ["loopd"]
