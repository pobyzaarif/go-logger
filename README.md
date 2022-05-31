# GO-LOGGER
go-logger is a custom log using [gommon/log](https://github.com/labstack/gommon) as a base log that can use in your golang project. the aim is to cover all logs (in common use such as inbound/outbound network log, query DB, error event, etc) with standard JSON log and ease for another engine to consume it. so we can focus in business logic and system flow. just tell the Infrastructure/DevOps team "pssst bro, pls take my app log, its shows up in stdout, then throw it to something, its trash but can be rare in some cases ;)".

## Features
For now, it's only integrated some popular go packages (from my own perspective cause in many cases I usually use them) such as echo v4, gorm, mongo-driver, and http.

## Implementation
Sample project that implements go-logger [POS_LITE](https://github.com/pobyzaarif/pos_lite)
