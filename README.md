# 📊🧉 ExpenseMate
ExpenseMate is an expense tracking app with a handy Telegram bot and a front end. It uses [Go](https://go.dev/), [MySQL](https://www.mysql.com/), [SQLite](https://www.sqlite.org/index.html), [Fiber](https://gofiber.io/), [TailwindCSS](https://tailwindcss.com/) and [HTMX](https://htmx.org/).

### Features
* Create, Update, Delete or Modify expenses. (Web) 
* Import expenses from [N26](https://n26.com/en-es) CSV extract. (Web)
* Create Rules to automatically categorize imported expenses. (Web)
* Manage multiple users and attach them to expenses. (Web)
* Get prompted expenses without category and categorize them. (Telegram) 
* Create expenses. (Telegram) 
* Get an expense summary by category. (Telegram) 

You can try a [demo](https://demo1.contre.io), it resets every 3hs. Telegram bot is not enabled.

### 🦭 Run in a container 
```sh
podman run -d \
  --name expense-mate \
  --restart always \
  --env STORAGE_ENGINE="sqlite"  \  # Required (sqlite|mysql)
  --env SQLITE_PATH="./exp.db" \  # Required
  --env LOAD_SAMPLE_DATA="true" \  # Optional (default:true)
  --env VISION_MODEL="llama3.2-vision:11b-instruct-q4_K_M"  # Optional
  --env TEXT_MODEL="llama3.2:3b-instruct-q6_K"  # Optional
  --env OLLAMA_ENDPOINT="http://localhost:11434/api/generate"  # Optional
  --env JSON_STORAGE_PATH="./users.json" \ # Required
  --env TELEGRAM_APITOKEN="<TELEGRAM_BOT_TOKEN>"\ # Optional
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
/ai - Analyze image/text for expenses. Send /cancel to quit.
/help - Displays this menu.
```
