The application has no external dependencies, besides Docker and docker-compose,
if you want to use them to run it.

To run the application with `docker-compose` do

    docker-compose up -d

The application listens to port 8282, on localhost.

To run the application with plain `docker` do

    docker -t pservice build .
    docker -d -p 8282:8282 run pservice

If you do not have or want to use Docker, a simple

    go run

will suffice.

The endpoint follows the full specification described in the PDF. It accepts
periods of the form "xh", "xd", "xmo", "xy", where `x` is any number greater
than zero.  The timezone argument is optional and defaults to Local, if empty or
omitted.

Successful requests return a JSON array with the results (which can be empty if
t1=t2) and status code 200.  Unsuccessful requests return an error object as
shown in the PDF and status code 400.
