4. **Configurar .env**
```bash
cd proyecto
cp .env.example .env
# Editar con valores seguros
nano .env
# IMPORTANTE: Setear VITE_GO_API_URL con la IP pública ANTES de hacer build
# Ejemplo: VITE_GO_API_URL=http://54.XXX.XXX.XXX:8080
# Esto es necesario porque VITE_GO_API_URL se hornea en tiempo de compilación (build arg)
```
5. **Levantar servicios**
```bash
docker-compose up -d --build
```
