# Complex set art

Generate artistic Julia Set renders. More about Julia Set [here](https://en.wikipedia.org/wiki/Julia_set).<br>
You can check out another (web based) Julia Set visualization of mine [here](https://bartolomej.github.io/julia-set).

| | | |
|:-------------------------:|:-------------------------:|:-------------------------:|
|<img width="500" alt="C = 0 + 0i" src="docs/examples/0.00.0i.png">  C = 0 + 0i |  <img width="500" alt="C = -0.43 - 0.2i" src="docs/examples/-0.43-0.2t.png"> C = -0.43 - 0.2i|
|<img width="500" alt="C = 0.61 + 0.52i" src="docs/examples/grey_0.61_0.52.png">  C = 0.61 + 0.52i |  <img width="500" alt="C = -0.81 + 0.0i" src="docs/examples/gold_-0.81_0.0.png"> C = -0.81 + 0.0i|
|<img width="500" src="docs/examples/giphy1.gif">  |  <img width="500" src="docs/examples/giphy3.gif"> |

## Usage

```bash
// download package
go get -u github.com/bartolomej/complex-set-art

// move to project root
cd ~/go/src/github.com/bartolomej/complex-set-art

// run example render
go run *.go default-image
```

Julia set is generated using a function of a complex domain: `Z = Z^2 + C` <br>
All generated files are saved at `github.com/bartolomej/complex-set-art/out/`.


### Cli

Arguments syntax:
```bash
go run *.go <image-size> <Re(C)> <Im(C)> <output-file>
```

Example cli run command:
```bash
go run *.go 3000 0.37 -0.4
```

Run with default config:
```bash
go run *.go default-video OR default-image
```

### Custom config

Create `renderers.json`  configuration file in root directory.<br>
Example configuration for generating static renders:

##### 1. Static image rendering
```json
{
  "id": "smooth-boundaries",
  "resolution": 1000,
  "returnMode": "I",
  "encoding": "png",
  "maxIterations": 20.0,
  "maxThreshold": 30.0,
  "color": "HSV(c*2 + 190, tanh(c), 1)",
  "static": {
     "originX": 0.0,
     "originY": 0.0,
     "axisSpan": 1.5,
     "realC": 0.0,
     "imagC": 0.8,
     "realExp": 3,
     "imagExp": 0.0
  }
}
```

##### 2. Video rendering
```json
{
  "id": "stripy-video",
  "resolution": 1500,
  "returnMode": "D",
  "encoding": "png",
  "maxIterations": 30.0,
  "maxThreshold": 30.0,
  "color": "HSV(c * 10, 100, 0)",
  "duration": 10,
  "fps": 30,
  "start": {
     "originX": 0.0,
     "originY": 0.0,
     "axisSpan": 2.0,
     "realC": 0.3,
     "imagC": 0.3,
     "realExp": 2.0,
     "imagExp": 0.0
  },
  "end": {
     "originX": 0.0,
     "originY": 0.0,
     "axisSpan": 1.5,
     "realC": -0.3,
     "imagC": -0.3,
     "realExp": 3.0,
     "imagExp": 0.0
  }
}
```

Run above configuration with the following commands:
```bash
// pass config id as first param
go run *.go zoomed-in
```

## TODO
- display video rendering progress
- add extensive documentation
- add configs examples
- array of animation keyframes ?
- rendering job queue
- support custom set formulas (quadratic,..)
- improve performance with multithreading
- simple gui app ? (https://github.com/andlabs/ui)