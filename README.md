# Convenciones del Proyecto

Este documento describe las convenciones de desarrollo que estamos siguiendo para el proyecto.

## 1. Estructura del Proyecto

Adoptamos una estructura de proyecto orientada a paquetes. Cada paquete encapsula una funcionalidad específica y relacionada.

### 1.1 Interfaces

Se definirán interfaces claras para los repositorios y servicios dentro de cada estructura de paquete. Esto promueve la modularidad y facilita la sustitución de implementaciones.

**Nomenclatura:**

- Para los servicios, la interfaz se nombrará `entidad_service.go` y su implementación `entidad_default.go`.
- Para los repositorios, la interfaz será `entidad_repository.go` y su implementación `entidad_map.go`.

### 1.2 Persistencia de Datos

Por el momento, el repositorio se manejará como un mapa en memoria, en lugar de una base de datos relacional. Esta decisión se basa en el contenido cubierto hasta ahora en la asignatura. En futuras etapas, se evaluará la integración con una base de datos relacional.

### 1.3 Gestión de Errores

La estructura de los errores se define en el archivo `errors.go` ubicado dentro del paquete `pkg`. Este archivo centraliza la gestión de errores para mantener la consistencia en todo el proyecto.

## 2. Flujo de Trabajo con Ramas

Cada integrante del equipo trabajará en una user story específica, generando una rama dedicada para el desarrollo de dicha funcionalidad. Una vez completado y revisado el trabajo, esta rama se fusionará con la rama principal (`main`).

## 3. Pruebas Unitarias

Se implementarán pruebas unitarias para asegurar la calidad y el correcto funcionamiento del código. Como mínimo, se testearán las capas de servicio y repositorio para verificar su comportamiento individual y aislado.
