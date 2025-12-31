package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"golang.org/x/net/proxy"
)

const (
	TorProxyAddr  = "127.0.0.1:9050"
	TargetFile    = "targets.yaml"
	ScreenshotDir = "screenshots"
	ReportFile    = "scan_report.log"
)

type Site struct {
	URL  string
	Name string
}

func main() {

	os.MkdirAll(ScreenshotDir, 0755)

	logFile, _ := os.Create(ReportFile)
	defer logFile.Close()

	printInfo("Tor bağlantısı kontrol ediliyor (Port: 9050)...")

	if checkAndPrintIP() {
		printSuccess("Tor bağlantısı başarılı!")
	} else {
		printErr("Tor bağlantısı başarısız! Servis açık mı?")
		return
	}

	sites, err := readTargets(TargetFile)
	if err != nil {
		printErr(fmt.Sprintf("Dosya okunamadı: %v", err))
		return
	}
	printInfo(fmt.Sprintf("targets.yaml okundu - %d site yüklendi", len(sites)))

	for {
		fmt.Println()
		fmt.Println("=== Thor Scraper Menu ===")
		for i, s := range sites {
			fmt.Printf("%d. %s\n", i+1, s.Name)
		}
		fmt.Printf("%d. Hepsini Tara\n", len(sites)+1)
		fmt.Println("0. Çıkış")
		fmt.Print("Seçiminiz: ")

		var choice int
		fmt.Scan(&choice)

		if choice == 0 {
			break
		} else if choice == len(sites)+1 {

			for _, s := range sites {
				takeScreenshot(s, logFile)
			}
		} else if choice > 0 && choice <= len(sites) {

			takeScreenshot(sites[choice-1], logFile)
		} else {
			printErr("Geçersiz seçenek.")
		}
	}
}

func takeScreenshot(site Site, logFile *os.File) {
	printInfo(fmt.Sprintf("Taranıyor: %s (%s)", site.Name, site.URL))

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.ProxyServer("socks5://"+TorProxyAddr),
		chromedp.Flag("headless", true),
		chromedp.WindowSize(1280, 800),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	var buf []byte

	err := chromedp.Run(ctx,
		chromedp.Navigate(site.URL),
		chromedp.Sleep(5*time.Second),
		chromedp.FullScreenshot(&buf, 90),
	)

	if err != nil {
		printErr(fmt.Sprintf("%s -> UYARI: %v", site.Name, err))
		logFile.WriteString(fmt.Sprintf("[UYARI] %s (%s) -> %v\n", site.Name, site.URL, err))
		return
	}

	filename := fmt.Sprintf("%s/%s_%d.png", ScreenshotDir, site.Name, time.Now().Unix())
	if err := os.WriteFile(filename, buf, 0644); err != nil {
		printErr("Dosya kaydedilemedi")
		return
	}

	fileInfo, _ := os.Stat(filename)
	size := fileInfo.Size()
	printSuccess(fmt.Sprintf("Screenshot yakalandı %s (%d bytes)", site.Name, size))
	printSuccess(fmt.Sprintf("Screenshot kaydedildi: %s", filename))

	logFile.WriteString(fmt.Sprintf("[BAŞARILI] %s -> %s\n", site.Name, filename))
}

func checkAndPrintIP() bool {
	dialer, err := proxy.SOCKS5("tcp", TorProxyAddr, nil, proxy.Direct)
	if err != nil {
		printErr(fmt.Sprintf("Proxy hatası: %v", err))
		return false
	}

	client := &http.Client{
		Transport: &http.Transport{Dial: dialer.Dial},
		Timeout:   20 * time.Second,
	}

	resp, err := client.Get("http://check.torproject.org/api/ip")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	fmt.Printf("IP Adresi Doğrulanıyor... Tor IP Adresiniz: %s\n", strings.TrimSpace(string(body)))
	return true
}

func readTargets(path string) ([]Site, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var sites []Site
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && strings.Contains(line, "|") {
			parts := strings.Split(line, "|")
			sites = append(sites, Site{URL: parts[0], Name: parts[1]})
		}
	}
	return sites, scanner.Err()
}

func printInfo(msg string) {
	fmt.Printf("[BİLGİ] %s\n", msg)
}
func printSuccess(msg string) {
	fmt.Printf("[BAŞARILI] %s\n", msg)
}
func printErr(msg string) {
	fmt.Printf("[HATA] %s\n", msg)
}
