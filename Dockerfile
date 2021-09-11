FROM alpine AS base
ENV PORT=80
WORKDIR /app
EXPOSE 80

FROM golang:1.16-alpine AS deps
WORKDIR /src
COPY go.mod .
COPY go.sum .
RUN go mod download

FROM deps AS build
WORKDIR /src
COPY . .
RUN GOOS=linux CGO_ENABLED=1 GOARCH=amd64 go build -o /app/congress

FROM base AS final
WORKDIR /app
COPY --from=build /app/congress .
ENTRYPOINT ["/app/congress"]
