#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>

#define STASSID "your-wifi-ssid"
#define RELAY_PIN D5
#define STAPSK "your-wifi-password"

const char *ssid = STASSID;
const char *password = STAPSK;

ESP8266WebServer server(80);

void handleToggle()
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
    server.send(200, "text/json", "{\"value\":" + String(current_state) + "}");
}

void handleStatus()
{
    server.send(200, "text/json", "{ \"status\": \"online\", \"type\": \"binary_switch\" }");
}

void handleGetValue()
{
    int current_state = digitalRead(RELAY_PIN);
    server.send(200, "text/json", "{\"value\":" + String(current_state) + "}");
}

void handlePostValue()
{
    Serial.println("handlePostValue called");
    Serial.println(server.arg("value"));

    if (!server.hasArg("value"))
    {
        Serial.println("no 'value' argument");
        server.send(400, "text/plain", "400: Invalid Request, no 'value' argument found");
        return;
    }

    Serial.println(server.arg("value"));

    if (server.arg("value") == "1")
    {
        Serial.println("post value 1");
        digitalWrite(RELAY_PIN, HIGH);
    }
    else if (server.arg("value") == "0")
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
    server.on("/value", HTTP_PUT, handlePostValue);
    server.on("/status", HTTP_GET, handleStatus);
    server.on("/toggle", HTTP_POST, handleToggle);

    server.begin();
    Serial.println("HTTP server started");
    pinMode(RELAY_PIN, OUTPUT);
}

void loop(void)
{
    server.handleClient();
    MDNS.update();
}
