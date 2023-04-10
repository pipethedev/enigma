# Engima

Engima is a CLI application used to generate hermes token
## Installation

Clone the repository to install enigma.

```bash
$ git clone https://github.com/pipethedev/enigma.git
```

```bash
$ cd enigma
```

```bash
$ go mod vendor
```

```bash
$ go run ./cmd/enigma
```
>Note: You can download the executable from [here](https://github.com/pipethedev/enigma/tags) instead if you don't want to set it up using the repository.

## Usage

Create new hermes token

```sh
go run ./cmd/enigma -create
```

Fetch existing hermes token

```sh
go run ./cmd/enigma -get
```

## License

[MIT](https://choosealicense.com/licenses/mit/)