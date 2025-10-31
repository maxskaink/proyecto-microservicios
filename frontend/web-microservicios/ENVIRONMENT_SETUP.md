# Configuración de Environment

## 🔧 Setup inicial

1. Copia el archivo de ejemplo:
```bash
cp src/app/environment/environment.example.ts src/app/environment/environment.ts
```

2. Edita `environment.ts` con tus configuraciones reales de Firebase:
   - Obtén las credenciales desde Firebase Console
   - Actualiza `apiUrl` con la URL de tu backend

## 🔒 Seguridad

- ❌ **NUNCA** commitees `environment.ts` con credenciales reales
- ✅ Solo commitea `environment.example.ts` con valores de ejemplo
- 🔑 Mantén las credenciales de Firebase seguras

## 🚀 Despliegue

Para producción, configura las variables de entorno en tu plataforma de hosting (Vercel, Netlify, etc.)
