# 📊🧉 ExpenseMate
ExpenseMate is an expense tracking app with a handy Telegram bot and a front end.
Stack:
 * Go
 * MySQL or SQLite
 * HTMX

You can try a [demo](https://demo1.contre.io), it resets every 3hs. Telegram bot is not enabled.

### 🦭 Run in a container 
```sh
podman run -d \
  --name expenses-app \
  --restart always \
  --env STORAGE_ENGINE=sqlite  \
  --env SQLITE_PATH=./exp.db \
  --env LOAD_SAMPLE_DATA=true \
  --env JSON_STORAGE_PATH=./users.json \
  --env TELEGRAM_APITOKEN= \
  -p 8080:8080 \
  contre95/expense-mate:latest
```

### 💻 Run locally
All configurations are set in the `.env` file and passed as environment variables. You can access from [localhost:8080](http://localhost:8080)
```sh
  # Clone the repository
  git clone https://github.com/contre95/expense-mate.git && cd expense-mate
  # Install the dependencies
  go mod tidy
  # Set the .env (you can leave the defaults for a quick start)
  mv .env.example .env
  # Source the env variables
  . <(cat .env | grep -v -e '^$' | grep -v "#" | awk '{print "export " $1}')
  # Run with go
  go run ./cmd/main.go
  # Alternatively run with air
  # air -c air.toml
```


### 🟦 Telegram
Optionally you can set up a Telegram bot to interact with the app. All you need is a new bot that can be creating chatting to [@BotFather](https://t.me/BotFather)
```
Check the menu for available commands, please.
/categories - Sends you all the categories available.
/summary - Sends summar of last month's expenses.
/unknown - Categorize unknown expenses. /done and continue in another moment.
/new - Creates a new expense. /fix if you made made a mistake.
/ping - Checks bot availability and health.
/help - Displays this menu.
```
