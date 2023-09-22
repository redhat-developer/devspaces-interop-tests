FROM registry.access.redhat.com/ubi8/go-toolset:1.18.4-8 as builder

WORKDIR /github.com/redhat-developer/devspaces-interop-tests
USER root

# compile test binary
COPY . .
RUN make

FROM registry.access.redhat.com/ubi8/go-toolset:1.18.4-8
WORKDIR /test-harness/
COPY --from=builder /github.com/redhat-developer/devspaces-interop-tests/bin/che-operator-test-harness ./che-operator-test-harness
ENTRYPOINT [ "/test-harness/che-operator-test-harness" ]
