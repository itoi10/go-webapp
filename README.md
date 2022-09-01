# Go API アプリ

## 起動方法

```sh
$ make up
```

```./auth/cert/```以下に```public.pem```と```secret.pem```が必要

ルーティング
```./mux.go```

ハンドラ
```./handler```

エンドポイント例 (ユーザ登録、ログイン)
```sh
$ curl -X POST localhost:8080/register -d '{"name": "normal_user1", "password":"test", "role":"user"}'
{"id":1}
curl -X POST localhost:8080/login -d '{"user_name":"normal_user1", "password":"test"}'
{"access_token":"<JWT token>"}

```


## 参考書籍

[清水 陽一郎. 詳解Go言語Webアプリケーション開発. シーアンドアール研究所. 2022](https://www.amazon.co.jp/dp/4863543727/)
