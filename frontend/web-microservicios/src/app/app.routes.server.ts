import { RenderMode, ServerRoute } from '@angular/ssr';

export const serverRoutes: ServerRoute[] = [
  {
    path: 'login', // ← Sin barra inicial para coincidir con app.routes.ts
    renderMode: RenderMode.Prerender
  },
  {
    path: 'home',
    renderMode: RenderMode.Prerender
  },

  {
    path: 'user',
    renderMode: RenderMode.Server // Contenido dinámico para el usuario
  },
  {
    path: 'publishProduct',
    renderMode: RenderMode.Server // Contenido dinámico para el usuario
  },
  {
    path: '', // ← Ruta raíz también sin barra
    renderMode: RenderMode.Prerender
  },
  {
    path: '**', // ← Catch-all para otras rutas
    renderMode: RenderMode.Server
  }
];
