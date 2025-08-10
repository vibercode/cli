package mcp

// StartMCPServer inicia el servidor MCP usando el protocolo real
func StartMCPServer() error {
	// Crear e iniciar el protocolo MCP sin mensajes de UI
	// El protocolo MCP debe comunicarse SOLO via JSON-RPC sobre stdout
	protocol := NewMCPProtocol()

	// Iniciar el bucle de comunicación (sin mensajes informativos)
	return protocol.Start()
}
