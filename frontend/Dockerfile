# syntax=docker/dockerfile:1

# Stage 1: Build the Vite Vue.js app
FROM node:24 AS builder

ARG VITE_API_URL
ENV VITE_API_URL=$VITE_API_URL

WORKDIR /app

# Copy the package.json and yarn.lock files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy the project files
COPY . .

# Build the app
RUN npm run build

# Stage 2: Run Nginx with the built app
FROM nginx:alpine

# Remove the default Nginx configuration
RUN rm -rf /etc/nginx/conf.d/*

# Add new config compatible with Vue.JS
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

# Copy the built app from the previous stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Expose the Nginx port
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]
