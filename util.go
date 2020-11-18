/**
 * @Author: ZhaoYadong
 * @Date: 2020-11-18 10:53:59
 * @LastEditors: ZhaoYadong
 * @LastEditTime: 2020-11-18 11:35:28
 * @FilePath: /src/wxpay/util.go
 */
package wxpay

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"log"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/pkcs12"
)

func XmlToMap(xmlStr string) Params {

	params := make(Params)
	decoder := xml.NewDecoder(strings.NewReader(xmlStr))

	var (
		key   string
		value string
	)

	for t, err := decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement: // 开始标签
			key = token.Name.Local
		case xml.CharData: // 标签内容
			content := string([]byte(token))
			value = content
		}
		if key != "xml" {
			if value != "\n" {
				params.SetString(key, value)
			}
		}
	}

	return params
}

func MapToXml(params Params) string {
	var buf bytes.Buffer
	buf.WriteString(`<xml>`)
	for k, v := range params {
		buf.WriteString(`<`)
		buf.WriteString(k)
		buf.WriteString(`><![CDATA[`)
		buf.WriteString(v)
		buf.WriteString(`]]></`)
		buf.WriteString(k)
		buf.WriteString(`>`)
	}
	buf.WriteString(`</xml>`)

	return buf.String()
}

// 用时间戳生成随机字符串
func nonceStr() string {
	return strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
}

// 将Pkcs12转成Pem
func pkcs12ToPem(p12 []byte, password string) tls.Certificate {

	blocks, err := pkcs12.ToPEM(p12, password)

	// 从恐慌恢复
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()

	if err != nil {
		panic(err)
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}
	return cert
}

// 解析Pem
func parsePem(certData, keyData []byte, password string) tls.Certificate {
	// 从恐慌恢复
	defer func() {
		if x := recover(); x != nil {
			log.Print(x)
		}
	}()
	keyPEMBlock, rest := pem.Decode(keyData)
	if len(rest) > 0 {
		panic("Decode key failed!")
	}

	if x509.IsEncryptedPEMBlock(keyPEMBlock) {
		keyDePEMByte, err := x509.DecryptPEMBlock(keyPEMBlock, []byte(password))
		if err != nil {
			panic(err)
		}

		// 解析出其中的RSA 私钥
		key, err := x509.ParsePKCS1PrivateKey(keyDePEMByte)
		if err != nil {
			panic(err)
		}
		// 编码成新的PEM 结构
		keyData = pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(key),
			},
		)

	}
	cert, err := tls.X509KeyPair(certData, keyData)
	if err != nil {
		panic(err)
	}
	return cert

}
