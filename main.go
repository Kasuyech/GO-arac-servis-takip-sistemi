package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/microsoft/go-mssqldb"
)

// Veri Modelleri
type Arac struct {
	ID    int
	Plaka string
	Marka string
	Model string
	Sahip string
}

func main() {
	// Bağlantı Ayarları
	connString := "server=localhost;port=1433;database=go_proje;integrated security=true;encrypt=disable;"

	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Bağlantı ayarı hatalı:", err)
	}
	defer db.Close()

	// Tabloları kontrol et ve yoksa oluştur
	tablolariHazirla(db)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n=====================================")
		fmt.Println("   🔧 ARAÇ SERVİS TAKİP SİSTEMİ 🔧   ")
		fmt.Println("=====================================")
		fmt.Println("1. Yeni Araç Kayıt Et")
		fmt.Println("2. Kayıtlı Araçları Listele")
		fmt.Println("3. Arıza / Servis Kaydı Ekle")
		fmt.Println("4. Aracın Servis Geçmişini Sorgula")
		fmt.Println("5. Toplam Gelir Raporu")
		fmt.Println("6. Çıkış")
		fmt.Print("Seçiminiz: ")

		secim, _ := reader.ReadString('\n')
		secim = strings.TrimSpace(secim)

		switch secim {
		case "1":
			aracEkle(db, reader)
		case "2":
			aracListele(db)
		case "3":
			servisKaydiEkle(db, reader)
		case "4":
			servisGecmisiSorgula(db, reader)
		case "5":
			toplamGelirRaporu(db)
		case "6":
			fmt.Println("Sistemden çıkılıyor. İyi günler!")
			os.Exit(0)
		default:
			fmt.Println("Geçersiz seçim, lütfen 1-6 arası bir değer girin!")
		}
	}
}

func tablolariHazirla(db *sql.DB) {
	sorgu := `
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='Araclar' AND xtype='U')
	CREATE TABLE Araclar (
		id INT PRIMARY KEY IDENTITY(1,1),
		plaka NVARCHAR(20) NOT NULL UNIQUE,
		marka NVARCHAR(50),
		model NVARCHAR(50),
		sahip NVARCHAR(100)
	);
	
	IF NOT EXISTS (SELECT * FROM sysobjects WHERE name='ServisKayitlari' AND xtype='U')
	CREATE TABLE ServisKayitlari (
		id INT PRIMARY KEY IDENTITY(1,1),
		arac_id INT FOREIGN KEY REFERENCES Araclar(id),
		islem NVARCHAR(MAX),
		teknisyen NVARCHAR(50),
		tutar DECIMAL(10,2),
		tarih DATETIME DEFAULT GETDATE()
	);`

	_, err := db.Exec(sorgu)
	if err != nil {
		log.Fatal("Tablolar oluşturulurken hata çıktı:", err)
	}
}

func aracEkle(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n--- Yeni Araç Girişi ---")
	fmt.Print("Plaka: ")
	plaka, _ := reader.ReadString('\n')

	fmt.Print("Marka: ")
	marka, _ := reader.ReadString('\n')

	fmt.Print("Model: ")
	model, _ := reader.ReadString('\n')

	fmt.Print("Sahibi: ")
	sahip, _ := reader.ReadString('\n')

	sorgu := "INSERT INTO Araclar (plaka, marka, model, sahip) VALUES (@p1, @p2, @p3, @p4)"
	_, err := db.Exec(sorgu, strings.TrimSpace(plaka), strings.TrimSpace(marka), strings.TrimSpace(model), strings.TrimSpace(sahip))

	if err != nil {
		fmt.Println("❌ Hata: Araç eklenemedi! (Plaka zaten kayıtlı olabilir)")
		return
	}
	fmt.Println("✅ Araç başarıyla kaydedildi!")
}

func aracListele(db *sql.DB) {
	rows, err := db.Query("SELECT id, plaka, marka, model, sahip FROM Araclar")
	if err != nil {
		log.Println("Listeleme hatası:", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n--- Sistemdeki Tüm Araçlar ---")
	for rows.Next() {
		var a Arac
		rows.Scan(&a.ID, &a.Plaka, &a.Marka, &a.Model, &a.Sahip)
		fmt.Printf("🚗 Plaka: %s | Marka/Model: %s %s | Sahibi: %s\n", a.Plaka, a.Marka, a.Model, a.Sahip)
	}
}

func servisKaydiEkle(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n--- Servis Kaydı Oluşturma ---")
	fmt.Print("İşlem yapılacak aracın plakası: ")
	plaka, _ := reader.ReadString('\n')
	plaka = strings.TrimSpace(plaka)

	// Önce plakadan aracın ID'sini bulalım
	var aracID int
	err := db.QueryRow("SELECT id FROM Araclar WHERE plaka = @p1", plaka).Scan(&aracID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("❌ Sistemde böyle bir plaka bulunamadı. Önce aracı kaydetmelisiniz.")
		} else {
			fmt.Println("❌ Veri tabanı hatası:", err)
		}
		return
	}

	fmt.Print("Araçtaki Şikayet / Yapılan İşlem: ")
	islem, _ := reader.ReadString('\n')

	fmt.Print("İlgilenen Teknisyen: ")
	teknisyen, _ := reader.ReadString('\n')

	fmt.Print("Toplam Maliyet (TL): ")
	tutarStr, _ := reader.ReadString('\n')
	tutar, err := strconv.ParseFloat(strings.TrimSpace(tutarStr), 64)
	if err != nil {
		fmt.Println("❌ Hatalı tutar girdiniz! Sadece rakam kullanın.")
		return
	}

	sorgu := "INSERT INTO ServisKayitlari (arac_id, islem, teknisyen, tutar) VALUES (@p1, @p2, @p3, @p4)"
	_, err = db.Exec(sorgu, aracID, strings.TrimSpace(islem), strings.TrimSpace(teknisyen), tutar)

	if err != nil {
		fmt.Println("❌ Kayıt başarısız oldu:", err)
		return
	}
	fmt.Println("✅ Servis kaydı başarıyla eklendi!")
}

func servisGecmisiSorgula(db *sql.DB, reader *bufio.Reader) {
	fmt.Println("\n--- Servis Geçmişi Sorgulama ---")
	fmt.Print("Sorgulanacak Plaka: ")
	plaka, _ := reader.ReadString('\n')
	plaka = strings.TrimSpace(plaka)

	// İki tabloyu (Araclar ve ServisKayitlari) birleştiren JOIN sorgusu
	sorgu := `
		SELECT s.tarih, s.islem, s.teknisyen, s.tutar 
		FROM ServisKayitlari s
		INNER JOIN Araclar a ON s.arac_id = a.id
		WHERE a.plaka = @p1
		ORDER BY s.tarih DESC
	`

	rows, err := db.Query(sorgu, plaka)
	if err != nil {
		fmt.Println("❌ Sorgulama hatası:", err)
		return
	}
	defer rows.Close()

	kayitVarMi := false
	fmt.Printf("\n>>> %s Plakalı Aracın Servis Kayıtları <<<\n", plaka)
	for rows.Next() {
		kayitVarMi = true
		var tarih, islem, teknisyen string
		var tutar float64

		rows.Scan(&tarih, &islem, &teknisyen, &tutar)
		// Tarih formatını kısaltmak için ilk 19 karakteri (YYYY-MM-DD HH:MM:SS) alıyoruz
		if len(tarih) > 19 {
			tarih = tarih[:19]
		}

		fmt.Printf("📅 Tarih: %s\n", tarih)
		fmt.Printf("🛠 İşlem: %s\n", islem)
		fmt.Printf("👨‍🔧 Teknisyen: %s\n", teknisyen)
		fmt.Printf("💰 Tutar: %.2f TL\n", tutar)
		fmt.Println("---------------------------------")
	}

	if !kayitVarMi {
		fmt.Println("Bu araca ait geçmiş servis kaydı bulunmamaktadır.")
	}
}

func toplamGelirRaporu(db *sql.DB) {
	var toplamGelir sql.NullFloat64 // Null gelme ihtimaline karşı güvenli tip
	var toplamIslem int

	err := db.QueryRow("SELECT SUM(tutar), COUNT(id) FROM ServisKayitlari").Scan(&toplamGelir, &toplamIslem)
	if err != nil {
		fmt.Println("❌ Rapor hesaplanamadı:", err)
		return
	}

	fmt.Println("\n📊 --- SERVİS GENEL DURUM RAPORU ---")
	fmt.Printf("Bugüne kadar yapılan toplam işlem sayısı: %d\n", toplamIslem)
	if toplamGelir.Valid {
		fmt.Printf("Kasadaki Toplam Gelir: %.2f TL\n", toplamGelir.Float64)
	} else {
		fmt.Println("Kasadaki Toplam Gelir: 0.00 TL")
	}
	fmt.Println("------------------------------------")
}
