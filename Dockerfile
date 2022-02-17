FROM golang:1.17 as dev

ARG REPO_USER
ARG REPO_KEY

WORKDIR /go/shorturlapi

COPY . .

RUN go get ./...

EXPOSE 8080

HEALTHCHECK --interval=1m --timeout=30s --start-period=1m --retries=3 \
  CMD curl -f http://localhost:8080/healthcheck || exit 1

CMD ["go", "run", "/go/shorturlapi/main.go"]



FROM golang:1.17 as build

ARG REPO_USER
ARG REPO_KEY

WORKDIR /go/shorturlapi

COPY . .

COPY --from=dev /go/shorturlapi .

RUN CGO_ENABLED=0 GOOS=linux go build -a .




FROM golang:1.17-alpine as prod

COPY --from=build /go/shorturlapi /shorturlapi

EXPOSE 8080

HEALTHCHECK --interval=1m --timeout=30s --start-period=1m --retries=3 \
  CMD curl -f http://localhost:8080/healthcheck || exit 1

ENTRYPOINT ["/shorturlapi"]