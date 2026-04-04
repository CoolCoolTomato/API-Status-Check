FROM node:22-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile
COPY web/ ./
RUN pnpm build

FROM golang:1.25-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
RUN go build -o server ./cmd/server/

FROM alpine:latest
WORKDIR /app
RUN mkdir -p data
COPY --from=backend /app/server .
COPY --from=backend /app/web/dist ./web/dist
EXPOSE 8080
CMD ["./server"]
