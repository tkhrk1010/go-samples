

### setup shell file
```
chmod +x migrate.sh
```

### init postgresql db in local
```
$ cd db
$ docker-compose up -d
```

check db exist by yourself


### migration
```
$ cd ../migration
```

create migrate file
```
$ make create name=create_hoge_table
```
and write sql

migrate
```
$ make up
```