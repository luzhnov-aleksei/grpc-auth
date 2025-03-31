
Для запуска генерации proto файла по новой
protoc --go_out=internal --go-grpc_out=internal internal/protos/auth.proto

### **2.1 Запуск контейнера**

Создайте и запустите контейнер с PostgreSQL:

```
docker run --name postgres-db -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=admin -e POSTGRES_DB=simple_service -p 5432:5432 -d postgres:latest

```
