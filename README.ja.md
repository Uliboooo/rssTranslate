# rssTranslate

`rssTranslate`は、RSSフィードを取得し、Gemini APIを使用してタイトルを日本語に翻訳し、クリーンで読みやすいウェブページに表示するシンプルなウェブサーバーです。

## 前提条件

-   Goがインストールされていること。
-   Gemini APIキー。

## セットアップ

1.  **リポジトリをクローンする:**
    ```bash
    git clone <repository-url>
    cd rssTranslate
    ```

2.  **Gemini APIキーを設定する:**
    アプリケーションは`GEMINI_API_KEY`環境変数が設定されている必要があります。
    ```bash
    export GEMINI_API_KEY="あなたのAPIキー"
    ```

## 実行方法

1.  **サーバーを起動する:**
    プロジェクトのルートディレクトリで次のコマンドを実行します。
    ```bash
    go run main.go
    ```
    デフォルトでは、サーバーはポート`58877`で起動し、Hacker Newsの「newest」RSSフィードを使用します。

2.  **アプリケーションを表示する:**
    ウェブブラウザで`index.html`ファイルを開き、翻訳されたRSSフィードを表示します。

## 設定

コマンドラインフラグを使用してサーバーの動作をカスタマイズできます。

-   `-p`: サーバーのポートを設定します。
    -   例: `go run main.go -p 8080`
    -   デフォルト: `58877`
-   `-url`: 取得するRSSフィードのURLを指定します。
    -   例: `go run main.go -url "https://example.com/feed.xml"`
    -   デフォルト: `https://hnrss.org/newest`
-   `-api`: APIのURLを設定します。
    -   デフォルト: `http://localhost`

## 仕組み

Goバックエンドは、`/api/rts`エンドポイントを持つウェブサーバーを起動します。このエンドポイントが呼び出されると、指定されたRSSフィードを取得し、タイトルを抽出し、Gemini APIに送信して翻訳します。翻訳されたタイトルは、リンクと公開日とともにJSONオブジェクトとして返されます。

フロントエンドはシンプルな`index.html`と`script.js`で構成されています。JavaScriptは`/api/rts`エンドポイントに`fetch`リクエストを送信し、翻訳されたコンテンツを動的にページに表示します。
