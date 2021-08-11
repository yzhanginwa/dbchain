### dbchain nginx 配置

nginx配置有两个配置，主要是控制外网访问api

#### 配置一(外网配置)：

```
server {
		
	    listen 80; #访问端口
        server_name  yitaibox.com www.yitaibox.com; #外网访问地址
        root         /var/www/ytboxweb/dist/;

        location / {
            try_files $uri $uri/ /index.html;
        }

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }

        location /relay/ipfs/ {
            rewrite /relay/ipfs/(.*) /ipfs/$1  break;
            proxy_pass http://172.20.0.101:8080;
            proxy_set_header origin 'http://172.20.0.101:8080';
        }

        location /relay/dbchain/upload/ {
          rewrite /relay/dbchain/upload/(.*) /dbchain/upload/$1  break;
          proxy_pass http://172.20.0.101:1318;
          proxy_set_header origin 'http://172.20.0.101:1318';
        }
        
        location /relay/dbchain/oracle/bsn/ {
        	return 404
        }

        location /relay/dbchain/oracle/ {
          rewrite /relay/dbchain/oracle/(.*) /dbchain/oracle/$1  break;
          proxy_pass http://172.20.0.101:1318;
          proxy_set_header origin 'http://172.20.0.101:1318';
        }

        location /relay/ {
            rewrite /relay/(.*) /$1 break;
            proxy_pass http://172.20.0.101:1317;
            proxy_set_header origin 'http://172.20.0.101:1317';
        }
}

```

#### 配置二(内网访问)：

```
server {
		
	    listen 8000; #访问端口，应该与外网
        server_name  192.168.xx.xx; #内部地址
        root         /var/www/ytboxweb/dist/;

        location / {
            try_files $uri $uri/ /index.html;
        }

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   html;
        }

        location /relay/ipfs/ {
            rewrite /relay/ipfs/(.*) /ipfs/$1  break;
            proxy_pass http://172.20.0.101:8080;
            proxy_set_header origin 'http://172.20.0.101:8080';
        }

        location /relay/dbchain/upload/ {
          rewrite /relay/dbchain/upload/(.*) /dbchain/upload/$1  break;
          proxy_pass http://172.20.0.101:1318;
          proxy_set_header origin 'http://172.20.0.101:1318';
        }

        location /relay/dbchain/oracle/ {
          rewrite /relay/dbchain/oracle/(.*) /dbchain/oracle/$1  break;
          proxy_pass http://172.20.0.101:1318;
          proxy_set_header origin 'http://172.20.0.101:1318';
        }

        location /relay/ {
            rewrite /relay/(.*) /$1 break;
            proxy_pass http://172.20.0.101:1317;
            proxy_set_header origin 'http://172.20.0.101:1317';
        }
}
```

