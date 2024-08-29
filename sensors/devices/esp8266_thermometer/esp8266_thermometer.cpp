#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>
#include <OneWire.h>
#include <DallasTemperature.h>

#define STASSID "YOUR_WIFI_SSID"
#define ONE_WIRE_BUS D6
#define STAPSK "YOUR_WIFI_PASSWORD"

const char *ssid = STASSID;
const char *password = STAPSK;

ESP8266WebServer server(80);
OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);

void handleValue()
{
    sensors.requestTemperatures();
    float tempC = sensors.getTempCByIndex(0);
    server.send(200, "text/json", "{\"value\":" + String(tempC) + "}");
}

void handleStatus()
{
    server.send(200, "text/json", "{ \"status\": \"online\", \"type\": \"decimal_sensor\" }");
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

    server.on("/value", handleValue);
    server.on("/status", handleStatus);

    server.begin();
    Serial.println("HTTP server started");
}

void loop(void)
{
    server.handleClient();
    MDNS.update();
}