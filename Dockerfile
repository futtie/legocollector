FROM golang:1.18-alpine as builder
RUN apk add --no-cache git
ENV CGO_ENABLED=0
ENV GOOS=linux 
ENV GO111MODULE=on
WORKDIR /app
COPY . /app
RUN go build -a -installsuffix cgo -o legocollector .
# RUN go get github.com/go-swagger/go-swagger/cmd/swagger@0.16.0
# RUN swagger generate spec -o /app/swagger.json
# RUN sed -i 's/@bamboo.buildNumber/${bamboo.buildNumber}/g' /app/swagger.json

FROM scratch
COPY --from=builder /app/legocollector /legocollector
ENTRYPOINT ["/legocollector"]
