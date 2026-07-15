## Описание проекта
-------------------

1. При включении микроконтроллеера, если в файловой системе отсуствует файл config.json, микроконтроллер поднимает точку доступа c ssid = "esp8266-setup" и паролем "1qazxsw2"; запускает веб-сервер для обработки REST запроса по адресу 192.168.4.1 на 80-м порту. 
2. После получения  REST запроса POST /config вычитывает полученные JSON параметры:
    WIFI_SSID - идентифкатор wifi  сети
    WIFI_PASS - пароль wifi  сети
    SSDP      - ssdp name объекта ssdp
    MDNS      - mdns
    HTTP_PORT - порт, на котором будет отвечать web-сервер
    MQTT_HOST - хост для подключения к MQTT брокеру.
    MQTT_PORT - порт подключения к MQTT брокеру.
    MQTT_USER - имя пользователя MQTT брокера.
    MQTT_PASS - пароль пользователя MQTT.
    MQTT_ROOT - корневой топик MQTT.

3. Из полученных данных формируется JSON объект:
{
    "wifi": {
        "ssid": "<WIFI_SSID>",
        "pass": "<WIFI_PASS>"
    },
    "ssdp" : {
        "name": "<SSDP>",
        "url": "/",
        "serial": "",
        "model": {
            "name":"",
            "number":"",
            "url":""
        },
        "manufacturer": {
            "name":"",
            "url":""
        },
        "device_type": "upnp:rootdevice",
        "schema_url": "description.xml",
        "http_port": "<HTTP_PORT>"
    },
    "mdns": "<MDNS>",
    "mqtt": {
        "host":"<MQTT_HOST>",
        "port":"<MQTT_PORT>",
        "user":"<MQTT_USER>",
        "pass":"<MQTT_PASS>",
        "root":"<MQTT_ROOT>"
    }
}
 Если в полученном запросе поля "SSDP" и "MDNS" не заданы, то для <MDNS> используется значение "local.esp8266", а вместо <SSDP> создается название из префикса "esp8266-" и текущего даты-времени в формате unixtime (количество секунд, прошедших с 1 января 1970-го года). Пример для <SSDP> - "esp8266-1781289531". Если параметр <HTTP_PORT> не задан, то используется порт 80. Если для <MQTT_PORT> порт не задан, то используется 1883. 
4. Данный объект сохраняется в JSON формате, в файловой системе микроконтроллера под именем config.json.
5. Точка доступа закрывается, а микроконтрроллер перезагружается.
6. Если при загрузке микроконтроллера присуствует файл config.json, то он загружается в память, для дальнейшего использования содержащихся в нем значений.
7. Микроконтроллер подключается к wifi сети, используя параметры <wifi.ssid> и <wifi.pass> из файла конфигурации.
8. Микроконтроллер получает unixtime от NTP сервера, например "pool.ntp.org" и сохраняет его в паременной started.
9. Микроконтроллер запускает веб-сервер на порту указанном в <ssdp.http_port>.
10. Микроконтроллер запускает SSDP сервис, с параметрами из config.json.
Пример:
    SSDP.setDeviceType("upnp:rootdevice");
    SSDP.setSchemaURL("description.xml");
    SSDP.setHTTPPort(<ssdp.http_port>);
    SSDP.setName(<ssdp.name>);
    SSDP.setURL(<ssdp.url>);
    SSDP.setSerialNumber(<ssdp.serial>);
    SSDP.setModelName(<ssdp.model.name>);
    SSDP.setModelNumber(<ssdp.model_number>);
    SSDP.setModelURL(<ssdp.model.url>);
    SSDP.setManufacturer(<ssdp.manufacturer.name>);
    SSDP.setManufacturerURL(<ssdp.manufacturer.url>);
    SSDP.begin();
11. Микроконтроллер подключется к брокеру MQTT, используя параметры из файла config.json и подписываясь на топик вида <mqtt.root>/<ssdp.name>#
12. Микроконтроллер отправляет сообщение в топик <mqtt.root>/<ssdp.name>/state  со следующим содержимым:
{
  "SSDP": "<ssdp.name>",
  "MDNS": "<mdns>",
  "Started": <started>,
  "Updated": <текущее значение unixtime>,
  "Pins": {
    <значение pin микроконтроллера>
    "D1": "...",
    "D2": "...",
    "D3": "...",
    ...
  }
}
Такое сообщение должно отправляться каждый раз при изменении значение какого-нибудь pin микроконтролера, но не чаще чем раз в 1 секунду. Максимальная частота должна быть задана переменной внутри исходного кода.
13. Если в топик <mqtt.root>/<ssdp.name>/action приходит объект
{
    "config":{...}
}
то все значения из этого объекта применяются к соотвествующим полям объекта config, сам объект config сохраняется в файловой системе в виде файла config.json, а микроконтроллер перегружается.
Необходимо предсмотреть возможность получения и других объектов данных из топика mqtt.root>/<ssdp.name>/action, но, пока ничего не делать.

14. Если распознается нажатие кнопки более трех секунд, то файл config.json удаляется, а микроконтроллер перезагружается.