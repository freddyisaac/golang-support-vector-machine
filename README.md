# Golang Support Vector Machine Example implementation

## Introduction

need an intro here

## Generating Test data

Build the gen_data program using 

```go build gen_data.go utils.go math.go```

The script ``gen.sh`` shows how to create a simple 2d point set with two clusters.

It should crate an output file data.txt which can be viewed in R using the functions in plot_routines.r

Enter R and at the command prompt type

```
> source("plot_routines.r")
> plotpoints("data.txt")
```

If that goes well you should see something looking similar to this

![points](https://github.com/freddyisaac/support-vector-machine/assets/40456262/5281a67c-a451-40e1-b2ad-5f9ea349aae7)

## Using the SVM example

Now that we have some data in data.txt with two clusters of points centered around ``(0.25,0.25)`` and ``(0.75,0.75)``

create the example with

```go build -o 2d_ex 2d_ex.go svm.go math.go utils.go```

It takes quite a few params and is a bit noisy to the terminal but essentially entering

```2d_ex -if data.txt -lbx 0.0 -ubx 1.0 -lby 0.0 -uby 1.0 -of out.txt```

All being well it should produce a 50 by 50 grid which can be plotted as contours from within R again using

```
> source("plot_routines.r")
> cplot("out.txt", 50)
```

If all is well then you should see something like this

![contours](https://github.com/freddyisaac/support-vector-machine/assets/40456262/91451beb-d723-4c7a-9678-9c1085759a36)







