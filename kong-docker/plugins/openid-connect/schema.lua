local function allowed_methods(methods)
    local allowed_auth_methods = {
        "password",
        "client_credentials",
        "authorization_code",
        "bearer",
        "introspection",
        "kong_oauth2",
        "refresh_token",
        "session"
    }
    for _, method in ipairs(methods) do
        local match = false
        for _, allowed_method in ipairs(allowed_auth_methods) do
            if method == allowed_method then
                match = true
            end
        end

        if not match then
            return false
        end
    end

    return true
end

return {
    fields = {
        anonymous = { type = "string" },
        auth_methods = { type = "array", func = allowed_methods },
        issuer = { type = "url", required = true },
        consumer_claim = { type = "array" },
        consumer_by = { type = "array" }
    }
}
