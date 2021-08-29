# syntax=docker/dockerfile:1

# An attempt to make the deployment multistage. Not sure if I succeeded.

FROM golang:1.16-alpine AS build

RUN mkdir -p  /pservice
WORKDIR /pservice

COPY go.* ./
RUN go mod download

COPY main.go ./

RUN mkdir -p ./ptlist
COPY ptlist/*.go ./ptlist/

RUN mkdir -p ./time_utils
COPY time_utils/*.go ./time_utils/

RUN CGO_ENABLED=0 go build -o /pservice/pservice.out

FROM gcr.io/distroless/base-debian10
WORKDIR /pservice
COPY --from=build /pservice /pservice
# run with docker run -p 8383, to override
EXPOSE 8282
USER nonroot:nonroot
ENTRYPOINT [ "/pservice/pservice.out" ]
