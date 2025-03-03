# klines-stat
get historical klines and build new ones from resent trades data

## Base primitives

 - Storage: a worker reponsible for saving historical klines, new klines and resent trade data. Service launches one isntance of this type
 - REST Puller: a worker to pull historical klines via REST api. Service launches one instance for `pair-timeframe` combination. E.g `BTC_USDT-1m`, `BTC_USDT-15m`, `BTC_USDT-60m`.
 - WS Puller: a worker for getting actual recent trading data (RT). Service launches one instance of this type
 - Kline Builder: a worket to build klines from RT data. Service launches one instance for `pair-timeframe` combination. E.g `BTC_USDT-1m`, `BTC_USDT-15m`, `BTC_USDT-60m`.


## Base flow

REST Pullers load and save all completed historical klines since speicified time. Each worker of this type returns uncompleted kline for the specific pair and timeframe.
Kline builders start work with these uncompleted klines as initial input values. They RT data comes to it worker via channel from the WS Puller. If current kline is completed a worker sends it to saver via channel.


## What needs to be done

Build VBS data for historical klines and calculate it for klines which is getting build from RT data.


## Launch

The solution might be launched in a completed docker environment with automatically applied DB migrations.

```
## Run the solution in docker-compose environment
make compose-run

## connection to the DB
make db-connect

## Other commands description
make help
```


## NOTE

Please keep in mind that import of `github.com/wuhewuhe/bybit.go.api v1.0.18` will be replaced by my fork of this repo.
I found this solution very useful but there are several issues causes build and import error. My fork contains a quick fix.
