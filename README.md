# go-pipeline-timeout
外部プログラム コマンド を実行し、その　標準出力、標準エラー 、終了コード、　エラーを返却します。   
linuxの無名パイプ付きコマンドに対応し、コマンドのタイムアウト、標準入力、環境変数　に対応しています。   

ruby の　[Open3.#popen3](https://docs.ruby-lang.org/ja/latest/method/Open3/m/popen3.html) を参考にしてます。


### Usage
詳しい使い方は　、testコードを参照してください。

```go
import (
	. pipeline_timeout
	fmt
)

stout, sterr , code , err := Exec("cat /etc/passwd|egrep ^root")
```

### Thanks
このモジュールは　[mattn/go-shellwords](https://github.com/mattn/go-shellwords) を使用しています。   

### Author
keiichi ishioka
