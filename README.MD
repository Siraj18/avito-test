# BALANCE SERVICE

Тестовое задание

### Инструкция по запуску

#### 1. Билд контейнеров
```
docker-compose build
```

#### 2. После того как образы собраны, запускаем контейнеры
```
docker-compose up -d
```


#### 3. Для приостановки контейнеров используем следующую команду
```
docker-compose down
```


### Запросы и примеры ответов
Все запросы собраны в коллекции Postman. Доступ к коллекции можно получить по [ссылке](https://www.getpostman.com/collections/7e29a2b327292090f9b5 "Ссылка на postman коллекцию")   
Если у вас отсутствует Postman, то далее расписаны curl запросы.


#### Запрос на изменение баланса
Данный запрос меняет баланс пользователя по его uuid. Если поле "money" положительное, то баланс увеличивается, если отрицательное, то уменьшается.
Изначально сервис не содержит в себе информацию о пользователях. Пользователь добавляется в базу при первом зачислении денег на счет.
```
$ curl --location --request POST 'localhost:8000/changeBalance' \
    --header 'Content-Type: application/json' \
    --data-raw '{"id": "34be95d0-9a41-11ec-b909-0242ac120002", "money": 50}'
```  
#### Запрос на получение баланса
Данный запрос получает баланс пользователя по его uuid. Баланс пользователя хранится в рублях. При желании можно прописать параметр currency с указанием   
необходимой валюты.
```
$ curl --location --request POST 'localhost:8000/balance?currency=EUR' \
    --header 'Content-Type: application/json' \
    --data-raw '{"id": "34be95d0-9a41-11ec-b909-0242ac120002"}'
```

#### Запрос на перевод средств
Данный запрос приминает в себя uuid пользователей и также необходимую сумму.
```
curl --location --request POST 'localhost:8000/transferBalance' \
    --header 'Content-Type: application/json' \
    --data-raw '{"from_id": "34be95d0-9a41-11ec-b909-0242ac120002", "to_id": "34be95d0-9a41-11ec-b909-0242ac120003", "money": 50}'
```

#### Запрос на получение списка транзакций
Данный запрос принимает на вход uuid пользователя. На вывод он дает все транзакции, в которых замешан данные uuid. Также предусмотрена фильтрация и пагинация   
данных. Фильтры: "date_asc", "date_desc", "money_asc", "money_desc".
```
$ curl --location --request POST 'localhost:8000/allTransactions' \
    --header 'Content-Type: application/json' \
    --data-raw '{"id": "34be95d0-9a41-11ec-b909-0242ac120002", "sort_type":"date_asc", "limit": 10, "page":1}'
```
### Запуск тестов

#### Команда для запуска тестов
```
go test ./... -cover
```
