FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

ARG VITE_API_URL=/api/v1
ARG VITE_YANDEX_GEOCODER_API_KEY=bca105ec-e532-40dc-8c78-5a3ae701d5bf


ENV VITE_API_URL=$VITE_API_URL
ENV VITE_YANDEX_GEOCODER_API_KEY=$VITE_YANDEX_GEOCODER_API_KEY

RUN npm run build

FROM nginx:1.27-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80
