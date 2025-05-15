# Usa una imagen con herramientas necesarias para compilar C
FROM golang:1.24

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos de definición de dependencias primero (para aprovechar la cache)
COPY go.mod ./
COPY go.sum ./

# Descarga las dependencias
RUN go mod download

# Copia el resto del código fuente
COPY . .

# Compila la aplicación
RUN go build -o todo-api

# Expone el puerto que usa tu servidor
EXPOSE 8080

# Comando por defecto para ejecutar la app
CMD ["./todo-api"]
