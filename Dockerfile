FROM golang:1.16.15 AS builder
WORKDIR /opt/kuligo
COPY . /opt/kuligo/
RUN GOOS=linux GOARCH=amd64 go build -o kuligo main.go

FROM centos:7.9.2009 AS runner
WORKDIR /opt/kuligo
COPY --from=builder /opt/kuligo/kuligo .
COPY views/ /opt/kuligo/views/
EXPOSE 10099
CMD ["./kuligo"]