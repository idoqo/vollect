FROM golang:1.15.2-alpine3.12 as builder

COPY go.mod go.sum /go/src/gitlab.com/idoko/vollect/
WORKDIR /go/src/gitlab.com/idoko/vollect
RUN go mod download
COPY . /go/src/gitlab.com/idoko/vollect
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/vollect gitlab.com/idoko/vollect

FROM alpine

RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/gitlab.com/idoko/vollect/build/vollect /usr/bin/vollect

EXPOSE 8080 8080

ENTRYPOINT ["/usr/bin/vollect"]