FROM nginx:alpine

WORKDIR /usr/share/nginx/html

COPY nginx /etc/nginx

EXPOSE 80
