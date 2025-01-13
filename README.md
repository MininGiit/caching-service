# caching-service

# Сервис обработки заказов
Cервис с простейшим интерфейсом для in-memoru кеширования. В сервисе реализован алгоритм кеширования LRU c TTL
Сервис предоставляет возможность получить данные, загрузить данные и удалить их.

## Требования
* Версия Go 1.22.2
* Docker

## Конфигурация
Сервис поддерживает следующие параметры конфигурации:
    - `server_host_port`: конфигурируется через переменную окружения `SERVER_HOST_PORT` или через флаг `-server-host-port`, по умолчанию равен `localhost:8080`
    - `cache_size`: конфигурируется через переменную окружения `CACHE_SIZE` или флаг `-cache-size`, по умолчанию равен `10`
    - `default_cache_ttl`: конфигурируется через переменную окружения `DEFAULT_CACHE_TTL` или флаг `-default-cache-ttl`, по умолчанию равен `1 минуте`
    - `log_level`: конфигурируется через переменную окружения `LOG_LEVEL` или флаг `-log-level`, по умолчанию выставляется в `WARN`

## Тестирование 
Для запуска тестов с расчётом покрытия используется скрипт:
    ```
    go test -v -coverpkg=./... -coverprofile=coverage.out -covermode=count ./... && go tool cover -func coverage.out | grep total | awk '{print $3}'
    ```
    
## Инструкция по запуску 

1) Сборка
```
docker build -t cache-service .
```
2) Запуск сервиса
```
docker run -e SERVER_HOST_PORT="0.0.0.0:8080" -e CACHE_SIZE="200" -e DEFAULT_CACHE_TTL="120" -e LOG_LEVEL="DEBUG" -p 8080:8080 cache-service
```

## Использование сервиса

#### Обработчик `POST /api/lru`

Добавление данных в кэш. 

Примеры запросов:
1. Запрос с указанием TTL в секундах и строковым значением
    ```
    POST /api/lru HTTP/1.1
    Content-Type: application/json
    ...
    
    {
        "key": "some_key",
        "value": "some_cached_value",
        "ttl_seconds": 30
    }
    ```
2. Запрос без указания TTL и с числовым значением
    ```
    POST /api/lru HTTP/1.1
    Content-Type: application/json
    ...
    
    {
        "key": "some_key",
        "value": 331.1
    }
    ```
   
Возможные коды ответа:
1. `201` - данные записаны в кэш
2. `400` - невалидные входные данные

#### Обработчик `GET /api/lru/{key}`

Получение данных по ключу

Пример запроса:

```
GET /api/lru/some_key HTTP/1.1
Content-Length: 0
```

Возможные ответы сервера: 
1. `200` - Данные получены
   ```
   200 OK HTTP/1.1
   Content-Type: application/json
   ...

    {
        "key": "some_key",
        "value": 311.1,
        "expires_at": 1718278493
    }
   ```
2. `404` - ключ не найден

#### Обработчик `GET /api/lru`

Получение всего текущего наполнения кэша в виде двух списков: списка ключей и списка значений. Пары ключ-значение располагаются на соответствующих индексах. Других ограничений на порядок нет.

Возможные ответы сервера:
1. `200`
    ```
   200 OK HTTP/1.1
   Content-Type: application/json
   ...

    {
        "keys": ["some_key3", "some_key1", "some_key2"],
        "values": ["some_value3", 1.1, "some_value2"]
    }
   ```
2. `204` - кэш пустой

#### Обработчик `DELETE /api/lru/{key}`

Удаление пары ключ/значение 

Пример запроса:

```
DELETE /api/lru/some_key HTTP/1.1
Content-Length: 0
```

Возможные ответы сервера:
1. `204` - успешное удаление
2. `404` - ключ не найден

#### Обработчик `DELETE /api/lru`

Полная очистка кэша. 

Пример запроса:

```
DELETE /api/lru HTTP/1.1
Content-Length: 0
```

Возможные ответы сервера:
1. `204` - успешное удаление
```
