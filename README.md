go-rocket-sample-app
====================

[github.com/acidlemon/rocket](https://github.com/acidlemon/rocket) のサンプルアプリ

### 内容

herokuで動かすことを前提にしたリマインダアプリ

指定のメールアドレスに対して、件名に日時指定、本文にリマインダ内容を書いてメールを送ると指定日時にリマインダが届く

### 使用APIサービス

* cloudmailin
 * メールを受信してアプリのエンドポイントを叩くのに使用
* mailgun
 * メール送信に使用

### 動かし方

設定は環境変数でする　設定は以下の通り

```
REGISTER_ADDRESS=メールアドレス登録用アドレス(cloudmailin)
SIGNUP_ADDRESS=リマインダ登録用アドレス(cloudmailin)
APP_NAME=アプリ名
APP_ADDRESS=アプリからの送信メールアドレス
APP_URL=herokuのアプリURL
MAILGUN_APIKEY=mailgunのapikey
MAILGUN_DOMAIN=mailgunのdomain
DATABASE_URL=postgresのdsn
```
