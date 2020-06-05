# shortr
Simple url shortener in various languages and frameworks

_**Note:** this is just a simple service to learn various languages and frameworks. Efficiency nor security or mantainabilty are intended in this project._

## Flows
### `GET` <span style="font-weight: normal; font-size: 0.8em;">/*<span/>
### `POST` <span style="font-weight: normal; font-size: 0.8em;">/*<span/>
### `DELETE` <span style="font-weight: normal; font-size: 0.8em;">/*<span/>
### `PUT` <span style="font-weight: normal; font-size: 0.8em;">/*<span/>
### `GET` <span style="font-weight: normal; font-size: 0.8em;">/*/stats<span/>

## Languages
- **Golang**
    - [Echo](go/echo/README.md)
- **JavaScript**
    - [Deno](js/deno/README.md)
    - [Express](js/express/README.md)
- **Python**
    - [Flask](py/flask/README.md)
- **Rust**
    - [Actix-Web](rs/actix/README.md)

## Setup
- Install:
    - [`docker 19.03.6 >=`](https://docs.docker.com/get-docker/)
    - [`docker-compose 1.21.0 >=`](https://docs.docker.com/compose/install/)
- Run `make $language/$framework` ( for example `go/echo` )

See [`makefile`](makefile) for further commands.

## Overall Comparison
These are not _good_ comparisons nor benchmarks, but gives a quick overlook at language and framework efficiency. The benchmarking tool used is [Apache's AB](https://httpd.apache.org/docs/2.4/programs/ab.html) with `ab -n 1000000 -k -c 30 -q http://localhost:80/benchmark`.

### `GET CACHED` <span style="font-weight: normal; font-size: 0.8em;">/*<span/>
### `GET` <span style="font-weight: normal; font-size: 0.8em;">/*<span/>

## License
This project is licensed under the [MIT License](https://opensource.org/licenses/MIT) - read the [LICENSE](LICENSE) file for details.