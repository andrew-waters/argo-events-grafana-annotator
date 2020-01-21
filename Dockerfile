FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM golang:1.13-rc-alpine as build
WORKDIR /go/src/github.com/andrew-waters/argocd-grafana-annotator
RUN apk update && apk upgrade && \
  apk add --no-cache bash git openssh
RUN adduser -D -g '' binuser
COPY ./ /go/src/github.com/andrew-waters/argocd-grafana-annotator
RUN go get -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o binary .

FROM scratch
ENV PATH=/bin
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/andrew-waters/argocd-grafana-annotator/binary /bin/argocd-grafana-annotator
COPY --from=build /etc/passwd /etc/passwd
USER binuser
CMD ["./bin/argocd-grafana-annotator"]
