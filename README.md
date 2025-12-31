# 🚀 Thor-Scraper 
Bu proje, **Go** dili ve **Chromedp** altyapısını kullanarak web sayfalarını otomatik tarayan ve ekran görüntüsü alan bir tooldur.

* 🌐 **Headless Browser:** Chrome altyapısıyla görünmez ve hızlı tarayıcı yönetimi.
* 📸 **Smart Screenshot:** Sayfanın tamamını otomatik yakalama.
* ⚡  **Yüksek Performans:** Go'nun eşzamanlılık (concurrency) gücüyle hızlı işlem.


## 🛠️ Kurulum ve Kullanım

Aracı çalıştırmak için aşağıdaki adımları sırasıyla uygulamanız yeterlidir:
1️⃣ Ön Hazırlık
Sisteminizde **Go** ve **Google Chrome** yüklü olduğundan emin olun.
Bu projede tor ağı kullanılacağı için;

**Tor servisi kurulumu **
```bash
A) Linux (Kali / Ubuntu / Debian) Kurulumu
# Tor servisini yükle
sudo apt update && sudo apt install tor -y
# Servisi başlat
sudo service tor start
# Çalışıp çalışmadığını kontrol et (Active: active (running) görmelisin)
sudo service tor status

B) macOS (Homebrew ile) Kurulumu
# Tor servisini yükle
brew install tor
# Servisi başlat
brew services start tor

2️⃣ **Depoyu bilgisayarınıza indirin:**
   ```bash
   git clone https://github.com/SquidWardWasHere/Thor-Scraper.git
   cd Thor-Scraper

Kullanım için;
cd Thor-Scraper

3️⃣ Bağımlılıkları Çekin
Gerekli tüm kütüphaneleri (chromedp, sysutil, pdf vb.) yüklemek için:
go mod tidy

4️⃣ Çalıştırın
Her şey hazır! Uygulamayı başlatmak için:
go run main.go

| Kütüphane | Görevi |
| :--- | :--- |
| **chromedp** | Tarayıcı kontrolü ve otomasyon. |
| **pixelmatch** | Görüntü karşılaştırma ve fark bulma. |
| **ledongthuc/pdf** | PDF okuma ve işleme desteği. |
| **sysutil** | Sistem düzeyinde yardımcı fonksiyonlar. |
