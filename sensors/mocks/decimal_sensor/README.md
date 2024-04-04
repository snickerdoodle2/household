# Decimal sensor mock
### Wymagania:
- go `1.22`

Port można zmienić przy użyciu flagi `-port` (domyślnie `8888`). 
Minimalną i maksymalną wartość zwracaną przez sensor ustawia się flagami:
odpoiwednio `-min` oraz `-max`

## Endpointy
1. `GET /status` - zwraca status czujnika
1. `GET /value` - zwraca wartość odczytywaną przez czujnik (na potrzeby "symulacji" czujnika wartość jest obliczana przy wywoływaniu)
