local plugin = require("kong.plugins.base_plugin"):extend()

function plugin:new()
  plugin.super.new(self, "request-transformer-advanced")
end

function plugin:access(plugin_conf)
  plugin.super.access(self)
end


plugin.PRIORITY = 1000

return plugin
