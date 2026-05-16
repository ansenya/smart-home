FROM node:latest AS build-stage
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY ./ .
ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

FROM nginx AS production-stage
RUN mkdir /app
COPY --from=build-stage /app/dist /app
COPY nginx.conf /etc/nginx/nginx.conf
