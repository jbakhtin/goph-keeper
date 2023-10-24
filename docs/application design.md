# Проектирование и дизайн приложения
## Клиент-серверное взаимодействие
![Взаимодействие Клиент-Сервер](images/client-server-structure-2.png)

## Структура базы данных Серверного приложения
![Структура базы данных](images/db-structure.drawio.png)

## Организация директорий приложения
```
/bin
/cmd
    /client
    /server
/cert
/docs
    /images
/gen/
    /kv/v1
    /creditcrd/v1
    /file/v1
/internal
    /client
    /server
        /application
            /apperror
                /errors.go
                /handler.go
            /config
            /logger
            /type
            app.go
        /composites
            /grpc
        /domain
            /model
            /repository
            /service
        /infastructure
            /database
                /postgres
                    /migration
                    /query
                    /repository
            /objectstorage
            /grpc
                /interseptor
                /handler
                    /kv/v1
                    /creditcard/v1
                    /file/v1
/pkg
    /proto
        /kv/v1
        /creaditcard/v1
        /file/v1
/storage
```