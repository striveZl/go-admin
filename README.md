### docker

启动应用：

```console
docker compose up --build app
```
执行迁移：

```console
make migrate
```
### swagger

```console
# 新增接口，补完注释后执行
swag init -g main.go -o docs
```
