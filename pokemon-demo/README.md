# Pokémon App: inicio rápido con Docker

Esta carpeta permite ejecutar la app y PostgreSQL en otra computadora con Docker Desktop. El código de la aplicación se toma de la carpeta raíz del repositorio, así que hay que clonar o actualizar el repositorio completo.

## Requisitos

- Docker Desktop instalado y en ejecución.
- Git, para clonar el repositorio o traer los cambios más recientes.

## Primera instalación

En PowerShell:

```powershell
git clone https://github.com/Interino2/Pokemons.git
cd Pokemons\pokemon-demo
Copy-Item .env.example .env
docker compose up -d --build
```

Abre <http://localhost:8080>.

## Actualizar una copia existente

Desde la carpeta del repositorio:

```powershell
git pull
cd pokemon-demo
docker compose up -d --build
```

`git pull` descarga el código, pero no reconstruye ni reinicia los contenedores. `docker compose up -d --build` construye la app con el código actualizado y la inicia. La base de datos conserva sus datos en un volumen de Docker.

## Comandos útiles

Ejecutar desde `pokemon-demo`:

```powershell
docker compose ps
docker compose logs -f app
docker compose down
```

`docker compose down` detiene y elimina los contenedores, pero conserva los datos. **No uses `docker compose down -v` salvo que quieras borrar también el volumen y todos los datos de PostgreSQL.**

Si los puertos `8080` o `5434` ya están ocupados, cambia `APP_HOST_PORT` o `DB_HOST_PORT` en el archivo `.env` antes de iniciar los servicios.
