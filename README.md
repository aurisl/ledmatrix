# LED Matrix

Application to control AdaFruit LED Matrix. The project made on basis of https://www.adafruit.com/product/2345

## Installing

git clone git@github.com:aurisl/ledmatrix.git

`./build.sh hardware` //or software if building for emulator

@Todo write needed C binding modifications

Add init.d service
```
....
```

Run on raspberryPI startup ....

## Hardware

The hardware mode uses Ada Fruit RGB matrix together with their provide hat for raspberryPI.
To communicate with matrix it is using modified "github.com/mcuadros/go-rpi-rgb-led-matrix" (@Todo fork project to add modifications)

## Software (emulated)

The software emulated mode does not need hardware, it can run all animations on your browser. (Make sure you build binary with software flag.)
Start application and open http://localhost:8081 in your browser. You should see matrix rendered on canvas.

## Widgets

There are number of widgets which allows you to display certain information on matrix. Like weather, play gif etc.
The application exposes HTTP endpoint under "http://localhost:8081/exec?widget=weather" from which widgets can be invoked.

### Weather

Weather widget provides weather information for the configured location. 
The weather is updated every 15 min. The data comes from "https://openweathermap.org/api".

Trigger  http://localhost:8081/exec?widget=weather

Configuration "widget-weather"

| Field        | Value
|--------------|------------
| api-token    | The token from open weather API
| city         | The city for which to get weather data.

### GIF

The GIF widget displays animated picture in LED matrix display. Picture must be 32x32 pixels size.

No configuration required. To send image use http endpoint. http://localhost:8081/exec?widget=gif&url=http://domain.tld/animation.gif
It is also possible to set animation execution limit in seconds. Add addition GET parameter "&=duration=60". Will execute for one minute.
The widget periodically queries api if there are any torrents downloading, if any then display on screen information.

Trigger http://localhost:8081/exec?widget=torrent

### Torrent

The torrent widget integrates with uTorrent WEB API. http://help.utorrent.com/customer/portal/topics/664593/articles.
It will display downloaded torrent progress information. File name, remaining amount and remaining percentage. 
When download is finished screen will start flashing red. First you need to enable WEB API in uTorrent desktop client. 
Set user name and password. Then get local IP of the machine it is running. 

Configuration "widget-torrent-status"

| Field          | Value
|----------------|--------------------
| torrent-web-api-url | The url to WEB API, for example "http://192.168.1.12:8081/gui/"
| username  | The username of WEB API user
| password  | The password of WEB API user

### Meter

The meter display readings from CO2 meter device. https://www.amazon.com/dp/B00H7HFINS.
It shows CO2 and temperature of room where device is placed.

Trigger http://localhost:8081/exec?widget=meter

Configuration "widget-co2-meter" 

| Field         | Value
|---------------|-------------
| path-to-device | The path to device driver for example "/dev/hidraw0"
| warning-threshold | The CO2 level warning threshold, the integer value from what level to display alert


### Location

The location widget display distance between two points using GPS coordinates. 
To calculate distance google maps API is used.
 
Trigger http://localhost:8081/exec?widget=location

Configuration "widget-location"

| Field      | Value
|------------|---------------
| google-maps-token | The google maps API token
| stationary-location-gps-coordinates | The coordinates of LED matrix display, starting point.
| location-provider-url | The http endpoint which returns json response with two fields. See bellow 

Location provider response example http://www.domain.tld/geo.txt

Returns
``
{
  "lon": "13.4341874",
  "lat": "52.5396728"
}
``

### Explosion

The explosion widget plays simple explosion animation.  
  
Trigger http://localhost:8081/exec?widget=boom
  
