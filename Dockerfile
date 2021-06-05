FROM golang:1.13 as build
WORKDIR /app
COPY . .
ENV GOPROXY=https://goproxy.io,direct
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
  go build -ldflags -s -a -installsuffix cgo -o bin/elf-server cmd/*.go

FROM scratch as run
WORKDIR /app
VOLUME [ "/app/upload" ]
EXPOSE 5000
COPY --from=build /app/bin/elf-server .
COPY resources .
COPY theme .
ENTRYPOINT [ "/app/elf-server" ]
