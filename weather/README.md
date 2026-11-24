# Waybar weather indicator

Weather for Waybar, powered by [wttr.in](https://github.com/chubin/wttr.in).
Shows the current and feels-like temperatures, with a detailed forecast on hover.

## Setup

Install Go and Task, then run:

```sh
task prod
```

Add this definition to your Waybar configuration and include `custom/weather`
in your module list. It updates every 30 minutes:

```jsonc
"custom/weather": {
    "max-length": 15,
    "return-type": "json",
    "format": "WX:{}",
    "exec": "$HOME/.config/waybar/scripts/weather",
    "interval": 1800
}
```

Location is detected from your IP address by default. To choose a city, prefix
the `exec` command with `CITY_WEATHER=London`.

Forecast markers: `r` means rain, `s` means snow.

## Fallback

Results are saved automatically. If an update fails, the indicator uses the last
successful reading for the same city, keeping its original timestamp. After the
timestamp, the tooltip shows `Fallback not used` for fresh results and `Fallback used`
when a saved reading is used. If no saved reading exists, the command logs an error
and exits.

## Preview

![preview](/assets/preview.png)
