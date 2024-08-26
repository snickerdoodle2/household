# Binary switch mock
### Wymagania:
- go `1.22`

Port można zmienić przy użyciu flagi `-port` (domyślnie `8888`)

## Endpointy
1. `GET /status` - zwraca status przełącznika
1. `GET /value` - zwraca wartość przełącznika (domyślnie false)
1. `POST /value` - ustawia wartość przełącznika na przeciwną
1. `PUT /value` - ustawia wartość przełącznika na podaną przez użytkownika

# binary_switch.cpp
Przykładowy kod na płytkę ESP8266MOD, do której podpięto przekaźnik SRD-05V-SL-C. Przekaźnik podpięty jest do zasilania 5V z płytki i do pinu nr 5. Po uruchomieniu płytka łączy się z WiFi wg podanego SSID i hasła i na wyjście Serial wypisuje swój adres IP. Domyślnie uruchamia się na porcie 80.

## Endpointy
1. `GET /status` - zwraca status przełącznika
1. `GET /value` - zwraca wartość przełącznika (1/0)
1. `POST /toggle` - ustawia wartość przełącznika na przeciwną
1. `PUT /value` - ustawia wartość przełącznika na podaną przez użytkownika (parametr `value` (1/0))
