## Backend Project with CRUD 

### Tech used
- Go (Gin) for server
- PostgreSQL for Database (using docker)
- SQLC (ORM for using the db)
- Viper (env Variables)
- gomock (to mock the database and test api)



### Setup Database

1) run (and pull) postgres image from docker
```
docker run --name bankapp -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -p 5432:5432 -d postgres
```

--- just for testing 
2) setup the schema in `schema.sql`

3) copy in docker and execute the command
```
docker cp schema.sql bankapp:/schema.sql
```

4) apply the schema
```
psql -U postgres -d postgres -f /schema.sql
```

then check the database schema

```
psql -U postgres -d postgres

```
```
\dt
```
5) database url here is: `postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable`

### Setup database migration

1) Run
```
migrate create -ext sql -dir db/migration -seq init_schema
```

2) Add postgres, createdb, dropdb command to the make file

3) now use migrate to apply migrations and add it to the make file
```
migrate -path db/migration -database "postgres://postgres:postgres@localhost:5432/simple_bank?sslmode=disable" -verbose up
```


### Initialize the project

1) initialize the project
```
go mod init github.com/devrvk/simplebank
```

### gomock testing

1) generate the mock 
```
mockgen -package mockdb -destination db/mock/mock.go github.com/devrvk/simplebank/db/sqlc Store 
```

```
mockgen -package <package name in output> -destination <folder where the mock db is generated> <DB package location > < Interface to mock name>
```

### add user table using migrate

1) generate the table migration versions

```
migrate create -ext sql -dir db/migration -seq add_users
```

2) fill up the up and down scripts

3) update the make file to add one version up and down scripts


## Docker 

1) build image
```
docker build -t simplebank:latest .
```

2) run the image
```
docker run  --name simplebank -p 8080:8080 -e GIN_MODE=release simplebank:latest
```

### Connecting docker containers

-- Way 1
1) run `docker container inspect <name or id>` to see the `IPAddress` under Network settings
2) replace the local host to this IPAddress (take the url as a command line argument while running the backend container)
3) but then we need to change everytime we stop the container

-- Way 2
create a network and add the postgres container to the network

1) run
```bash
raghav.korde@HYDNALUE256711:$ docker network create bank-network
aab2cf957a7e6afcf56b39833d48759e9ecd365ebc52dfef428f89f5592b79f0
raghav.korde@HYDNALUE256711:$ docker network connect bank-network bankapp 
raghav.korde@HYDNALUE256711:$ docker network inspect bank-network 
[
    {
        "Name": "bank-network",
        "Id": "aab2cf957a7e6afcf56b39833d48759e9ecd365ebc52dfef428f89f5592b79f0",
        "Created": "2025-11-07T09:07:28.292717798Z",
        "Scope": "local",
        "Driver": "bridge",
        "EnableIPv4": true,
        "EnableIPv6": false,
        "IPAM": {
            "Driver": "default",
            "Options": {},
            "Config": [
                {
                    "Subnet": "172.19.0.0/16",
                    "Gateway": "172.19.0.1"
                }
            ]
        },
        "Internal": false,
        "Attachable": false,
        "Ingress": false,
        "ConfigFrom": {
            "Network": ""
        },
        "ConfigOnly": false,
        "Containers": {
            "892041fa1db4bd808800fd48606b1ea6aa76707a2663b1ce59cfb396af34aa66": {
                "Name": "bankapp",
                "EndpointID": "2450aaaef8a24f2c16a803d9b32df72d9c50de11c12079cf0bc6e1ecc13f880d",
                "MacAddress": "7e:59:b0:9a:a3:86",
                "IPv4Address": "172.19.0.2/16",
                "IPv6Address": ""
            }
        },
        "Options": {
            "com.docker.network.enable_ipv4": "true",
            "com.docker.network.enable_ipv6": "false"
        },
        "Labels": {}
    }
]
```

2) bankapp (db) is connected to two networks
```bash
docker container inspect bankapp
```
```
...
"Networks": {
                "bank-network": {
                    "IPAMConfig": {},
                    "Links": null,
                    "Aliases": [],
                    "MacAddress": "7e:59:b0:9a:a3:86",
                    "DriverOpts": {},
                    "GwPriority": 0,
                    "NetworkID": "aab2cf957a7e6afcf56b39833d48759e9ecd365ebc52dfef428f89f5592b79f0",
                    "EndpointID": "2450aaaef8a24f2c16a803d9b32df72d9c50de11c12079cf0bc6e1ecc13f880d",
                    "Gateway": "172.19.0.1",
                    "IPAddress": "172.19.0.2",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "DNSNames": [
                        "bankapp",
                        "892041fa1db4"
                    ]
                },
                "bridge": {
                    "IPAMConfig": null,
                    "Links": null,
                    "Aliases": null,
                    "MacAddress": "a2:6b:cf:30:da:bb",
                    "DriverOpts": null,
                    "GwPriority": 0,
                    "NetworkID": "014340612ec7851563efeb494380c612ac6c67c0a37ddd4057e9b6690644bcc7",
                    "EndpointID": "410cd2a2dfd49e16fa55e425bf7e8ecc0772f84cc33b2ae5331e88aaf61dfdb6",
                    "Gateway": "172.17.0.1",
                    "IPAddress": "172.17.0.2",
                    "IPPrefixLen": 16,
                    "IPv6Gateway": "",
                    "GlobalIPv6Address": "",
                    "GlobalIPv6PrefixLen": 0,
                    "DNSNames": null
                }
            }
...
```
3) replace the DB_SOURCE (pass through the command line args) change `localhost` to `bankapp` (container name of DB)

-- Way 3
1) Setup docker compose file (creates a network by default and we can use service name directly in the url)

NOTE:

- Docker compose already runs all the services in one network so no need to create a network
- We need to add wait-for script to migrate the db after the postgres container starts.
- 