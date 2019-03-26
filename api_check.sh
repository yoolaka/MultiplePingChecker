curl -X GET 'http://127.0.0.1:9335/ping/www.naver.com?wait=true' 
curl -X POST -d 'server=www.naver.com&count=100' 'http://127.0.0.1:9335/ping'
curl -X POST -d 'server=google.com&count=50' 'http://127.0.0.1:9335/ping'
curl -X POST -d 'server=google.com&count=100' 'http://127.0.0.1:9335/ping'
#curl -X GET 'http://127.0.0.1:9335/ping/www.naver.com?wait=true' 
curl -X DELETE 'http://127.0.0.1:9335/ping/google.com'
curl -X GET 'http://127.0.0.1:9335/ping' --output output
