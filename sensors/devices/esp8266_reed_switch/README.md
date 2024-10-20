# Reed Switch - kontaktron - binary sensor (aktywny)
## reed_switch
Ten przykładowy sensor powstał na bazie ESP8266 i czujnika kontaktoronowego. Czujnik podpięty jest do zasilania 3,3V i pinu D7. Gdy elementy czujnika zbliżą się do siebie, na pinie pojawia się napięcie. Gdy elementy są blisko siebie zwraca wartość 1 (np drzwi zamknięte), w przeciwnym wypadku 0
### Endpointy
1. `GET /status` - zwraca status urządzenia
1. `GET /value` - zwraca wartość czujnika ([0/1])
## reed_switch_active
### Aktywność
Gdy urządzenie wykryje zmianę wartości na czujniku (na przykład otwarcie/zamknięcie drzwi), wysyła `POST` request do serwera (na ten moment mock pythonowy do testowania). W parametrach wysyła on swoje id i nową wartość na czujniku.

# TODO
Na ten moment aktywność czujnika ma być rozwiązana tak, że po stronie serwerowej będzie wątek odpowiedzialny za nasłuchiwanie requestów od aktywnych czujników. Gdy dostanie taki request, może zidentyfikować czujnik po id, zapisać wartość do bazy i podjąć kolejne działania.

ID i adres serwera czujnik powinien otrzymać przez protokół, który jeszcze nie powstał.