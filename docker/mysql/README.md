# MySQL Dockerizado para Mercado Fresco

Este directorio contiene la configuración necesaria para ejecutar la base de datos MySQL en un contenedor Docker.

## Configuración

La base de datos está configurada con los siguientes parámetros:

- **Usuario**: root
- **Contraseña**: ""
- **Base de datos**: frescos
- **Puerto**: 3306
- **Host**: localhost (cuando se accede desde la máquina host)
- **Host**: mysql (cuando se accede desde otros contenedores en la misma red)

## Cómo usar

### Iniciar la base de datos

Desde la raíz del proyecto, ejecuta:

```bash
docker-compose up --build
```

Esto construirá la imagen de Docker y ejecutará el contenedor en segundo plano.

### Verificar el estado

Para verificar que el contenedor está funcionando:

```bash
docker-compose ps
```

### Detener la base de datos

```bash
docker-compose down
```

### Acceder a la base de datos

Para conectarte a la base de datos desde la línea de comandos:

```bash
docker exec -it mercado_fresco_mysql mysql -uroot -p"" frescos
```

## Configuración de la aplicación

Para conectar tu aplicación a esta base de datos dockerizada, asegúrate de que tu archivo `config.yml` tenga la siguiente configuración:

```yaml
database:
  user: root
  password: ""
  host: mysql
  port: 3306
  name: frescos
```

## Datos iniciales

El contenedor carga automáticamente el script `database_setup_with_seed.sql` que crea todas las tablas y carga datos de prueba.
