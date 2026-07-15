//
//  ContentView.swift
//  esp.config
//
//  Created by Dmitriy Soloshenko on 19.06.2026.
//

import SwiftUI

struct ContentView: View {
    
    @State var vm = ContentVM()

    var body: some View {
        VStack {
            StatusSection()
            ScrollView {
                VStack {
                    WiFiSection()
                    GeneralSection()
                    MqttSection()
                    ManagementSection()
                }
                .padding()
            }
            .frame(maxWidth: .infinity, maxHeight: .infinity)
        }
    }
    
    private func StatusSection() -> some View {
        Group {
            if self.vm.connected {
                VStack {
                    Text("Connected to microcontroller")
                        .foregroundStyle(.white)
                        .padding()
                }
                .frame(maxWidth: .infinity)
                .background(.green)
            }
        }
    }

    
    private func WiFiSection() -> some View {
        DisclosureGroup("Wi-Fi", isExpanded: self.$vm.wifiExpanded) {
            VStack{
                HStack {
                    Text("SSID:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter Wi-Fi SSID here", text: self.$vm.wifiSsid)
                        .padding()
                        .foregroundColor(.black)

                }
                .frame(maxWidth: .infinity)

                HStack {
                    Text("Pass:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter Wi-Fi password here", text: self.$vm.wifiPass)
                        .padding()
                        .foregroundColor(.black)

                }
                .frame(maxWidth: .infinity)
            }
            .padding(.vertical)
        }
    }

    private func GeneralSection() -> some View {
        DisclosureGroup("General", isExpanded: self.$vm.generalExpanded) {
            VStack{
                HStack {
                    Text("SSDP:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter SSDP name here", text: self.$vm.ssdp)
                        .padding()
                        .foregroundColor(.black)
                    Button("↺") {
                        self.vm.createSSDPName()
                    }
                    .frame(width: 44, height: 44)
                    .foregroundColor(.blue)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)

                HStack {
                    Text("mDNS:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter mDNS here", text: self.$vm.mdns)
                        .padding()
                        .foregroundColor(.black)
                    Button("↺") {
                        self.vm.resetMDNS()
                    }
                    .frame(width: 44, height: 44)
                    .foregroundColor(.blue)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)

                HStack {
                    Text("Port:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter HTTP port here", text: self.$vm.httpPort)
                        .padding()
                        .foregroundColor(.black)
                    Button("↺") {
                        self.vm.resetHttpPort()
                    }
                    .frame(width: 44, height: 44)
                    .foregroundColor(.blue)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)
            }
            .padding(.vertical)
        }
    }

    private func MqttSection() -> some View {
        DisclosureGroup("Mqtt", isExpanded: self.$vm.mqttExpanded) {
            VStack {
                HStack {
                    Text("Host:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT host here", text: self.$vm.mqttHost)
                        .padding()
                        .foregroundColor(.black)

                }
                .frame(maxWidth: .infinity)
                HStack {
                    Text("Port:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT port here", text: self.$vm.mqttPort)
                        .padding()
                        .foregroundColor(.black)
                    Button("↺") {
                        self.vm.resetMqttPort()
                    }
                    .frame(width: 44, height: 44)
                    .foregroundColor(.blue)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)
                HStack {
                    Text("User:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT user name here", text: self.$vm.mqttUser)
                        .padding()
                        .foregroundColor(.black)

                }
                .frame(maxWidth: .infinity)
                HStack {
                    Text("Pass:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT user password here", text: self.$vm.mqttPass)
                        .padding()
                        .foregroundColor(.black)

                }
                .frame(maxWidth: .infinity)
                HStack {
                    Text("Root:")
                        .foregroundColor(.black)
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT root topic here", text: self.$vm.mqttRoot)
                        .padding()
                        .foregroundColor(.black)
                    Button("↺") {
                        self.vm.createMqttRoot()
                    }
                    .frame(width: 44, height: 44)
                    .foregroundColor(.blue)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)
            }
            .padding(.vertical)
        }
    }

    private func ManagementSection() -> some View {
        DisclosureGroup("Management", isExpanded: self.$vm.managementExpanded) {
            VStack {
                Button("Save configuration") {
                    self.vm.save()
                }
                .padding()
                .frame(maxWidth: .infinity)
                .foregroundColor(.white)
                .background(.indigo)
                .cornerRadius(8)

                Button("Reconnection") {
                    self.vm.reconnection()
                }
                .padding()
                .frame(maxWidth: .infinity)
                .foregroundColor(.white)
                .background(.indigo)
                .cornerRadius(8)

                Button("Send") {
                    self.vm.send()
                }
                .padding()
                .frame(maxWidth: .infinity)
                .foregroundColor(.white)
                .background(.indigo)
                .cornerRadius(8)
            }
            .padding(.vertical)
        }
    }

}

#Preview {
    ContentView()
}
