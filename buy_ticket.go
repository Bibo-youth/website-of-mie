package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
	"time"
)

// 模拟登录
func login(client *resty.Client, username, password string) error {
	loginURL := "https://example.com/login" // 替换为实际登录的 URL
	payload := map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := client.R().
		SetFormData(payload).
		Post(loginURL)
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("login failed with status %d", resp.StatusCode())
	}

	// 检查是否登录成功
	document, err := goquery.NewDocumentFromReader(resp.RawResponse.Body)
	if err != nil {
		return err
	}
	if document.Find(".login-success").Length() == 0 {
		return fmt.Errorf("login failed")
	}

	return nil
}

// 填写预约信息并提交
func submitReservation(reservationURL string, data map[string]string) error {
	//resp, err := client.R().
	//	SetFormData(data).   // 填写表单信息
	//	Post(reservationURL) // 提交预约请求
	//if err != nil {
	//	return err
	//}
	//
	//if resp.StatusCode() != 200 {
	//	return fmt.Errorf("failed to submit reservation, status code: %d", resp.StatusCode())
	//}
	//
	//// 解析返回的 HTML 判断是否提交成功
	//document, err := goquery.NewDocumentFromReader(resp.RawResponse.Body)
	//if err != nil {
	//	return err
	//}
	//
	//// 根据返回的 HTML 判断是否预约成功（此处假设通过某个标识来判断）
	//if document.Find(".reservation-success").Length() == 0 {
	//	return fmt.Errorf("reservation submission failed")
	//}
	//
	//fmt.Println("Reservation submitted successfully!")

	jsonPayload, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Error marshaling JSON for callback: %v", err)
	}

	req, err := http.NewRequest("POST", reservationURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error creating callback request: %v", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making callback request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: server returned status code %v", resp.StatusCode)
		return err
	}
	fmt.Printf("response header %v", resp.Header)
	fmt.Printf("response body %v", resp.Body)
	return nil
}

func main() {
	// 登录信息
	username := "your-username"
	id := "1234568"

	//// 登录
	//err := login(client, username, id)
	//if err != nil {
	//	log.Fatalf("Login failed: %v", err)
	//}

	// 预约信息
	reservationURL := "https://eapply.abchina.com/coin/coin/CoinIssuesDistribution?typeid=202307"
	reservationData := map[string]string{
		"name":  username,
		"id":    id,
		"phone": "1234567890",
		"date":  "2024-12-25",
	}

	// 启动多个 goroutine 来并发提交预约请求
	for i := 0; i < 5; i++ {
		go func() {
			err := submitReservation(reservationURL, reservationData)
			if err != nil {
				fmt.Printf("Failed to submit reservation: %v", err)
			}
		}()
	}

	// 阻塞主线程，确保所有 goroutine 执行完毕
	time.Sleep(5 * time.Second)
}
