# BROWSER environment variable

On non-Windows systems, Git Town will first read the `BROWSER` environment
variable to determine the browser command. If it isn't set, Git Town will try
various common commands like `open`, `xdg-open`, or `x-www-browser`.
