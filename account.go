/**
 * @Author: ZhaoYadong
 * @Date: 2020-11-18 10:53:59
 * @LastEditors: ZhaoYadong
 * @LastEditTime: 2020-11-18 10:56:26
 * @FilePath: /src/wxpay/account.go
 */
package wxpay

import (
	"io/ioutil"
	"log"
)

type Account struct {
	appID     string
	mchID     string
	apiKey    string
	certData  []byte
	isSandbox bool
	isPem     bool
}

// 创建微信支付账号
func NewAccount(appID string, mchID string, apiKey string, isSanbox bool) *Account {
	return &Account{
		appID:     appID,
		mchID:     mchID,
		apiKey:    apiKey,
		isSandbox: isSanbox,
	}
}

// 设置证书
func (a *Account) SetCertData(certPath string) {
	certData, err := ioutil.ReadFile(certPath)
	if err != nil {
		log.Println("读取证书失败")
		return
	}
	a.certData = certData
}

// 设置证书类型
func (a *Account) SetIsPem() {
	a.isPem = true
}
