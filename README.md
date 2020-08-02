# Complex set art

Generate artistic images of Julia Set. More about [Julia Set](https://en.wikipedia.org/wiki/Julia_set).
I made this program while I was learning Go language for the first time.

I've also made another web based visualization of Julia set [here](https://bartolomej.github.io/julia-set).

| | | |
|:-------------------------:|:-------------------------:|:-------------------------:|
|<img width="500" alt="C = 0 + 0i" src="docs/examples/blue_0.0_0.8.png">  C = 0 + 0i |  <img width="500" alt="C = -0.43 - 0.2i" src="docs/examples/-0.43-0.2t.png"> C = -0.43 - 0.2i|
|<img width="500" alt="C = 0.61 + 0.52i" src="docs/examples/grey_0.61_0.52.png">  C = 0.61 + 0.52i |  <img width="500" alt="C = -0.81 + 0.0i" src="docs/examples/gold_-0.81_0.0.png"> C = -0.81 + 0.0i|
|<img width="500" src="docs/examples/giphy1.gif">  |  <img width="500" src="docs/examples/giphy3.gif"> |

## Installation

Clone project in your project directory.
```bash
git clone github.com/bartolomej/julia-set && cd julia-set
```

## Usage

Requires a configured [Go environment](https://golang.org/doc/install).

### Default config

Run `default` scene configuration stored in `example.config.json`.
```
go run main.go
```

Basic usage syntax:
```bash
go run *.go <image-size> <Re(C)> <Im(C)> <output-file>
```

### Custom config

Create `config.json` scene configuration file in root directory.
Check out examples in `example.config.json` file.

Run custom scene configuration with the following syntax:
```bash
go run *.go <scene-id>
```