# Complex set art

Generate artistic Julia set images and save them to image format.

![](https://i.ibb.co/cbwxH60/out1.png)

## Usage

```bash
go get -u github.com/bartolomej/complex-set-art
cd ~/go/src/github.com/bartolomej/complex-set-art
```

Julia set is generated using a function of a complex domain: `Z = Z^2 + C` <br>
All generated files are saved at `github.com/bartolomej/complex-set-art/out/`.

You can use this tool in 3 different ways:
- run without any arguments (default settings) `go run *.go`
- pass 4 arguments via cli
- define configurations in a file and run with


### Cli
When running main.go you have to pass parameters for:
 - output image size 
 - real part of *C* constant
 - imaginary part of *C* constant
 - set generation mode 
    - "-i" - by iteration limit
    - "-t" - by value threshold)

```bash
go run *.go <image-size> <Re(c)> <Im(c)> <generation-mode> <outputfile>
```

Example cli run command:
```bash
go run *.go 3000 0.37 -0.4 -i
```

### Configuration file
Create `renderers.json`  configuration file in root directory.<br>
Example configuration for generating static image:
```json
[
  {
      "id": "1",
      "resolution": 3000,
      "renderMode": "-t",
      "encoding": "jpeg",
      "static": {
        "centerX": 0,
        "centerY": 0,
        "axisSpan": 4,
        "mode": "-t",
        "realC": -0.6,
        "imagC": 0.5
      }
    }
]
```

Run above configuration with this command:
```bash
go run *.go 1
```