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



