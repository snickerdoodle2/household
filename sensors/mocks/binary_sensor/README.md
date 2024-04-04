# Binary sensor mock
### Wymagania:
- go `1.22`

Port można zmienić przy użyciu flagi `-port` (domyślnie `8888`)

## Endpointy
1. `GET /status` - zwraca status czujnika
1. `GET /value` - zwraca wartość odczytywaną przez czujnik (na potrzeby "symulacji" czujnika wartość jest obliczana przy wywoływaniu)
