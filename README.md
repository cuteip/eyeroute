# eyeroute
WIP: Looking Glass

## dev メモ

リポジトリルートディレクトリ（`./`）に Go のバックエンドが、`./front/` に React、TypeScript のフロントエンドがある。

修正しながらやるときは、以下のように、フロントエンド側は `npm start` で起動し、バックエンド側は `--dev-upstream-front-url` でフロントエンド側 URL を指定する。

```
$ cd front
$ npm start
```

```
$ make eyeroute
$ ./eyeroute server run --dummy --dev-upstream-front-url http://127.0.0.1:3000
```

`--dummy` を渡すと、動作確認時に毎度 `mtr` を実行するのではなく、ダミーのデータを返すようになる。
