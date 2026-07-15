# Console Development Commands / Консольные команды разработки.
============================

**Установка утилиты миграции**
brew install golang-migrate

**Установка нативной GO надстройки SQLITE**
go get -u modernc.org/sqlite

**Создание файлов миграции**
migrate create -ext sql -dir ./internal/database/sqlite/schemes -seq init

**Установка надстройки PostgreSQL**
github.com/jmoiron/sqlx

**Устанавливаем GIN - Rest роутер**
go get -u github.com/gin-gonic/gin

**Устанавливаем доступ к переменным окружения**
go get -u github.com/joho/godotenv

**Установка MQTT провайдера**
go get -u github.com/eclipse/paho.mqtt.golang

**Установка Kafka провайдера**
go get -u github.com/twmb/franz-go/pkg/kgo
