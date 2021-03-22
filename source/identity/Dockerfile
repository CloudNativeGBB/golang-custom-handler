FROM golang:latest AS builder

WORKDIR /app
ADD handler/. /app
RUN go mod download && go get -d -v && go build -o main .

FROM mcr.microsoft.com/azure-functions/dotnet:3.0

ENV AzureWebJobsScriptRoot=/home/site/wwwroot \
    AzureFunctionsJobHost__Logging__Console__IsEnabled=true

WORKDIR /home/site/wwwroot
COPY --from=builder /app/main .
COPY function/. .
RUN rm host.json
RUN mv host-container.json host.json

EXPOSE 80