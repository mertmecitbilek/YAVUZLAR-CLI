package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type User struct {
	Username string
	Password string
	Role     string // "admin" veya "customer"
}

// Kullanıcıların tutulduğu kısım burası
var users = []User{
	{"admin", "admin123", "admin"},
	{"customer", "cust123", "customer"},
}

// Log dosyasına mesaj kaydetme fonksiyonu
func logToFile(message string) {
	file, _ := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("%s: %s\n", timestamp, message)
	file.WriteString(logMessage)
}

// Kullanıcı doğrulama fonksiyonu
func authenticate(username, password string) *User {
	for i, user := range users {
		if user.Username == username && user.Password == password {
			return &users[i] // Kullanıcıyı referansla döndürüyoruz
		}
	}
	return nil
}

// Admin için oluşturulan kısım

func adminMenu() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nAdmin Menüsü:")
		fmt.Println("1. Müşteri Ekle")
		fmt.Println("2. Müşteri Sil")
		fmt.Println("3. Çıkış Yap")
		fmt.Print("Seçiminizi girin: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {

		case "1":
			fmt.Print("Yeni müşterinin kullanıcı adını girin: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			fmt.Print("Yeni müşterinin şifresini girin: ")
			password, _ := reader.ReadString('\n')
			password = strings.TrimSpace(password)

			users = append(users, User{Username: username, Password: password, Role: "customer"})
			fmt.Println("Müşteri başarıyla eklendi.")
			logToFile(fmt.Sprintf("Admin yeni müşteri ekledi: %s", username))

		case "2":
			fmt.Print("Silinecek müşterinin kullanıcı adını girin: ")
			username, _ := reader.ReadString('\n')
			username = strings.TrimSpace(username)

			found := false
			for i, user := range users {
				if user.Username == username && user.Role == "customer" {
					users = append(users[:i], users[i+1:]...)
					fmt.Println("Müşteri başarıyla silindi.")
					logToFile(fmt.Sprintf("Admin müşteri sildi: %s", username))
					found = true
					break
				}
			}
			if !found {
				fmt.Println("Müşteri bulunamadı.")
			}

		case "3":
			fmt.Println("Çıkış yapılıyor...")
			return

		default:
			fmt.Println("Geçersiz seçim. Lütfen tekrar deneyin.")
		}
	}
}

// Kullanıcılar için oluşturulan kısım

func customerMenu(user *User) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nMüşteri Menüsü:")
		fmt.Println("1. Profil Görüntüle")
		fmt.Println("2. Şifre Değiştir")
		fmt.Println("3. Çıkış Yap")
		fmt.Print("Seçiminizi girin: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			fmt.Printf("\nProfil:\nKullanıcı Adı: %s\n", user.Username)

		case "2":
			fmt.Print("Yeni şifrenizi girin: ")
			newPassword, _ := reader.ReadString('\n')
			newPassword = strings.TrimSpace(newPassword)

			// Kullanıcının şifresinin güncellendiği kısım
			user.Password = newPassword
			fmt.Println("Şifre başarıyla değiştirildi.")
			logToFile(fmt.Sprintf("Müşteri %s şifresini değiştirdi", user.Username))

		case "3":
			fmt.Println("Çıkış yapılıyor...")
			return

		default:
			fmt.Println("Geçersiz seçim. Lütfen tekrar deneyin.")
		}
	}
}

// İlk giriş kısmı burası

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nGiriş Sistemi")
		fmt.Print("Kullanıcı adınızı girin: ")
		username, _ := reader.ReadString('\n')
		username = strings.TrimSpace(username)
		fmt.Print("Şifrenizi girin: ")
		password, _ := reader.ReadString('\n')
		password = strings.TrimSpace(password)

		// Doğrulama yapılan kısım

		user := authenticate(username, password)
		if user == nil {
			fmt.Println("Hatalı giriş. Lütfen tekrar deneyin.")
			logToFile(fmt.Sprintf("Başarısız giriş denemesi: %s", username))
			continue
		}

		// log.txt'ye logların yazılması

		logToFile(fmt.Sprintf("Kullanıcı %s (%s) giriş yaptı", user.Username, user.Role))
		if user.Role == "admin" {
			adminMenu()
		} else {
			customerMenu(user)
		}
	}
}
