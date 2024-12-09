#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>

// Replace with your network credentials
#define STASSID "YOUR_SSID"
#define STAPSK "YOUR_PASSWORD"
// Set your sensor type here
const char *sensor_type = "SENSOR_TYPE";

// Set your sensor pins here

ESP8266WebServer server(80);
const char *ssid = STASSID;
const char *password = STAPSK;

float measure()
{
    // Implement your sensor reading here
    return 0.0;
}

void handleGetValue()
{
    float value = measure();
    server.send(200, "text/json", "{\"value\":" + String(value) + "}");
}

void handleToggle()
{
    // Implement your toggle logic here
    return;
}

void handleSetValue()
{
    if(!server.hasArg("value")){
        server.send(400, "text/plain", "400: Invalid Request, no 'value' argument found");
        return;
    }else if (server.arg("value") != "1" && server.arg("value") != "0") {
        server.send(400, "text/plain", "400: Invalid Request, 'value' argument must be 0 or 1");
        return;
    }
    // Implement your setValue logic
    return;
}

void handleStatus()
{
    server.send(200, "text/json", "{ \"status\": \"online\", \"type\": \"" + String(sensor_type) + "\" }");
}

void setup(void)
{
    Serial.begin(115200);
    WiFi.mode(WIFI_STA);
    WiFi.begin(ssid, password);
    Serial.println("");

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
    server.on("/toggle", HTTP_POST, handleToggle);
    server.on("/value", HTTP_PUT, handleSetValue);

    server.begin();
    Serial.println("HTTP server started");
}

void loop(void)
{
    server.handleClient();
    MDNS.update();
}