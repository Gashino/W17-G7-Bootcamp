FROM mysql:8.0

# Copiar el script SQL al directorio docker-entrypoint-initdb.d
# Los scripts en este directorio se ejecutan automáticamente al iniciar el contenedor
COPY ./docs/sql/database_setup_with_seed.sql /docker-entrypoint-initdb.d/

# Variables de entorno para configurar MySQL
ENV MYSQL_ALLOW_EMPTY_PASSWORD=yes
ENV MYSQL_DATABASE=frescos

# Exponer el puerto 3306
EXPOSE 3306

# Configuración para mejorar el rendimiento y la compatibilidad
CMD ["mysqld", "--character-set-server=utf8mb4", "--collation-server=utf8mb4_unicode_ci", "--skip-log-bin"]
