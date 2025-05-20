# static-buildpack

A buildpack static websites, aka for:

* gohugo
* mdBook

---

Here's a high-level overview about how this works:

```mermaid
graph TD
    subgraph BUILD ["build phase"]
        direction LR
        version-->dl["download tool"]
    end

    subgraph DETECT ["detect phase"]
        direction LR
        w["check work dir"]-->t["determine type"]
    end

    DETECT-->d(("gohugo or mdBook"))
    d-->BUILD
    BUILD-->s["container with nginx/apache + your site"]
```

Some configuration is available, please see [api](./api/) for details.

## Dependencies

- `paketo-buildpacks/nginx`
- `paketo-buildpacks/httpd`

Either of these can be customized through various environment variables or a full config file for the web server. More details are available [on our documentation](https://www.runway.horse/docs/recipes/webservers/).

## Usage / License

Feel free to use this, as you see fit. For the turn-key zero-config solution, please check out [our PaaS service Runway](https://www.runway.horse/).