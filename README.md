## Локальный запуск

Для локального запуска требуется Docker и Docker-compose  
Команда
```
docker-compose up -d
```


Для запуска генерации proto файла по новой
```
protoc --go_out=internal --go-grpc_out=internal internal/protos/auth.proto
```

