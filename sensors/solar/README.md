# GROWATT solar panels integration
Serwer napisany w pythonie, odpowiedzialny za przetwarzanie danych i komunikację z instalacją fotowoltaiczną. Growatt nie ma oficjalnego API, więc użyłem endpointów znalezionych w projektach z czeluści githuba. 
## Endpointy
- `GET /value` - zwraca średnią moc instalacji z ostatnich 5 minut
