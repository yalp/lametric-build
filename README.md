# LaMetric app for Travis and Coveralls status

![screencast](lametric-build.gif)

# LaMetric setup

Go to the [LaMetric Developper site](https://developer.lametric.com)

And create an app with 3 screens:
- the first is a simple "Name" screen for the Travis build status
- the second is a "Goal"  screen for the coverage report
- the third is a simple "Name" for the coverage change metric

Set the app to use Push mode.

Publish you app as a private app

Keep a copy of the URL and Access Token.

# Webhooks setup

The app expect Travis and Coverals to calls though
their respective webhooks.
Each service must target a specific path.

## Travis CI

For Travis CI, add the url to service with:

    notifications:
        webhooks: http://server:port/travis
or

    notifications:
        webhooks:
          - http://other/server
          - http://server:port/travis

Don't forget to append the `/travis` path to the end the URL.

## Coveralls

For Coveralls, add the webhook url to the notifications section.

Don't forget to append the `/coversalls` path to the end of the URL.

# Build

This requires the go toolchain:

    go get github.com/yalp/lametric-build

# Running

To run properly, the service requires some env vars:
- `PORT`: port to use
- `LAMETRIC_URL`: URL to your LaMetric private app
- `LAMETRIC_TOKEN`: Token for your LaMetric private app

Example:

    PORT=8082 LAMETRIC_URL="http://localhost:9090" LAMETRIC_TOKEN="ABCD" ./lametric-build

# Licence

Public domain.
I'm not responsible if your house burns, nor anything else.

# TODO

* Add configuration file
* Make Travis or Coveralls services optional
* Support more cover and build services
* Make some screens optional
