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
        ZStack {
            Color.black
            VStack {
                StatusSection()
                ScrollView {
                    VStack {
                        WiFiSection()
                        GeneralSection()
                        MqttSection()
                    }
                    .padding()
                }
                .frame(maxWidth: .infinity, maxHeight: .infinity)
                Spacer()
                ManagementSection()
            }
        }
    }
    
    private func StatusSection() -> some View {
        Group {
            VStack {
                if self.vm.connected {
                    VStack {
                        Text("Connected to microcontroller")
                            .foregroundStyle(.white)
                            .padding()
                    }
                    .frame(maxWidth: .infinity)
                    .background(.green)
                }
             
                if self.vm.process {
                    ProgressView()
                        .foregroundColor(.white)
                        .padding()
                }
            }
        }
    }

    
    private func WiFiSection() -> some View {
        DisclosureGroup("Wi-Fi", isExpanded: self.$vm.wifiExpanded) {
            VStack{
                HStack {
                    Text("SSID:")
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter Wi-Fi SSID here", text: self.$vm.wifiSsid)
                        .padding()
                        .textInputAutocapitalization(.never)
                }
                .foregroundColor(.white)
                .frame(maxWidth: .infinity)

                HStack {
                    Text("Pass:")
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter Wi-Fi password here", text: self.$vm.wifiPass)
                        .padding()
                        .textInputAutocapitalization(.never)

                }
                .foregroundColor(.white)
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
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter SSDP name here", text: self.$vm.ssdp)
                        .padding()
                        .textInputAutocapitalization(.never)
                    Button("↺") {
                        self.vm.createSSDPName()
                    }
                    .foregroundColor(.white)
                    .frame(width: 44, height: 44)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)

                HStack {
                    Text("mDNS:")
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter mDNS here", text: self.$vm.mdns)
                        .padding()
                        .textInputAutocapitalization(.never)
                    Button("↺") {
                        self.vm.resetMDNS()
                    }
                    .foregroundColor(.white)
                    .frame(width: 44, height: 44)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)

                HStack {
                    Text("Port:")
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter HTTP port here", text: self.$vm.httpPort)
                        .padding()
                        .textInputAutocapitalization(.never)
                    Button("↺") {
                        self.vm.resetHttpPort()
                    }
                    .foregroundColor(.white)
                    .frame(width: 44, height: 44)
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
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT host here", text: self.$vm.mqttHost)
                        .padding()
                        .textInputAutocapitalization(.never)

                }
                .foregroundColor(.white)
                .frame(maxWidth: .infinity)
                HStack {
                    Text("Port:")
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT port here", text: self.$vm.mqttPort)
                        .padding()
                        .textInputAutocapitalization(.never)
                    Button("↺") {
                        self.vm.resetMqttPort()
                    }
                    .foregroundColor(.white)
                    .frame(width: 44, height: 44)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)
                HStack {
                    Text("User:")
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT user name here", text: self.$vm.mqttUser)
                        .padding()
                        .textInputAutocapitalization(.never)

                }
                .foregroundColor(.white)
                .frame(maxWidth: .infinity)
                HStack {
                    Text("Pass:")
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT user password here", text: self.$vm.mqttPass)
                        .padding()
                        .textInputAutocapitalization(.never)

                }
                .foregroundColor(.white)
                .frame(maxWidth: .infinity)
                HStack {
                    Text("Root:")
                        .frame(width:60, alignment: .trailing)
                    TextField("Enter MQTT root topic here", text: self.$vm.mqttRoot)
                        .padding()
                        .textInputAutocapitalization(.never)
                    Button("↺") {
                        self.vm.createMqttRoot()
                    }
                    .foregroundColor(.white)
                    .frame(width: 44, height: 44)
                    .padding(.leading)
                }
                .frame(maxWidth: .infinity)
            }
            .padding(.vertical)
        }
    }

    private func ManagementSection() -> some View {
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
        .padding()
    }
}

#Preview {
    ContentView()
}
