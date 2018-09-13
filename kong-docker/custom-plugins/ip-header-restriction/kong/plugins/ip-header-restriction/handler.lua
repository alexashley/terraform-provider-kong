local plugin = require("kong.plugins.base_plugin"):extend()

function plugin:new()
  plugin.super.new(self, "ip-header-restriction")
end

function plugin:access(plugin_conf)
  plugin.super.access(self)

  ngx.req.set_header("X-Ip-Header-Restriction-Dummy", "Hi from plugin")
end


plugin.PRIORITY = 1000

return plugin
