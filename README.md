# Complex set art

Generate artistic Julia set images and save them to image format.

![](https://i.ibb.co/cbwxH60/out1.png)

## Usage
Julia set is generated using a function of a complex domain: `Z = Z^2 + C` <br>
When running main.go you have to pass parameters for output image
size, real part of *C* constant and imaginary part of *C* constant.

```bash
go get -u github.com/bartolomej/complex-set-art
cd ~/go/src/github.com/bartolomej/complex-set-art
go run main.go <image-size> <Re(c)> <Im(c)> <outputfile>
```