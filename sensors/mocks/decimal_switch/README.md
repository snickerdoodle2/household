# Decimal switch mock
### Wymagania:
- go `1.22`

Port można zmienić przy użyciu flagi `-port` (domyślnie `8888`)

## Endpointy
1. `GET /status` - zwraca status przełącznika
1. `GET /value` - zwraca wartość przełącznika (domyślnie 0.0)
1. `PUT /value` - ustawia wartość przełącznika na podaną przez użytkownika
