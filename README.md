# go-mysql-radix-demo
Very simple demonstration of a Radix hosted api written in go accessing an Azure MySql service

**DB_PASS**: username:password@tcp(mysqlinstancename.mysql.database.azure.com:3306)/demo?tls=true

## Run locally

You need Go installed.

The DB_CONN environment variable must be set to a valid MySQL connection string, e.g. *username:password@tcp(mysqlinstancename.mysql.database.azure.com:3306)/demo?tls=true*

DB_CONN can be set in launch.json if you are using VS Code.


