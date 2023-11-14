plotpoints <- function(name) {
	a_data <- read.table(name)
	a_x <- a_data[[1]]
	a_y <- a_data[[2]]
	plot(a_x, a_y)
}

cplot <- function(name, num) {
	my_data <- read.table(name)
	vec_x <- my_data[[1]]
	vec_y <- my_data[[2]]
	vec_z <- my_data[[3]]
	
	m_x <- matrix(vec_x, nrow=num)
	v_x <- m_x[,1]
	v_y <- v_x
	m <- matrix(vec_z, nrow=num)
	contour(m)
}

