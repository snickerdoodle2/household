# HouseHold
## About The Project
This code was developed as part of the project _"A system for managing smart home infrastructure optimizing the consumption of natural resources."_ The project was undertaken as part of a Bachelor of Engineering degree at AGH University of Science and Technology (AGH UST).

The goal of the project was to create a smart home management system, allowing the easy optimisation of resource consumption and allowing the simple integration of customised devices. Based on user-defined policies and sensor input, the system is able to control effectors, increasing user comfort while reducing resource consumption. The basic components of the system are devices, a web-based client application, and a backend that acts as a server for the client application and a client for the devices.

The devices could be differentiated into sensors and effectors. Sensors provide data to the system in an active or passive manner. Passive sensors are queried by the system every user-defined time interval, while active sensors send data on their own. Effectors have a state that can be changed by the user. These devices could be either physical or virtual, running locally or in the cloud. The only requirement for the devices is to expose an appropriate API through which the backend communicates with them. Software templates for ESP-based boards are provided, as well as a virtual weather sensor server and a virtual sensor for solar installation data.

The application provides rule and sequence mechanisms. A rule is a logical expression based on the state of the sensors. When the expression is satisfied, the system automatically changes the state of an effector or executes a sequence. A sequence is a change of state of specific effectors, in a specific order, with an optional time interval in between. A sequence can be executed as an effect of a rule or can be manually triggered by the user. 

The web application, the frontend of the system, presents the data collected from the sensors in a clear and aesthetically pleasing way. The main page of the application displays graphs of the data collected in real time. Using a clear and intuitive interface, the user is able, among other things, to manage devices, create and manage sequences and rules, and receive notifications.

Conscious use of all the application's components gives the user the opportunity to create a personalized system, facilitating household management while reducing resource consumption. Creating your own devices and adding them to the system is intuitive and simple, and the ability to create your own rules and sequences allows you to automate everyday tasks.

### Built With

* [![Svelte][Svelte.dev]][Svelte-url]
* [![Bun][Bun.dev]][Bun-url]
* [![Tailwind CSS][Tailwind.dev]][Tailwind-url]
* [![Go][Go.dev]][Go-url]
* [![Docker][Docker.dev]][Docker-url]
* [![Just][Just.dev]][Just-url]
* [![Migrate][Migrate.dev]][Migrate-url]
* [![Postgres][Postgres.dev]][Postgres-url]

## Getting Started

### Backend

#### Requirements

1. `golang 1.22`
2. `docker`
3. `just` - [github](https://github.com/casey/just)
4. `watchexec` - [github](https://github.com/watchexec/watchexec)
5. `migrate cli` - [github](https://github.com/golang-migrate/migrate)

#### Setup

##### `.env` Configuration
Fill in the `.env` file:
```bash
DATABASE_ADDRESS=localhost:5432
POSTGRES_PASSWORD=<DB_PASSWORD>
POSTGRES_DB=<DB_NAME>
DATABASE_URL=postgresql://postgres:<DB_PASSWORD>@localhost:5432/<DB_NAME>?sslmode=disable
```

#### First run
1. `docker compose up -d`
2. `just up`
3. `just run`

### Frontend
#### Requirements

1. `bun` - [bun.sh](https://bun.sh/)

#### Setup

2. `bun install`

#### Run

1. `bun run dev`

## License

Distributed under the XXXXXXXX See `LICENSE.txt` for more information.

[Bun.dev]: https://img.shields.io/badge/Bun-4A4A55?style=for-the-badge&logo=bun&logoColor=ffffff
[Bun-url]: https://bun.sh/

[Tailwind.dev]: https://img.shields.io/badge/Tailwind%20CSS-4A4A55?style=for-the-badge&logo=tailwindcss&logoColor=38B2AC
[Tailwind-url]: https://tailwindcss.com/

[Go.dev]: https://img.shields.io/badge/Go-4A4A55?style=for-the-badge&logo=go&logoColor=00ADD8
[Go-url]: https://go.dev/

[Docker.dev]: https://img.shields.io/badge/Docker-4A4A55?style=for-the-badge&logo=docker&logoColor=2496ED
[Docker-url]: https://www.docker.com/

[Just.dev]: https://img.shields.io/badge/Just-4A4A55?style=for-the-badge&logo=just&logoColor=ffffff
[Just-url]: https://github.com/casey/just

[Migrate.dev]: https://img.shields.io/badge/Migrate-4A4A55?style=for-the-badge&logo=postgresql&logoColor=336791
[Migrate-url]: https://github.com/golang-migrate/migrate

[Postgres.dev]: https://img.shields.io/badge/Postgres-4A4A55?style=for-the-badge&logo=postgresql&logoColor=336791
[Postgres-url]: https://www.postgresql.org/

[Svelte.dev]: https://img.shields.io/badge/Svelte-4A4A55?style=for-the-badge&logo=svelte&logoColor=FF3E00 
[Svelte-url]: https://svelte.dev/
