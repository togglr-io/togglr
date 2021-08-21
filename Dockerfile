FROM golang:1.17 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /opt/app

COPY . /opt/app

RUN go build -o server cmd/server/main.go
RUN go build -o migrate cmd/migrate/main.go


FROM scratch

ENV TOGGLE_HOST=0.0.0.0
ENV TOGGLE_PORT=9001
ENV TOGGLE_DB_HOST=postgres
ENV TOGGLE_DB_PORT=5432
ENV TOGGLE_DB_USER=toggle
ENV TOGGLE_DB_PASSWORD=toggle

WORKDIR /opt/app
COPY --from=build /opt/app /opt/app
