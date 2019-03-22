curl -X POST -d 'server=google.com&count=100' 'http://127.0.0.1:9335/ping'
curl -X GET 'http://127.0.0.1:9335/ping/google.com?wait=true' 
