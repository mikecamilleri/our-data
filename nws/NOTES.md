# National Weather Service Client

- relies on github.com/mikecamilleri/cap
- the APIs use multiple ways of identifying information. The client by default should use the most precise way available given the endpoint and client configuration

## NWS Public Alerts (CAP) in ATOM feeds

- by state https://alerts.weather.gov/cap/us.php?x=0 where `us` may be replaced by a state code

- by zone https://alerts.weather.gov/cap/wwaatmget.php?x=ORZ006&y=0 where ORZ006 may be replaced by any zone from https://alerts.weather.gov

- by county https://alerts.weather.gov/cap/wwaatmget.php?x=ORC051&y=0 where ORC051 may be replaced by any county from https://alerts.weather.gov

- by zone and by county both hang with the connection open during testing with httpie

## NWS National Digital Forecast Database

- ref https://graphical.weather.gov/xml/rest.php

- use REST API

- confusing API -- look at later

## NWS Current Weather Conditions

- by station callsign http://w1.weather.gov/xml/current_obs/KPDX.xml where KPDX may be replaced with any call sign

- use XSD referenced in output for field list
