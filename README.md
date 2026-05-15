# Araç Servis Takip Sistemi (Go & MSSQL)

Bu proje, Go (Golang) ve MSSQL kullanılarak geliştirilmiş, konsol üzerinden çalışan bir araç servis takip ve kayıt sistemidir. 

## 🛠️ Kurulum ve Çalıştırma Gereksinimleri

Projenin başka bir bilgisayarda sorunsuz çalışması için aşağıdaki adımların izlenmesi gerekmektedir:

### 1. Veri Tabanı Hazırlığı
Bu proje verileri Microsoft SQL Server (MSSQL) üzerinde tutmaktadır. Projeyi çalıştırmadan önce:
* SQL Server Management Studio (SSMS) üzerinden **`go_proje`** adında boş bir veri tabanı oluşturulmalıdır.
* Proje kodları, veri tabanındaki tabloları (`Araclar` ve `ServisKayitlari`) ilk çalışmada **otomatik olarak** oluşturacaktır.

### 2. SQL Server Ağ Ayarları
Go'nun MSSQL'e bağlanabilmesi için TCP/IP protokolünün açık olması zorunludur:
* **SQL Server Configuration Manager**'ı açın.
* `SQL Server Network Configuration` -> `Protocols for MSSQLSERVER` (veya SQLEXPRESS) yolunu izleyin.
* **TCP/IP** durumunu **Enabled** (Etkin) yapın.
* Değişikliğin geçerli olması için SQL Server servisini **yeniden başlatın (Restart)**.

### 3. Kimlik Doğrulama
Projedeki bağlantı metni (Connection String) `integrated security=true` olarak ayarlanmıştır. Bu nedenle SSMS'e girerken **Windows Authentication** kullanılarak giriş yapılmalıdır.

### 4. Projeyi Çalıştırma
Projeyi bilgisayarınıza klonladıktan sonra terminal üzerinden proje dizinine giderek aşağıdaki komutu çalıştırmanız yeterlidir. Gerekli Go paketleri (MSSQL sürücüsü) otomatik olarak inecektir:

```bash
go run main.go
