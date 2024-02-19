# ib-cp-server

## What is it?

This is a proxy server with `POST /login` API for the [Interactive Brokers Client Portal API](https://interactivebrokers.github.io/cpwebapi/). Because the official server has no a simple `POST /login` api, it's hard to run as a headless service.

Ref: [Why can't automated login](https://interactivebrokers.github.io/cpwebapi/use-cases#automated-login)

## How it works?

1. This server offer a `POST /login` api to login to the official server.
2. When `POST /login` called, it will open the original login page by [chromedp](https://github.com/chromedp/chromedp) and automated fill the username and password (config by environment variables)
3. Then IB will send a push to your 2FA device, then you can confirm the login. So you have to enable a Two Factor Authentication (2FA) device. [iPhone](https://www.ibkrguides.com/iphone/sls/activating-ibkr-mobile.htm) or [Android](https://www.ibkrguides.com/android/sls/activating-ibkr-mobile.htm)
4. After login, you can simple use this server as the official server, because it will proxy all the requests to the official server.

Ref: [Official API Endpoints](https://interactivebrokers.github.io/cpwebapi/endpoints)

## How to use it?

### From source

1. Install [Go](https://go.dev/doc/install)
2. Download [clientportal.gw.zip](https://download2.interactivebrokers.com/portal/clientportal.gw.zip) or you can use the one in this repo. I will update it periodically. Unzip it to `clientportal.gw` folder or somewhere else.
3. Run `go run ./cmd/server` or `go build ./cmd/server` and run `./server`

### Docker

1. You can simple pull image by `docker pull ppaanngggg/ib-cp-server`
2. Or you can build it by `docker build -t ib-cp-server .`
3. Run it by `docker run -d --name=ib-cp-server -p 8000:8000 -e IB_USERNAME=your_username -e IB_PASSWORD=your_password ppaanngggg/ib-cp-server`

## Environment Variables

1. `IB_USERNAME`: Your IB Account username
2. `IB_PASSWORD`: Your IB Account password
3. `IB_EMBEDDED`: If you want main program help you to start the official server, set it to `true`. Default is `ture`
4. `IB_EXEC_DIR`: if `IB_EMBEDDED` is `true`, you can set this to the folder of `clientportal.gw` folder. Default is `./clientportal.gw`
5. `IB_URL`: If you start the official server by yourself, you can set this to the url of the official server. Default is `https://localhost:5000`
6. `SERVER_HOST`: The host of this server. Default is `0.0.0.0`
7. `SERVER_PORT`: The port of this server. Default is `8000`
8. `SERVER_TIMEOUT`: Request timeout. Default is `60s`

## Thanks

<a href="https://www.buymeacoffee.com/ppaanngggg" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: 41px !important;width: 174px !important;box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;-webkit-box-shadow: 0px 3px 2px 0px rgba(190, 190, 190, 0.5) !important;" ></a>
