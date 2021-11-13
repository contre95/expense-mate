# Expense App
My new excuse to keep on learning DDD and write some Go code.


# Configurations
To load the `.env` file please run the following command 
```sh
. <(cat .env | grep -v -e '^$' | grep -v "#" | awk '{print "export " $1}')
```
