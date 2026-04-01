### swagger

```console
# 新增接口，补完注释后执行
swag init -g main.go -o docs
```

### docker

启动数据库：

```
docker compose up -d pg
```

执行迁移：

```console
make migrate
```

首次迁移会插入一个默认用户：

- `id=1`
- `username=admin`
- `email=admin@example.com`

启动应用：

```console
docker compose up --build app
```

如果你希望显式指定命令，也可以直接执行：

```console
go run . migrate -d configs -c dev
```

新增表结构时，追加 migration 文件即可：

- `internal/db/migrations/000002_xxx.up.sql`
- `internal/db/migrations/000002_xxx.down.sql`
