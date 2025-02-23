package utils

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateAlphanumericCode(length int) string {
	// Definisikan karakter yang bisa digunakan dalam kode
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// Buat generator acak dengan seed waktu
	rand.Seed(time.Now().UnixNano())

	// Buat array untuk menyimpan hasil kode
	var code strings.Builder

	// Pilih karakter acak sesuai dengan panjang yang diminta
	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(chars)) // Pilih indeks acak dari karakter
		code.WriteByte(chars[randomIndex])   // Tambahkan karakter ke hasil kode
	}

	// Kembalikan hasil kode sebagai string
	return code.String()
}
