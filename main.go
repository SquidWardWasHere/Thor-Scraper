package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/proxy"
)

const (
	TorProxyAddr = "127.0.0.1:9050"
	TargetFile   = "targets.yaml"
	OutputDir    = "target_data"
	ReportFile   = "tarama_raporu.log"
)

func main() {
	fmt.Println("=== Thor'un Scraperi Başlatılıyor ===")

	os.MkdirAll(OutputDir, 0755)

	client, err := torBaglantisiKur()
	if err != nil {
		log.Fatal("Tor bağlantı hatası:", err)
	}
	fmt.Println("Bağlantı başarılı")
	fmt.Println("[BİLGİ] Tor ağına bağlanıldı")

	ipAdresiniGoster(client)

	urls, err := dosyaOku(TargetFile)
	if err != nil {
		log.Fatal("targets.yaml dosyası bulunamadı:", err)
	}
	fmt.Printf("[BİLGİ] %d adet hedef yüklendi.\n", len(urls))

	logDosyasi, _ := os.OpenFile(ReportFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer logDosyasi.Close()

	for _, url := range urls {
		fmt.Printf("Taranıyor: %s ... ", url)

		resp, err := client.Get(url)
		if err != nil {
			msg := fmt.Sprintf("[HATA] %v\n", err)
			fmt.Print(msg)
			logDosyasi.WriteString(url + " -> " + msg)
			continue
		}

		if resp.StatusCode == 200 {
			dosyaAdi := fmt.Sprintf("%s/veri_%d.html", OutputDir, time.Now().Unix())
			dosyaKaydet(dosyaAdi, resp.Body)

			msg := fmt.Sprintf("[BAŞARILI] Kaydedildi: %s\n", dosyaAdi)
			fmt.Print(msg)
			logDosyasi.WriteString(url + " -> " + msg)
		} else {
			msg := fmt.Sprintf("[UYARI] Kod: %d\n", resp.StatusCode)
			fmt.Print(msg)
			logDosyasi.WriteString(url + " -> " + msg)
		}
		resp.Body.Close()
	}
	fmt.Println("\n=== İşlem Tamamlandı ===")
}

func torBaglantisiKur() (*http.Client, error) {
	dialer, err := proxy.SOCKS5("tcp", TorProxyAddr, nil, proxy.Direct)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &http.Transport{Dial: dialer.Dial},
		Timeout:   45 * time.Second,
	}, nil
}

func ipAdresiniGoster(client *http.Client) {
	fmt.Print("IP Adresi Doğrulanıyor... ")
	resp, err := client.Get("http://check.torproject.org/api/ip")
	if err != nil {
		fmt.Println("[HATA] IP kontrolü başarısız!")
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Tor IP Adresiniz: %s\n", string(body))
}

func dosyaOku(yol string) ([]string, error) {
	file, err := os.Open(yol)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			urls = append(urls, line)
		}
	}
	return urls, scanner.Err()
}

func dosyaKaydet(isim string, veri io.ReadCloser) {
	out, _ := os.Create(isim)
	defer out.Close()
	io.Copy(out, veri)
}
