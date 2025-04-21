package main

import (
	"encoding/json"
	"os"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/icrowley/fake"
)

type Person struct {
	FIO string `json:"fio"`
	Fprice string	`json:"fprice"`
	Sprice string `json:"sprice"`
	Fpaid bool	`json:"fpaid"`
	Spaid bool `json:"spaid"`
	Birthday string `json:"birthday"`
	YearsOld string `json:"yearsold"`
	Email string `json:"email"`
	StudentNumber string `json:"studentnumber"`
}

func main() {
	LogPass := map[string]string{
		"palagnyuk.a.a@edu.mirea.ru" : "Aa19102006.",
		"bakyr.m.y@edu.mirea.ru" : "Mert251326Mert@",
		"andreev.a.r1@edu.mirea.ru" : "123EWQasdD!",
		"gorchakov.a.a@edu.mirea.ru" : "gognop-wyzpUp-1zyqru",
		"prygov.k.d@edu.mirea.ru" : "Kirill_200622",
	}

	var persons []Person

	for k, v := range LogPass{
		client := resty.New()

		client.SetHeader("User-Agent", fake.UserAgent())

		_, err := client.R().Get("https://lk.mirea.ru/auth.php")

		if err != nil{
			log.Fatalf("Ошибка GET-запроса", err)
		}

		data := map[string]string{
			"AUTH_FORM" : "Y",
			"TYPE" : "AUTH",
			"USER_LOGIN" : k,
			"USER_PASSWORD" : v,
			"USER_REMEMBER": "Y",
		}

		resp, err := client.R().SetFormData(data).Post("https://lk.mirea.ru/auth.php?login=yes")
		
		if err != nil{
			log.Fatalf("Ошибка POST-запроса")
		}

		fprice, sprice := takePrice(resp)
		fpaid, spaid := takePaid(resp)
		yearsold, birthday := takeALLDataOfBirthday(resp)

		person := Person{FIO: takeFIO(resp), Fprice: fprice, Sprice: sprice, Fpaid: fpaid, Spaid: spaid, Birthday: birthday, YearsOld: yearsold, Email: takeEmail(resp), StudentNumber: takeStudentNumber(resp)} 

		persons = append(persons, person)
	}

	for _, v := range persons{
		fmt.Println(v)
	}

	toJSON(persons)
}

func takeFIO(resp *resty.Response)	string  {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil{
		log.Fatalf("Ошибка создания файла", err)
	}

	return (doc.Find(".ml-6").Find("h1").Text())
}

func takePrice(resp *resty.Response) (string, string) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil{
		log.Fatalf("Ошибка создания файла", err)
	}

	return ExtraPrice(doc.Find("h4.text-lg.font-bold").Text())
}

func ExtraPrice(allprice string) (string, string){
	allprice = strings.ReplaceAll(allprice,"	", "")
	allprice = strings.ReplaceAll(allprice," ", "")
	
	return allprice[:12], allprice[12:]
}

func takePaid(resp *resty.Response) (bool, bool){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil{
		log.Fatalf("Ошибка создания файла", err)
	}

	return ExtraPaid(doc.Find("p.font-normal.text-gray-500.text-sm").Text())
}

func ExtraPaid(allPaid string) (bool, bool){
	if allPaid[:16] == "Оплачено" && allPaid[16:] == "Оплачено" {
		return true, true
	}

	return false, false
}

func takeALLDataOfBirthday(resp *resty.Response) (string, string)  {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil{
		log.Fatalf("Ошибка создания файла", err)
	}

	return ExtraBirthdayAndYearsold(doc.Find(".font-semibold.text-gray-900.text-sm").Eq(4).Text())
}

func ExtraBirthdayAndYearsold(birthdayAndyearsold string) (string, string){
	birthdayAndyearsold = strings.ReplaceAll(birthdayAndyearsold, " 	", "")
	birthdayAndyearsold = strings.ReplaceAll(birthdayAndyearsold, " ", "")
	birthdayAndyearsold = strings.ReplaceAll(birthdayAndyearsold, "\n", "")

	return birthdayAndyearsold[:8], birthdayAndyearsold[8:]
}

func takeEmail(resp *resty.Response) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil{
		log.Fatalf("Ошибка создания файла", err)
	}

	return doc.Find(".font-semibold.text-gray-900.text-sm").Eq(3).Text()
}

func takeStudentNumber(resp *resty.Response) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil{
		log.Fatalf("Ошибка создания файла", err)
	}

	return doc.Find(".font-semibold.text-gray-900.text-sm").Eq(2).Text()
}

func toJSON(persons []Person){
	file, err := os.Create("information.json")

	if err != nil{
		log.Fatalf("Ошибка создания JSON файла", err)
	}

	defer file.Close()

	encoder := json.NewEncoder(file)

	encoder.SetIndent("", "  ")
	err = encoder.Encode(persons)

	if err != nil {
	log.Fatalf("❌Ошибка записи JSON: %v", err)
	}

	fmt.Printf("✅Готово! Найдено людей: %d. Сохранено в information.json\n", len(persons))
}

