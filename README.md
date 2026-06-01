### docker

启动应用：

```console
make app-build
```

执行迁移：

```console
make migrate-docker
```

### swagger

```console
# 新增接口，补完注释后执行
swag init -g main.go -o docs
```
