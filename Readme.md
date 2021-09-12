# simple_reverse_proxy

A simple reverse proxy with middlewares

Right now it's only a prototype to help http debuging & profiling.

Supported middlewares are : 
- **log_request** to be able to have metrics on requests
- **cors** to enable cors on responses
- **dump_request** to be able to see what's exactly send to the server

## Configuration

Create a file config.yaml with the following content and adapt to your needs 

```yaml
server:
  listenPort: 8889
  forwardUrl: http://www.google.fr
middlewares:
  - cors
  - log_request
  - dump_request
```

## TODO

Add tests