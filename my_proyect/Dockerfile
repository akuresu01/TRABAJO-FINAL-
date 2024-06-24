# Usa la imagen base de golang
FROM golang:1.20-alpine

# Establece el directorio de trabajo dentro del contenedor
WORKDIR /app

# Copia los archivos go.mod y descarga las dependencias
COPY go.mod ./
RUN go mod download

# Copia el código fuente al directorio de trabajo
COPY . .

# Construye la aplicación
RUN go build -o /myapp

# Expone el puerto en el que la aplicación escuchará
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["/myapp"]
