//
//  ContentViewModel.swift
//  esp.config
//
//  Created by Dmitriy Soloshenko on 19.06.2026.
//
import Combine
import SwiftUI
import NetworkExtension

@Observable
class ContentVM {
    
    // MARK: - Types
    
    // MARK: - Publishers
    
    // MARK: - Public properties
    var wifiExpanded       = Settings.wifiExpanded
    var generalExpanded    = Settings.generalExpanded
    var mqttExpanded       = Settings.mqttExpanded
    var managementExpanded = Settings.managementExpanded
    
    // Wi-Fi section
    var wifiSsid = Settings.wifiSsid
    var wifiPass = Settings.wifiPass

    // General section
    var ssdp     = Settings.ssdp
    var mdns     = Settings.mdns
    var httpPort = Settings.httpPort

    // MQTT section
    var mqttHost = Settings.mqttHost
    var mqttPort = Settings.mqttPort
    var mqttUser = Settings.mqttUser
    var mqttPass = Settings.mqttPass
    var mqttRoot = Settings.mqttRoot

    
    var connected = false

    // MARK: - Private properties
    private var ssidAP = "esp8266-setup"
    private var passAP = "1qazxsw2"
    private var hostAP = "http://192.168.4"
    private var portAP = "80"
    
    private var configUrl: String {
        return "\(self.hostAP):\(self.portAP)/config"
    }
    
    // MARK: - Init
    deinit {
        self.disconnectFromWifi(ssid: self.ssidAP)
    }
    
    // MARK: - Public methods
    func createMqttRoot() {
        self.mqttRoot = UUID().uuidString
    }

    func resetMqttPort() {
        self.mqttPort = "1883"
    }

    func resetHttpPort() {
        self.mqttPort = "80"
    }

    func resetMDNS() {
        self.mdns = "local.esp8266"
    }
    
    func createSSDPName() {
        self.ssdp = "esp-" + "\(Int(Date().timeIntervalSince1970))"
    }
    
    func save() {
        Settings.wifiExpanded = self.wifiExpanded
        Settings.generalExpanded = self.generalExpanded
        Settings.mqttExpanded = self.mqttExpanded
        Settings.managementExpanded = self.managementExpanded
        
        Settings.wifiSsid = self.wifiSsid
        Settings.wifiPass = self.wifiPass

        Settings.ssdp = self.ssdp
        Settings.mdns = self.mdns
        Settings.httpPort = self.httpPort


        Settings.mqttHost = self.mqttHost
        Settings.mqttPort = self.mqttPort
        Settings.mqttUser = self.mqttUser
        Settings.mqttPass = self.mqttPass
        Settings.mqttRoot = self.mqttRoot
    }
    
    func reconnection() {
        if self.connected {
            self.disconnectFromWifi(ssid: self.ssidAP)
        }
        self.connectToWifi(ssid: self.ssidAP, password: self.passAP)
    }
    
    func send() {
        let jsonBody: [String: Any] = [
            "WIFI_SSID": self.wifiSsid,
            "WIFI_PASS": self.wifiPass,
            "MDNS": self.mdns,
            "SSDP": self.ssdp,
            "HTTP_PORT": self.httpPort,
            "MQTT_HOST": self.mqttHost,
            "MQTT_PORT": self.mqttPort,
            "MQTT_USER": self.mqttUser,
            "MQTT_PASS": self.mqttPass,
            "MQTT_ROOT": self.mqttRoot
        ]
        
        self.sendDataToEsp(jsonBody)
    }
    
    // MARK: - Private methods
    private func connectToWifi(ssid: String, password: String) {
        let hotspotConfig = NEHotspotConfiguration(ssid: ssid, passphrase: password, isWEP: false)
        
        // IMPORTANT FOR IoT: the network keeps only until the application is active on the screen.
            hotspotConfig.joinOnce = true

        NEHotspotConfigurationManager.shared.apply(hotspotConfig) { (error) in
            if let error = error {
                let nsError = error as NSError
                switch nsError.code {
                case NEHotspotConfigurationError.userDenied.rawValue:
                    print("User cancelled connection request.")
                case NEHotspotConfigurationError.invalidSSID.rawValue:
                    print("SSID is wrong.")
                default:
                    print("Connection error to Wi-Fi: \(error.localizedDescription)")
                }
            } else {
                print("Device connected successfully to \(ssid)")
                self.connected = true
            }
        }
    }

    
    private func disconnectFromWifi(ssid: String) {
        NEHotspotConfigurationManager.shared.removeConfiguration(forSSID: ssid)
        print("The network configuration \(ssid) was removed. Device comes back to previous network.")
        self.connected = false
    }

    
    private func sendDataToEsp(_ jsonBody: [String: Any]) {
        guard let url = URL(string: self.configUrl) else { return }
                
        guard let jsonData = try? JSONSerialization.data(withJSONObject: jsonBody, options: []) else {
            print("Error: Cannot create JSON")
            return
        }
        
        var request = URLRequest(url: url)
        request.httpMethod = "POST"
        request.setValue("application/json", forHTTPHeaderField: "Content-Type") // Говорим ESP, что шлем JSON
        request.httpBody = jsonData
        
        request.timeoutInterval = 10.0
        
        let task = URLSession.shared.dataTask(with: request) { data, response, error in
            if let error = error {
                print("Network error: \(error.localizedDescription)")
                return
            }
            
            if let httpResponse = response as? HTTPURLResponse {
                print("Check status code from ESP: \(httpResponse.statusCode)")
                
                if httpResponse.statusCode == 200 {
                    print("Data transferred TO esp8266 successfully!")
                    
                    // RUS:
                    // Если передача данных успешная, то микроконтроллер сам потушит точку доступа,
                    // а "телефон" автоматически переключится на прежнюю Wi-Fi сеть.

                    // ENG:
                    // If the data transfer is successful, the microcontroller
                    // automatically turns off the access point.
                    // The "phone" automatically switches to the previous Wi-Fi network.
                    
                }
            }
            
            if let data = data, let responseString = String(data: data, encoding: .utf8) {
                print("Microcontroller response: \(responseString)")
            }
        }
        
        task.resume()
    }

}
