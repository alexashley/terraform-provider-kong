return {
    no_consumer = true,
    fields = {
        whitelist = { type = 'array', required = true },
        ip_header = { type = 'string', required = true, default = 'x-forwarded-for' },
        override_global = { type = 'bool', required = true, default = false }
    }
}
