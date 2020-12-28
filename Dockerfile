FROM golang:1.15.3-buster AS build

WORKDIR /root

COPY . .

RUN go mod download

RUN make tests
RUN make build-static

FROM scratch

COPY --from=build /root/go-there /bin/go-there

ENTRYPOINT ["/bin/go-there"]