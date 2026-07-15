//
//  Settings.swift
//  esp.config
//
//  Created by Dmitriy Soloshenko on 19.06.2026.
//

import SwiftUI

let Settings = SettingsProvider.instance

class SettingsProvider {
    static fileprivate (set) var instance = SettingsProvider()
    
    @AppStorage("wifiExpanded") var wifiExpanded: Bool = false
    @AppStorage("generalExpanded") var generalExpanded: Bool = false
    @AppStorage("mqttExpanded") var mqttExpanded: Bool = false
    @AppStorage("managementExpanded") var managementExpanded: Bool = true

    @AppStorage("wifiSsid") var wifiSsid: String = ""
    @AppStorage("wifiPass")    var wifiPass: String = ""

    @AppStorage("ssdp")    var ssdp: String = ""
    @AppStorage("mdns")    var mdns: String = ""
    @AppStorage("httpPort")    var httpPort: String = "80"

    @AppStorage("mqttHost")     var mqttHost: String = ""
    @AppStorage("mqttPort")     var mqttPort: String = "1883"
    @AppStorage("mqttUser")     var mqttUser: String = ""
    @AppStorage("mqttPass")     var mqttPass: String = ""
    @AppStorage("mqttRoot")     var mqttRoot: String = UUID().uuidString

    
    func clear() {
    }
}

