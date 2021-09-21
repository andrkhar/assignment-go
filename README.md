# Coding Assignment in GO

## Services
1. **mysql** - container with a database
2. **echo** - server with API to store and retrieve data to the database.
3. **generator** - server with API to run a generator that makes requests to **echo** API to populate the database.
4. **dev** - container used for development only. It has vim installed with vim-go plugin enabled.


## Usage
```
cd src
cp template.env .env
```
Edit `.env` file to add MYSQL credentials

Run containers:
```
docker-compose up -d generator
```

*It runs 3 services: **generator**, **echo** and **mysql***.

Wait several seconds for database availability.

## Testing
### echo - port :80
1. Get all records from database: **GET** localhost:80/data
2. Get filtered records: **GET** localhost/data?start_timestamp=1632062101&end_timestamp=1632062111&ID1=1&ID2=F

### generator - port :8080
1. Start data generation: **POST** localhost:8080/start
2. Stop: **POST** localhost:8080/stop

## Developing
```
cd src
docker-compose up -d dev
docker attach src_dev_1
cd echo
vim server.go
```
*Note: When edit docker binded files (like "server.go") on the host and in the container at the same time they could be not synchronized.*
