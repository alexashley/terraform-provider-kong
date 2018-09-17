--- This is a modified copy of Kong's open source request-transformer plugin
--- See the README in this directory for license information.
--- Original: https://github.com/Kong/kong/blob/6acb86a67de50a7c50151b7350693ea9dc40e65d/kong/plugins/request-transformer/schema.lua#L1-L79

local find = string.find
-- entries must have colons to set the key and value apart
local function check_for_value(value)
    for i, entry in ipairs(value) do
        local ok = find(entry, ":")
        if not ok then
            return false, "key '" .. entry .. "' has no value"
        end
    end
    return true
end

local function check_method(value)
    if not value then
        return true
    end
    local method = value:upper()
    local ngx_method = ngx["HTTP_" .. method]
    if not ngx_method then
        return false, method .. " is not supported"
    end
    return true
end

return {
    fields = {
        http_method = { type = "string", func = check_method },
        remove = {
            type = "table",
            schema = {
                fields = {
                    body = { type = "array", default = {} }, -- does not need colons
                    headers = { type = "array", default = {} }, -- does not need colons
                    querystring = { type = "array", default = {} } -- does not need colons
                }
            }
        },
        rename = {
            type = "table",
            schema = {
                fields = {
                    body = { type = "array", default = {} },
                    headers = { type = "array", default = {} },
                    querystring = { type = "array", default = {} }
                }
            }
        },
        replace = {
            type = "table",
            schema = {
                fields = {
                    -- config.replace.uri
                    uri = {type = "string"},
                    body = { type = "array", default = {}, func = check_for_value },
                    headers = { type = "array", default = {}, func = check_for_value },
                    querystring = { type = "array", default = {}, func = check_for_value }
                }
            }
        },
        add = {
            type = "table",
            schema = {
                fields = {
                    body = { type = "array", default = {}, func = check_for_value },
                    headers = { type = "array", default = {}, func = check_for_value },
                    querystring = { type = "array", default = {}, func = check_for_value }
                }
            }
        },
        append = {
            type = "table",
            schema = {
                fields = {
                    body = { type = "array", default = {}, func = check_for_value },
                    headers = { type = "array", default = {}, func = check_for_value },
                    querystring = { type = "array", default = {}, func = check_for_value }
                }
            }
        }
    }
}
