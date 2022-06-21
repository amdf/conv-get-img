FROM golang:1.17-alpine as builder
ENV GOPATH=
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app ./cmd

FROM alpine:3
COPY --from=builder app /bin/app
#COPY *.ttf ./
ENTRYPOINT ["/bin/app"]
