# ESP8266 thermometer - decimal sensor
Urządzenie zbudowane na bazie ESP8266 z termometrem DS18B20. Termometr jest podłączony do zasilania, uziemienia i pinu D6. Aby odczytać wartości używany jest interfejs OneWire.
### Endpointy
1. `GET /status` - zwraca status urządzenia
1. `GET /value` - zwraca aktualną wartość temperatury w stopniach celsjusza