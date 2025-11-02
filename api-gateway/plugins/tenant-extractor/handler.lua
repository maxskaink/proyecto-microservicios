-- handler.lua
-- Plugin para extracción de tenant desde URL y agregarlo como header

local TenantExtractorHandler = {
  VERSION = "1.0.0",
  PRIORITY = 1000,
}

function TenantExtractorHandler:access(conf)
  local tenant_id = nil
  
  -- Obtener el path de la request
  local request_path = kong.request.get_path()
  
  -- Extraer el tenant desde el primer segmento del path usando pattern matching
  -- Busca un patrón como: /tenant_name/api/...
  local pattern = "^/([a-zA-Z0-9_%-]+)/api/"
  tenant_id = string.match(request_path, pattern)
  
  if tenant_id then
    -- Si se encontró un tenant, agregarlo como header
    kong.service.request.set_header(conf.header_name or "X-Tenant-ID", tenant_id)
    
    -- Log para debugging
    kong.log.debug("Tenant extracted: ", tenant_id)
    
    -- Remover el tenant del path antes de enviarlo al servicio
    -- Transforma /tenant/api/users/me en /api/users/me
    local new_path = string.gsub(request_path, "^/" .. tenant_id, "", 1)
    kong.service.request.set_path(new_path)
    
    kong.log.debug("Original path: ", request_path)
    kong.log.debug("New path: ", new_path)
  else
    -- Si no hay tenant en la URL, asegurarse de que no exista el header
    kong.service.request.clear_header(conf.header_name or "X-Tenant-ID")
    kong.log.debug("No tenant found in path")
  end
end

return TenantExtractorHandler
