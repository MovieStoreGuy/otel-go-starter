# Collatz Conjecutre

Using a simple math conjecture known as the Collatz Conjecture, we can simply build function that recursively calls itself to generate tracing data.

The rules of the conjecture is simple, any number of N > 1 will always reach 4, 2, 1 computation loop with the following rules:
- If _n_ is odd, then _n_ => (3n+1)/2
- If _n_ is even, then _n_ => n/2

## How to visual the trace data

To keep things simple, you can either run the docker compose or the docker run command directly.

```shell
> docker-compose up -d
# Or
> docker run --rm -d -p "9411:9411" openzipkin/zipkin-slim:latest
```

Once that has started up, you can simplying run:
```shell
> go run main.go
```