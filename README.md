# EProxy
Simple proxy server (improved server from [medium](https://medium.com/@mlowicki/http-s-proxy-in-golang-in-less-than-100-lines-of-code-6a51c2f2c38c).

# Configuration
Proxy has 4 startup params. To get help you have to run app with **-h** flag.
Usage of ./eproxy:
~~~sh
  -key string
    	path to key file (default "server.key")
  -pem string
    	path to pem file (default "server.pem")
  -port string
    	Proxy port (default "8888")
  -proto string
    	Proxy protocol (http or https) (default "https")
~~~

# Issues
Unfortently, HTTPS proxy doesn't work on Ubuntu with self signed certificates (won't fix).
