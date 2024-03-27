# jack-henry-eval-weather
Evaluation exercise for a position with Jack Henry

Configuration options below:

| Flag | Env Var | Default | Description |
| -----| -------- | -------- | ----------- |
| -listen | WEATHER_LISTEN | :8080 | Port for the service to listen on |
| -openweathermap-key | WEATHER_OPENWEATHERMAP_KEY | | API key for Open Weather Map | 

Architecture:
- Endpoints can be added as individual modules in the handlers directory.
- Configuration is shared among all modules through the config module.
- The main executable is minimal to ensure that functionality is in the handlers
