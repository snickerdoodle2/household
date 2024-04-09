# Binary switch mock
### Wymagania:
- go `1.22`

Port można zmienić przy użyciu flagi `-port` (domyślnie `8888`).
Opóźnienie, po którym przycisk się "wyłącza" można zmienić  przy użyciu flagi
`-delay` (domyślnie 2 sekundy)

## Endpointy
1. `GET /status` - zwraca status przełącznika
1. `GET /value` - zwraca wartość przełącznika (domyślnie false)
1. `POST /value` - włącza przycisk na określoną ilość czasu
