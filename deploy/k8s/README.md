## Kubernetes Infrastructure Stack


[ENG]
# This directory contains manifests and scripts for deploying the IoT infrastructure on Kubernetes.
**🚀 Main Stack (Unified Manifest)**
* All components are consolidated into a single file to simplify dependency management and deployment order:
- infra.yaml — The main manifest. It contains all PVCs, ConfigMaps, Deployments, Services, and Jobs.
**Command to deploy the entire stack:**
Bash
kubectl apply -f infra.yaml

*Important:* Before applying changes, if you have modified the kafka-init-topics-job, you must delete the existing object first:
kubectl delete job kafka-init-topics-job

**📂 Project Structure**
- infra.yaml — The current configuration of the entire infrastructure.
- archive/ — Original YAML files (logical blocks) that served as the basis for consolidation.
- *.sh — Helper scripts for cluster management (startup, port forwarding).

**🛠 Management**
Local Startup: Use en-start-cluster or ru-start-cluster depending on your preference.
Monitoring: Data is collected via Vector, processed and sent to VictoriaMetrics, with visualization provided in Grafana.


[RUS]
# Этот каталог содержит манифесты и скрипты для развертывания IoT-инфраструктуры в Kubernetes.
**🚀 Основной стек (Единый манифест)**
* Все компоненты объединены в один файл для упрощенного управления зависимостями и порядком запуска:
- infra.yaml — основной манифест. Содержит PVC, ConfigMaps, Deployments, Services и Jobs.

**Команда для запуска всего стека:**
Bash
kubectl apply -f infra.yaml

*Важно:* Перед применением изменений, если вы меняли kafka-init-topics-job, удалите старый объект:
kubectl delete job kafka-init-topics-job

**📂 Структура проекта**
- infra.yaml — актуальная конфигурация всей инфраструктуры.
- archive/ — исходные YAML-файлы (логические блоки), послужившие основой для объединения.
- *.sh — вспомогательные скрипты для управления кластером (запуск, проброс портов).

**🛠 Управление**
* Локальный запуск: Используйте en-start-cluster или ru-start-cluster в зависимости от предпочтений.
* Мониторинг: Данные собираются через Vector, обрабатываются и передаются в VictoriaMetrics, визуализация в Grafana.