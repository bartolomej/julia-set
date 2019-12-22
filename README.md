# Complex set art

Generate artistic Julia set images and save them to image format.

![](https://i.ibb.co/cbwxH60/out1.png)

## Usage

Julia set is generated using a function of a complex domain: `Z = Z^2 + C` <br>

When running main.go you have to pass parameters for:
 - output image size 
 - real part of *C* constant
 - imaginary part of *C* constant
 - set generation mode 
    - "-i" - by iteration limit
    - "-t" - by value threshold)

```bash
go get -u github.com/bartolomej/complex-set-art
cd ~/go/src/github.com/bartolomej/complex-set-art
go run *.go <image-size> <Re(c)> <Im(c)> <generation-mode> <outputfile>
```

## Examples
This command will render .png image of the set (at `C = -0.37 + -0.4`), with a resolution of 3000px,
and export it to `github.com/bartolomej/complex-set-art/out/out4.png`.
```bash
go run *.go 3000 0.37 -0.4 -i out6.png
```