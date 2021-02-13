# mandelbrot
A Mandelbrot generating program written in Go. Written in about 2 hours at this point so a little basic, but I have introduced some concurrency.

## TODO
- fiddle with the colours. I'm sure there are some resources online for different palettes.
- make it a web app. The end goal with this is to implement smooth zoom with a mouse, and I'd quite like to be able to do this in a browser. I have no idea how feasible this is.
- zoom algorithm. I'm guessing this is where it gets complicated. Just off the top of my head, I'm seeing that there'd be "zoom thresholds", i.e. once the zoom area is small enough the image has to be reloaded. Again, I have no idea yet how this can be done smoothly. Also, what is the limit of zoom? THe straightforward approach is obviously bounded by the precision afforded by a float64, but are there some tricks I can pull off to get around this? Maybe that's just moving the problem.
