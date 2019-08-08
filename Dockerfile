FROM golang as builder

ADD go.??? /go/ursho/

WORKDIR /go/ursho/

RUN go mod download

COPY . /go/ursho

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ursho .

FROM scratch

ENV PORT 8080

COPY --from=builder /go/ursho/ursho /app/
ADD config/config.json /app/config/

WORKDIR /app

CMD ["./ursho"]
