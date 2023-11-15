# Golang Support Vector Machine xample implementation

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




