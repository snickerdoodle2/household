#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>
#include <ESP8266HTTPClient.h>

#define STASSID "your-wifi-ssid"
#define STAPSK "your-wifi-password"

const char *ssid = STASSID;
const char *password = STAPSK;

// hardcoded for now
const char *serverName = "http://10.0.0.55:5000/activesensorupdate";

#define SENSOR_PIN D7
#define RELAY_PIN D5

int previous_sensor_value = 0;

ESP8266WebServer reed_server(80);
ESP8266WebServer relay_server(8008);

void handleReedStatus()
{
    reed_server.send(200, "text/json", "{ \"status\": \"online\", \"type\": \"binary_sensor_active\" }");
}

void handleReedGetValue()
{
    int current_state = digitalRead(SENSOR_PIN);
    reed_server.send(200, "text/json", "{\"value\":" + String(current_state) + "}");
}

void handleRelayToggle()
{
    int current_state = digitalRead(RELAY_PIN);

    if (current_state == HIGH)
    {
        digitalWrite(RELAY_PIN, LOW);
    }
    else
    {
        digitalWrite(RELAY_PIN, HIGH);
    }

    current_state = digitalRead(RELAY_PIN);
    relay_server.send(200, "text/json", "{\"value\":" + String(current_state) + "}");
}

void handleRelayStatus()
{
    relay_server.send(200, "text/json", "{ \"status\": \"online\", \"type\": \"binary_switch\" }");
}

void handleRelayGetValue()
{
    int current_state = digitalRead(RELAY_PIN);
    relay_server.send(200, "text/json", "{\"value\":" + String(current_state) + "}");
}

void handleRelayPostValue()
{
    String body = server.arg("plain");
    Serial.println("Received body: " + body);

    JsonDocument doc;
    DeserializationError error = deserializeJson(doc, body);

    if (error)
    {
        Serial.println("Failed to parse JSON: " + String(error.c_str()));
        server.send(400, "text/plain", "400: Invalid Request, failed to parse JSON");
        return;
    }

    if (!doc["value"].is<int>())
    {
        Serial.println("No 'value' key in JSON or 'value' is not an integer");
        server.send(400, "text/plain", "400: Invalid Request, no 'value' key in JSON or invalid type");
        return;
    }

    int value = doc["value"].as<int>();
    Serial.println("Parsed value: " + String(value));

    if (value == 1)
    {
        Serial.println("post value 1");
        digitalWrite(RELAY_PIN, HIGH);
    }
    else if (value == 0)
    {
        Serial.println("post value 0");
        digitalWrite(RELAY_PIN, LOW);
    }
    else
    {
        Serial.println("post value incorrect");
        server.send(400, "text/plain", "400: Invalid Request, 'value' argument incorrect [0/1]");
        return;
    }

    int current_state = digitalRead(RELAY_PIN);
    Serial.println("pin set to " + String(current_state));
    server.send(200, "application/json", "{\"value\": " + String(current_state) + "}");
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

    reed_server.on("/value", HTTP_GET, handleReedGetValue);
    reed_server.on("/status", HTTP_GET, handleReedStatus);

    relay_server.on("/value", HTTP_GET, handleRelayGetValue);
    relay_server.on("/value", HTTP_PUT, handleRelayPostValue);
    relay_server.on("/status", HTTP_GET, handleRelayStatus);
    relay_server.on("/toggle", HTTP_POST, handleRelayToggle);

    reed_server.begin();
    Serial.println("reed_HTTP server started");
    relay_server.begin();
    Serial.println("relay_HTTP server started");
    pinMode(RELAY_PIN, OUTPUT);
}

void loop(void)
{
    int current_sensor_value = digitalRead(SENSOR_PIN);

    if (current_sensor_value != previous_sensor_value)
    {
        WiFiClient client;
        HTTPClient http;

        http.begin(client, serverName);
        http.addHeader("Content-Type", "application/json");
        int httpResponseCode = http.POST("{\"sensor_id\":\"abc123\",\"value\":\"" + String(current_sensor_value) + "\"}");

        http.end();
        previous_sensor_value = current_sensor_value;
    }
    reed_server.handleClient();
    relay_server.handleClient();
    MDNS.update();
}
