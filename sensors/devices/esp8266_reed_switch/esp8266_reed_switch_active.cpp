#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>
#include <ESP8266HTTPClient.h>

#define STASSID "your-wifi-ssid"
#define SENSOR_PIN D7
#define STAPSK "your-wifi-password"

const char *ssid = STASSID;
const char *password = STAPSK;

// test python server ip - to be changed
const char *serverName = "http://10.0.0.55:5000/activesensorupdate";

int previous_sensor_value = 0;

ESP8266WebServer server(80);

void handleStatus()
{
    server.send(200, "text/json", "{ \"status\": \"online\", \"type\": \"binary_sensor_active\" }");
}

void handleGetValue()
{
    int current_state = digitalRead(SENSOR_PIN);
    // Serial.println(current_state);
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
    int current_sensor_value = digitalRead(SENSOR_PIN);

    if (current_sensor_value != previous_sensor_value)
    {
        WiFiClient client;
        HTTPClient http;

        // Serial.println("value changed from " + String(previous_sensor_value) + " to " + String(current_sensor_value));

        http.begin(client, serverName);
        http.addHeader("Content-Type", "application/json");
        int httpResponseCode = http.POST("{\"sensor_id\":\"abc123\",\"value\":\"" + String(current_sensor_value) + "\"}");

        // Serial.println("HTTP Response code: ");
        // Serial.println(httpResponseCode);

        http.end();
        previous_sensor_value = current_sensor_value;
    }
    server.handleClient();
    MDNS.update();
}
