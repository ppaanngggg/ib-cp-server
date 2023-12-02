# ib-cp-server

## What is it?

This is a proxy server with `/login` API for the [Interactive Brokers Client Portal API](https://interactivebrokers.github.io/cpwebapi/).Because the official server has no a simple `/login` api, it's hard to run as a headless service.

Ref: [Why can't automated login](https://interactivebrokers.github.io/cpwebapi/use-cases#automated-login)

## How it works?

1. This server offer a `/login` api to login to the official server.
2. When `/login` called, it will open the original login page by [chromedp](https://github.com/chromedp/chromedp) and automated fill the username and password (config by environment variables)
3. Then IB will send a push to your 2FA device, then you can confirm the login. So you have to enable a Two Factor Authentication (2FA) device. [iPhone](https://www.ibkrguides.com/iphone/sls/activating-ibkr-mobile.htm) or [Android](https://www.ibkrguides.com/android/sls/activating-ibkr-mobile.htm)
4. After login, you can simple use this server as the official server, because it will proxy all the requests to the official server.

Ref: [Official API Endpoints](https://interactivebrokers.github.io/cpwebapi/endpoints)

## How to use it?

### From source

TODO

### Docker

TODO