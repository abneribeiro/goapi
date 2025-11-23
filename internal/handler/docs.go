package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/abneribeiro/goapi/internal/model"
)

type DocsHandler struct {
	openapiPath string
}

func NewDocsHandler(docsPath string) *DocsHandler {
	return &DocsHandler{
		openapiPath: filepath.Join(docsPath, "openapi.yaml"),
	}
}

func (h *DocsHandler) ServeOpenAPI(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(h.openapiPath)
	if err != nil {
		// Certifique-se de que a função respondJSON existe neste pacote ou importada
		respondJSON(w, http.StatusNotFound, model.ErrorResponse("NOT_FOUND", "OpenAPI specification not found"))
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(data)
}

func (h *DocsHandler) ServeScalarUI(w http.ResponseWriter, r *http.Request) {
	// Nota: O uso de crases (backticks) permite strings multilinhas no Go.
	// A configuração JS foi ajustada para corrigir o erro de importCollection e ocultar as tools.
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Equipment Rental API - Interactive Documentation</title>
    <meta name="description" content="Interactive API documentation for the Equipment Rental platform.">
    <link rel="icon" type="image/svg+xml" href="data:image/svg+xml,<svg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 100 100'><text y='.9em' font-size='90'>🔧</text></svg>">
    <style>
        body { margin: 0; padding: 0; }
        /* Segurança extra para ocultar botão flutuante caso o JS demore */
        .scalar-api-client__toggle { display: none !important; }
    </style>
</head>
<body>
    <script id="api-reference" data-url="/docs/openapi.yaml"></script>

    <script>
        var configuration = {
            theme: 'purple',
            layout: 'modern',
            showSidebar: true,
            darkMode: true,
            searchHotKey: 'k',
            
            // --- CONFIGURAÇÕES DE VISIBILIDADE ---
            showDeveloperTools: false,   // Oculta a barra de ferramentas (mesmo em localhost)
            hideModels: false,
            hideDownloadButton: false,
            hideTestRequestButton: false, // Mantém o botão de teste (se quiser remover, mude para true)
            
            // --- CORREÇÃO DO ERRO DE CRASH ---
            hideClients: true,            // Oculta gerador de código (cURL, JS) que estava causando o erro "Missing required param"
            // ---------------------------------

            defaultHttpClient: {
                targetKey: 'shell',
                clientKey: 'curl',
            },
            authentication: {
                preferredSecurityScheme: 'bearerAuth',
            },
            metaData: {
                title: 'Equipment Rental API',
                description: 'A comprehensive RESTful API for equipment rental management.',
                ogDescription: 'Browse equipment, make reservations, and manage your rentals.',
                ogTitle: 'Equipment Rental API - Interactive Documentation',
                twitterCard: 'summary_large_image',
            },
            customCss: ` + "`" + `
                .darklight-reference-promo { display: none !important; }
                .scalar-app {
                    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
                }
                .sidebar {
                    background: linear-gradient(180deg, #1e1b4b 0%, #312e81 100%);
                }
                .sidebar-heading {
                    font-weight: 600;
                }
            ` + "`" + `,
            onSpecLoaded: function() {
                // Se tiver um elemento de loading, oculta aqui
                console.log('Spec Loaded');
            }
        };

        // --- PONTO CRÍTICO DE CORREÇÃO ---
        // Injeta a configuração no elemento ANTES do script do Scalar rodar
        var apiReference = document.getElementById('api-reference');
        apiReference.dataset.configuration = JSON.stringify(configuration);
    </script>

    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference@latest"></script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, max-age=3600")
	w.Write([]byte(html))
}
