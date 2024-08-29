#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>

#define STASSID "you-wifi-ssid"
#define SENSOR_PIN D7
#define STAPSK "your-wifi-pwd"

const char *ssid = STASSID;
const char *password = STAPSK;

ESP8266WebServer server(80);

void handleStatus()
{
    server.send(200, "text/json", "{ \"status\": \"online\", \"type\": \"binary_sensor\" }");
}

void handleGetValue()
{
    int current_state = digitalRead(SENSOR_PIN);
    server.send(200, "text/json", "{\"value\":" + String(current_state) + "}");
}

void setup(void)
{
    Serial.begin(115200);
    WiFi.mode(WIFI_STA);
    WiFi.begin(ssid, password);
    Serial.println("");

    // Wait for connection
    while (WiFi.status() != WL_CONNECTED)
    {
        delay(500);
        Serial.print(".");
    }
    Serial.println("");
    Serial.print("Connected to ");
    Serial.println(ssid);
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());

    if (MDNS.begin("esp8266"))
    {
        Serial.println("MDNS responder started");
    }

    server.on("/value", HTTP_GET, handleGetValue);
    server.on("/status", HTTP_GET, handleStatus);

    server.begin();
    Serial.println("HTTP server started");
}

void loop(void)
{
    server.handleClient();
    MDNS.update();
}
