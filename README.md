# Weather Forecast REST API
Author of this app is Jakub Klonowski (jakubpklonowski@gmail.com).

Data is generated randomly. No connection to any data source.

## Manual
To host API use `build_win.ps1` script. This will generate `weatherForecast.exe` file that you need to run. Credentials for this api are: 
- login: `user`
- password: `password`


## Endpoints
Test

GET `http://127.0.0.1:10000/api/test/{data}`

Login

POST `http://127.0.0.1:10000/api/weather/login`

Weather

POST `http://127.0.0.1:10000/api/weather`

## Structs
### Login
Request

    {
        "login": "",
        "password": ""
    }

Response

    {
        "token": ""
    }

### Weather
Request

    {
        "token": "",
        "geo": {
            "lat": 0,
            "long": 0
        },
        "unit": ""
    }

Response

    {
        "weather": {
            "weather": "",
            "temperature": 0,
            "humidity": 0
        }
    }
