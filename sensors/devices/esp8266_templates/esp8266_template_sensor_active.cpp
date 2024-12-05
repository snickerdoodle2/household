#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>
#include <ESP8266HTTPClient.h>
#include <ArduinoJson.h>

// Replace with your network credentials
#define STASSID "YOUR_SSID"
#define STAPSK "YOUR_PASSWORD"

// Set your sensor type here
const char *sensor_type = "SENSOR_TYPE";

// Set your sensor pins here

ESP8266WebServer server(80);
const char *ssid = STASSID;
const char *password = STAPSK;

struct SensorConfig
{
    String serverUri;
    String initAckEndpoint;
    String measurementsEndpoint;
    String idToken;
    bool isConfigured = false;
};

SensorConfig sensorConfig;

float measure()
{
    // Implement your sensor reading here
    return 0.0;
}

bool sendingCondition()
{
    // Implement your sending condition here
    return true;
}

void handleValue()
{
    float value = measure();
    server.send(200, "text/json", "{\"value\":" + String(value) + "}");
}

void handleStatus()
{
    server.send(200, "text/json", "{ \"status\": \"online\", \"type\": \"" + String(sensor_type) + "\" }");
}

void handleInit()
{
    if (server.method() != HTTP_POST)
    {
        server.send(405, "text/plain", "Method Not Allowed");
        return;
    }

    JsonDocument requestBody;
    DeserializationError error = deserializeJson(requestBody, server.arg("plain"));

    if (error)
    {
        server.send(400, "text/plain", "Invalid JSON");
        return;
    }

    String serverUri = requestBody["server-uri"].as<String>();
    if (serverUri.startsWith("http://") != 0 && serverUri.startsWith("https://") != 0)
    {
        serverUri = "http://" + serverUri;
    }
    sensorConfig.serverUri = serverUri;
    sensorConfig.idToken = requestBody["id-token"].as<String>();
    sensorConfig.measurementsEndpoint = requestBody["measurements-endpoint"].as<String>();
    sensorConfig.initAckEndpoint = requestBody["init-ack-endpoint"].as<String>();
    sensorConfig.isConfigured = true;

    server.send(200, "text/plain", "Initialization successful");

    sendInitAck();
}

void sendInitAck()
{
    JsonDocument requestBody;
    requestBody["id-token"] = sensorConfig.idToken;
    requestBody["server-uri"] = sensorConfig.serverUri;
    requestBody["init-ack-endpoint"] = sensorConfig.initAckEndpoint;
    requestBody["measurements-endpoint"] = sensorConfig.measurementsEndpoint;

    String jsonString;
    serializeJson(requestBody, jsonString);

    WiFiClient client;
    HTTPClient http;

    http.begin(client, sensorConfig.serverUri + sensorConfig.initAckEndpoint);
    http.addHeader("Content-Type", "application/json");
    int httpResponseCode = http.POST(jsonString);

    http.end();
}

void sendValue()
{
    JsonDocument requestBody;
    requestBody["message-type"] = "measurement";
    requestBody["sensor-type"] = "decimal_sensor";
    requestBody["value"] = measure();
    requestBody["id-token"] = sensorConfig.idToken;

    String jsonString;
    serializeJson(requestBody, jsonString);

    WiFiClient client;
    HTTPClient http;

    http.begin(client, sensorConfig.serverUri + sensorConfig.measurementsEndpoint);
    http.addHeader("Content-Type", "application/json");
    int httpResponseCode = http.POST(jsonString);

    http.end();
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

    server.on("/status", HTTP_GET, handleStatus);
    server.on("/value", HTTP_GET, handleValue);
    server.on("/init", HTTP_POST, handleInit);

    server.begin();
    Serial.println("HTTP server started");
}

void loop(void)
{
    server.handleClient();

    if (sensorConfig.isConfigured && sendingCondition())
    {
        sendValue();
    }

    MDNS.update();
}