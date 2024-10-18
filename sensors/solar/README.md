# GROWATT solar panels integration
Serwer napisany w pythonie, odpowiedzialny za przetwarzanie danych i komunikację z instalacją fotowoltaiczną. Growatt nie ma oficjalnego API, więc użyłem endpointów znalezionych w projektach z czeluści githuba. 
## Endpointy
- `GET /value` - zwraca średnią moc instalacji z ostatnich 5 minut
- `GET /status` - zwraca status i typ czujnika (decmial_sensor)
## Obsługa błędów
W przypadku błędu logowania, lub pobierania danych z serwera GroWatt, zwracany jest response z kodem HTTP 500, oraz stosownym komentarzem