FROM node:latest AS build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./ .
ARG VITE_OAUTH_PROVIDER_URL
ARG VITE_OAUTH_CLIENT_ID
ENV VITE_OAUTH_PROVIDER_URL=$VITE_OAUTH_PROVIDER_URL
ENV VITE_OAUTH_CLIENT_ID=$VITE_OAUTH_CLIENT_ID
RUN npm run build

FROM nginx AS production-stage
RUN mkdir /app
COPY --from=build-stage /app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf