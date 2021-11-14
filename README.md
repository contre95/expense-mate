# Expense App
My new excuse to keep on learning DDD and write some Go code. Don't expect much doc.

It's just an API made with [Fiber](https://github.com/gofiber/fiber) and [Gorm](https://gorm.io/) (not the best decision) which aims the help you with tracking expenses and let's you interact with them on [Grafana](https://grafana.com/)

## Features
* Import expenses from different source
    * Google sheets 
    * Sample Data ([exampleImporter.go](https://github.com/contre95/expenses-app/blob/main/pkg/gateways/importers/exampleImporter.go)) 
* Add / Delete individual expenses (TODO)


# Configurations
All configurations are set in the `.env` file and passed as environment variables
```sh
# Set the .env
mv .env.example .env
# Install the dependencies
go mod tidy
# Source the env variables
. <(cat .env | grep -v -e '^$' | grep -v "#" | awk '{print "export " $1}')
```

# Run 
```sh
# Development environment
docker-compose up -d
# The app
go run main.go
```

# Endpoints

## Google Sheets Importer
```sh
# /api/v1/importers/:importer_id
curl -H "Content-Type: application/json" -X POST \
    -d '{ "bypass_wrong_expenses": true }' \
    -X POST http://localhost:3000/api/v1/importers/sheets  | jq
```
```json
{
  "err": null,
  "msg": {
    "Msg": "All the expenses where imported",
    "SuccesfullImports": 206,
    "FailedImports": 0
  },
  "success": true
}
```

## Healthcheck
```sh
# /ping
curl -H "Content-Type: application/json" -X GET http://localhost:3000/ping | jq
```
```json
{
  "ping": "pong"
}
```
   
   
   
   
   
   
