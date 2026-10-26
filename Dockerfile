# Menggunakan image resmi Golang versi Alpine sebagai tahap builder
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Menyalin file manajemen modul Go
COPY go.mod ./
RUN go mod download

# Menyalin seluruh kode sumber ke dalam container
COPY . .

# Melakukan build file binary aplikasi dengan nama 'main'
RUN go build -o main .

# Menggunakan image Alpine yang sangat ringkas untuk tahap production
FROM alpine:latest

WORKDIR /app

# Menyalin hasil build binary dari tahap builder sebelumnya
COPY --from=builder /app/main .

# Membuka port 8080 agar aplikasi dapat diakses dari luar container
EXPOSE 8080

# Menjalankan file binary aplikasi saat container aktif
CMD ["./main"]