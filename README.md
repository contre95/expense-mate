# 📊 ExpenseMate :mate:
ExpenseMate is an expense tracking app with a handy telegram bot and a front end. It uses Go, MySQL, SQLite, Fiber, Tailwindcss and HTMX.

![image](https://github.com/contre95/expenses-app/assets/15664513/df1d0fc1-12a8-488e-940c-d950c1916948)

### Grafana (`refactoring`)
![image](https://user-images.githubusercontent.com/15664513/216789116-86d3cf33-5535-4bb9-b30c-8196c5ef1696.png)

### Run locally
All configurations are set in the `.env` file and passed as environment variables

```sh
# Install the dependencies
go mod tidy
# Set the .env
mv .env.example .env
# Source the env variables
. <(cat .env | grep -v -e '^$' | grep -v "#" | awk '{print "export " $1}')
# Run with air
air -c air.toml
```
### Container :whale:
```sh 
  docker run -d \
  --name expenses-app \
  --restart always \
  --env STORAGE_ENGINE=sqlite \
  --env LOAD_SAMPLE_DATA=true \
  --env SQLITE_PATH=./exp.db \
  --env JSON_STORAGE_PATH=./users.json \
  --env TELEGRAM_APITOKEN= \
  -p 8080:8080 \
  contre95/expense-mate:latest
```

### Telegram `/help`
```
Check the menu for available commands, please.
/categories - Sends you all the categories available.
/summary - Sends summar of last month's expenses.
/unknown - Categorize unknown expenses. /done and continue in another moment.
/new - Creates a new expense. /fix if you made made a mistake.
/ping - Checks bot availability and health.
/help - Displays this menu.
