-- schema.lua
-- Definición del esquema de configuración del plugin

local typedefs = require "kong.db.schema.typedefs"

return {
  name = "tenant-extractor",
  fields = {
    {
      config = {
        type = "record",
        fields = {
          {
            header_name = {
              type = "string",
              default = "X-Tenant-ID",
              required = true,
            }
          },
        },
      },
    },
  },
}
