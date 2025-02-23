package utils

func CalculateInterest(principal float64, annualRate float64, tenureMonths int64) float64 {
	// Menghitung bunga tahunan berdasarkan bunga sederhana
	tenureYears := float64(tenureMonths) / 12 // Konversi tenure dari bulan ke tahun
	interest := principal * (annualRate / 100) * tenureYears

	return interest
}
