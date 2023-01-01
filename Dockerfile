FROM golang:1.18-alpine as builder
RUN apk add --no-cache git
ENV CGO_ENABLED=0
ENV GOOS=linux 
ENV GO111MODULE=on
WORKDIR /app
COPY . /app
RUN go build -a -installsuffix cgo -o legocollector .

FROM scratch
COPY --from=builder /app/legocollector /legocollector
COPY --from=builder /app/files/ /files/
ENTRYPOINT ["/legocollector"]
