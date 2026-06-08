# Etapa 1: compila la aplicación Go.
# Usamos una imagen con Go porque aquí sí necesitamos las herramientas de compilación.
FROM golang:1.25-alpine AS builder

# Instala paquetes del sistema necesarios durante la compilación.
# `git` puede ser requerido por dependencias de Go.
# `ca-certificates` agrega certificados para conexiones seguras.
# `tzdata` incluye información de zonas horarias.
RUN apk add --no-cache git ca-certificates tzdata

# Define el directorio de trabajo dentro del contenedor.
# A partir de aquí, los comandos se ejecutan desde `/app`.
WORKDIR /app

# Copia primero los archivos de dependencias.
# Esto mejora la cache de Docker: si el código cambia pero las dependencias no,
# Docker reutiliza esta capa y no vuelve a descargarlas.
COPY go.mod go.sum ./

# Descarga los módulos definidos en go.mod/go.sum.
RUN go mod download

# Copia el resto del código fuente al contenedor.
COPY . .

# Compila la aplicación y genera el ejecutable `/app/main`.
# `CGO_ENABLED=0` intenta generar un binario estático.
# `GOOS=linux` asegura que el binario esté preparado para ejecutarse en Linux.
# `GOARCH=amd64` fija la arquitectura de salida.
# `-ldflags="-w -s"` reduce el tamaño del binario eliminando información de debug.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/main \
    ./cmd/api

# Etapa 2: crea la imagen final.
# Usamos una imagen distroless porque es más pequeña y más segura:
# solo contiene lo mínimo para ejecutar el binario.
FROM gcr.io/distroless/static:nonroot

# Directorio de trabajo de la imagen final.
WORKDIR /

# Copia únicamente el binario compilado desde la etapa anterior.
COPY --from=builder /app/main /main

# Copia la base de zonas horarias para que la app pueda manejar fechas/horas correctamente.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

# Ejecuta la aplicación con un usuario no root por seguridad.
USER nonroot:nonroot

# Documenta que la aplicación escucha en el puerto 8080.
EXPOSE 8080

# Comando que se ejecuta al iniciar el contenedor.
ENTRYPOINT ["/main"]
