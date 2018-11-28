FROM golang:1.11-alpine3.8 
COPY ./cloud_native_go /app/Cloud-Native-Go
RUN chmod +x /app/Cloud-Native-Go

ENV PORT 8080
EXPOSE 8080

ENTRYPOINT /app/cloud_native_go