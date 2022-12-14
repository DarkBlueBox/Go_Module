FROM heroiclabs/nakama-pluginbuilder:3.12.0 AS go-builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /backend

COPY . .

RUN go build --trimpath --mod=vendor --buildmode=plugin -o ./backend.so

FROM heroiclabs/nakama:3.12.0

COPY --from=go-builder /backend/backend.so /nakama/data/modules/
COPY local.yml /nakama/data/