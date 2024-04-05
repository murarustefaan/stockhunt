FROM golang:1.22-alpine as builder

RUN apk add --no-cache git gcc musl-dev

WORKDIR /service

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=1 go build -o out/stockhunt .
RUN ldd out/stockhunt | tr -s [:blank:] '\n' | grep ^/ | xargs -I % install -D % out/%
RUN ls -l out/
RUN ls -l out/lib/


FROM gcr.io/distroless/static AS runtime
LABEL maintainer="Stefan MURARU <muraru.stefaan@gmail.com>"
LABEL description="StockHunt - A stock market analysis tool"

USER nonroot

WORKDIR /service

COPY --from=builder --chown=nonroot /service/out/stockhunt ./
COPY --from=builder --chown=nonroot /service/out/lib/ /lib/
COPY --chown=nonroot ./migrations ./migrations

ENTRYPOINT ["/service/stockhunt"]