# Serwer Danych Pogodowych

Serwer zintegrowany z projektem Household - po skonfigurowaniu i uruchomieniu serwera, sensory automatycznie dodawane są do aplikacji Household, serwer pogodowy uruchamiany i zaczyna się standardowe działanie pasywnych sensorów.

Serwer korzysta z API open-meteo i umożliwia korzystanie z wszystkich danych znajdujących się w sekcjach **_Hourly Weather Variables_** i **_Current Weather_**. Jeden sensor zajmuje się przetwarzaniem tylko jednej zmiennej pogodowej.

## Konfiguracja

Konfiguracja serwera odbywa się poprzez plik config.json. W tym pliku należy zdefiniować parametry serwera pogodowego, serwera aplikacji _Household_ oraz listę sensorów pogodowych w odpowiednim formacie.

## Uruchomienie

Program uruchamiamy z linii komend wywołując:

```
python weather-sensor-server.py
```

### Parametry wywołania programu

- `-h, --help ` - wyświetlenie panelu pomocy
- `-u, --username` - konfiguracja nazwy użytkownika w _Household_
- `-p, --password` - konfiguracja hasła użytkownika w _Household_
- `-c, --configpath` - konfiguracja ścieżki do pliku konfiguracyjnego

Wszystkie parametry są opcjonalne, w przypadku niepodania username lub hasła, program podejmie próbę odczytania ich z pliku konfiguracyjnego. Domyślną lokalizacją pliku konfiguracyjnego jest katalog główny serwera (z niego wywołujemy komendę uruchamiającą), a jego nazwa to config.json

**Przykładowy fragment pliku konfiguracyjnego:**

```json
{
    "weather-server": {
        "host": "127.0.0.1",
        "port": 9001
    },
    "household-server": {
        "host": "172.30.227.16",
        "port": 8080,
        "username": "user",
        "password": "password"
    },
    "sensors": [
        {
            "name": "cloud-cover",
            "refresh_rate": 900,
            "processing": "current",
            "params": {
                "latitude": 50.3,
                "longitude": 19.28,
                "current": "cloud_cover"
            }
        },
        {
            "name": "rain sum next 24 hours",
            "refresh_rate": 3600,
            "processing": "sum_future",
            "number_of_hours": 24,
            "params": {
                "latitude": 50.3,
                "longitude": 19.28,
                "hourly": "rain",
                "forecast_days": 3
            }
        },
```

### Struktura najwyższego poziomu

- `weather-server` - konfiguracja hosta i portu na którym uruchmia się serwer pogodowy
- `household-server` - konfiguracja połączenia do systemu _Household_
- `sensors` - lista sensorów, każdy z odpowiednimi atrybutami

### obiekt `weather-server`

- `host` - host na którym uruchomi serwer pogodowy
- `port` - port na którym uruchomi się serwer pogodowy

### obiekt `household-server`

- `host` - host na którym działa serwer systemu _Household_
- `port` - port na którym działa serwer systemu _Household_
- `username` - [**opcjonalne**] nazwa użytkownika do autentykacji z serwerem _Household_, można go również podać jako parametr wywołania programu
- `password` - [**opcjonalne**] hasło użytkownika do autentykacji z serwerem _Household_, można je również podać jako parametr wywołania programu

### tablica `sensors`

Tablica sensors zawiera konfigurację każdego z sensorów. Każdy sensor to obiekt z następującymi polami:

- `name` - nazwa sensora
- `refresh-rate` - częstotliwość odświeżania danych z sensora w aplikacji _Household_ (w sekundach)
- `processing` - definiuje typ przetwarzania danych przez sensor. Dozwolone wartości:
  - `"sum_past"` - suma wartości przeszłych danych pogodowych
  - `"sum_future"` - suma wartości prognozowanych danych pogodowych
  - `"max_past"` - maksimum wartości przeszłych danych pogodowych
  - `"max_future"`- maksimum wartości prognozowanych danych pogodowych
  - `"min_past"` - minimum wartości przeszłych danych pogodowych
  - `"min_future"` - minimum wartości prognozowanych danych pogodowych
  - `"avg_past"` - średnia wartości przeszłych danych pogodowych
  - `"avg_future"` - średnia wartości prognozowanych danych pogodowych
  - `"current"` - dane aktualne
- `number_of_hours` - określa z ilu godzin w przód lub w tył (w zależności od pola processing) dane mają być przetwarzane. Jeśli `processing` jest `current` pole nie jest wymagane. Należy zadbać, aby dane otrzymane z _open-meteo_ miały odpowiednio długi timeframe by umożliwić przetwarzenie w określonej tutaj licznnie godzin
- `params` - obiekt parametrów zapytania wysyłanego do _open-meteo_. Najłatwiej jest otrzymać ten obiekt poprzez 'wyklikanie' interesujących nas danch na stronie https://open-meteo.com/en/docs , a następnie skopiowanie obiektu params z sekcji _Api Response_ > _Python_. Wspierane są dane otrzymane z sekcji _Hourly Weather Variables_ (zarówno dla _opcji forecast days_ jak i _past days_) oraz dane z sekcji _Current Weather_. Ważne aby obiekt ten nie zawierał zmiennych pogodowych poza tą jedną za której przetwarzanie odpowiedzialny jest sensor, gdyż może to spowodować nieoczekiwane działanie sensora.

Każdy sensor jest typu `decimal_sensor`

## Wymagania

- python (działa na 3.13.0)
- wymagane biblioteki:
  - `requests`
  - `flask`
  - `openmeteo_requests`
  - `numpy`
  - `pandas`

## Uwagi

- parametr `forecast_days` w `params` określa liczbę dni prognozy pogody włącznie z dniem bieżącym, niezalecane użycie wartości `1` gdyż będziemy mieli dostęp do prognozy wyłącznie z dnia bieżącego, wartość `3` zwróci prognozę dla dnia bieżącego, kolejnego i następnego po nim
