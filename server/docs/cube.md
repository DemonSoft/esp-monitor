# Kubernetes (Minikube) command line
-------------------

## Установка Minikube

**Установка Minikube**
brew install minikube *не работает!* 
Вместо этого следует использовать:
* 1. Переходим в домашнюю директорию пользователя
cd ~

* 2. Скачиваем бинарник (теперь файл точно запишется)
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-darwin-amd64

* 3. Устанавливаем его в систему (потребуется пароль от Макбука)
sudo install minikube-darwin-amd64 /usr/local/bin/minikube

* 4. Удаляем временный файл из домашней папки
rm minikube-darwin-amd64

**Установка Kuberctl**
brew install kubernetes-cli

**Проверка установки Kuberctl**
kubectl version --client

**Старт кластера**
minikube start --driver=docker

**Проверка запуска кластера**
kubectl get nodes

**Статус системных подов**
kubectl get pods -n kube-system

**Установка mosquitto**
kubectl apply -f mosquitto.yaml *НЕ РАБОТАЕТ*
minikube image load eclipse-mosquitto:latest

**Форвардинг mosquitto портов**
kubectl port-forward service/mosquitto-service 1883:1883 --address 0.0.0.0
*команда должна работать постоянно в отдельном окне терминала*

**Установка Kafka и Kafdrop из под minikube**
minikube image load apache/kafka:3.7.0
minikube image load obsidiandynamics/kafdrop:latest

**Запуск Kafka  манифеста**
kubectl apply -f kafka-stack.yaml

**Список загруженных образов**
minikube image ls --format=table



## Minikube building
**Сборка сервера в minikube**
*Переключение контекста Докера на minicube*
eval $(minikube docker-env)
*Сбока сервера*
docker build -t docker.io/library/esp-server:latest .
*Сбока воркеров*
docker build -t docker.io/library/kworker:latest .
docker build -t docker.io/library/ai-worker:latest .
docker build -t docker.io/library/hworker:latest .

*Применение конфигурации*
kubectl apply -f apps-stack.yaml

**Пробрасывание портов между локальной машиной и minikube**
sudo minikube tunnel

**Чтение последних 50 строк лог файла**
kubectl logs deployment/esp-server-deployment -c esp-worker --tail=50

**Перезапуск бэкенда в кубе**
kubectl rollout restart deployment/esp-server-deployment

**Проверка доступных сервисов и портов**
kubectl get svc

## Minikube Kafka
**Топики**
*Удаление задачи на кубе перед новым созданием*
kubectl delete job kafka-create-topics-job
*Создание новых топиков* через kubectl apply -f apps-stack.yaml

**Запуск Kafdrop через туннель**
*Kafdrop автоматически откроется*
minikube service kafdrop-service

## Minikube Others
**Показывает владельца порта**
sudo lsof -i :1883

**Убивает процесс, занявший порт**
*<PID> следует извлечь из предыдущей команды*
sudo kill -9 <PID>

**Добавляем права на старт кластера**
chmod +x start-cluster.sh

**Перекидывание образа между окружениями**
minikube image load esp-worker:v3

**Graphana**
Ссылка для связывания с Prometeus:
http://victoria-metrics-service:8428

**Непрерывное чтение логов (воркера) в консоли**
kubectl logs -f deployment/esp-worker-deployment -c esp-worker

**Задание секретного ключа в кластере**
kubectl create secret generic ai-worker-secrets \
  --from-literal=OPENROUTER_API_KEY="<OPENROUTER_API_KEY>"
