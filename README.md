<img src="./public/assets/img/logo.png)">
# ExpenseMate

## Expenses dashboard [Refactoring]

A Grafana dashboard is automatically created and can be accesses from [localhost:8080](http://localhost:8080)

![image](https://user-images.githubusercontent.com/15664513/216789116-86d3cf33-5535-4bb9-b30c-8196c5ef1696.png)

# Run locally
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
