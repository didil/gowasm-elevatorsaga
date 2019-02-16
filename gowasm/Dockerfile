FROM golang:1.11.5-alpine

WORKDIR /app/

ENV GOARCH wasm
ENV GOOS js

CMD ["go", "build", "-o", "app.wasm" , "."]