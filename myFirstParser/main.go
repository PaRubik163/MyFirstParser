package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/icrowley/fake"
)

type Person struct {
	FIO           string `json:"fio"`
	Fprice        string `json:"f_price"`
	Sprice        string `json:"s_price"`
	Fpaid         bool   `json:"f_paid"`
	Spaid         bool   `json:"s_paid"`
	Birthday      string `json:"birthday"`
	YearsOld      string `json:"yearsold"`
	Email         string `json:"email"`
	StudentNumber string `json:"student_number"`
	Login         string `json:"login"`
	WhichYear     string `json:"year_of_university"`
	GroupNumber   string `json:"group_number"`
}

func (p *Person) takeFIO(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	p.FIO = (doc.Find(".ml-6").Find("h1").Text())
}

func (p *Person) takePrice(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	p.Fprice, p.Sprice = extraPrice(doc.Find("h4.text-lg.font-bold").Text())
}

func (p *Person) takePaid(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	p.Fpaid, p.Spaid = extraPaid(doc.Find("p.font-normal.text-gray-500.text-sm").Text())
}

func (p *Person) takeALLDataOfBirthday(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	p.YearsOld, p.Birthday = extraBirthdayAndYearsold(doc.Find(".font-semibold.text-gray-900.text-sm").Eq(4).Text())
}

func (p *Person) takeEmail(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	p.Email = doc.Find(".font-semibold.text-gray-900.text-sm").Eq(3).Text()
}

func (p *Person) takeStudentNumber(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	p.StudentNumber = doc.Find(".font-semibold.text-gray-900.text-sm").Eq(2).Text()
}

func (p *Person) takeLogin(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	p.Login = doc.Find(".font-semibold.text-gray-900.text-sm").Eq(6).Text()
}

func (p *Person) takeWhichYearOfUniversity(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	p.WhichYear = doc.Find("p.text-sm.text-grey-900").Eq(0).Text()
}

func (p *Person) takeGroupNumber(resp *resty.Response){
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp.String()))

	if err != nil {
		log.Fatalf("Ошибка создания файла", err)
	}

	tmp := doc.Find("p.text-white.bg-persian-blue-800.rounded-2xl").Text()

	tmp = strings.ReplaceAll(tmp, "		", "")
	tmp = strings.ReplaceAll(tmp, "\n", "")
	tmp = strings.ReplaceAll(tmp, " ", "")

	p.GroupNumber = tmp
}

func (p *Person) print(persons *[]Person){
	for _, v := range *persons{
		fmt.Println(v)
	}
}

func extraPrice(allprice string) (string, string){
	allprice = strings.ReplaceAll(allprice, "	", "")
	allprice = strings.ReplaceAll(allprice, " ", "")

	return allprice[:12], allprice[12:]
}

func extraPaid(allPaid string) (bool, bool){
	if allPaid[:16] == "Оплачено" && allPaid[16:] == "Оплачено" {
		return true, true
	}

	return false, false
}

func extraBirthdayAndYearsold(birthdayAndyearsold string) (string, string){
	birthdayAndyearsold = strings.ReplaceAll(birthdayAndyearsold, " 	", "")
	birthdayAndyearsold = strings.ReplaceAll(birthdayAndyearsold, " ", "")
	birthdayAndyearsold = strings.ReplaceAll(birthdayAndyearsold, "\n", "")
	birthdayAndyearsold = strings.ReplaceAll(birthdayAndyearsold, "л", " л")

	return birthdayAndyearsold[:9], birthdayAndyearsold[9:]
}
func Loginned(LogPass map[string]string, person *Person, persons *[]Person) error {

	for k, v := range LogPass {
		client := resty.New()
		client.SetHeader("User-Agent", fake.UserAgent())

		if _, err := client.R().Get("https://lk.mirea.ru/auth.php"); err != nil{
			return errors.New("Ошибка GET-запроса")	
		}

		data := map[string]string{
			"AUTH_FORM":     "Y",
			"TYPE":          "AUTH",
			"USER_LOGIN":    k,
			"USER_PASSWORD": v,
			"USER_REMEMBER": "Y",
		}

		resp, err := client.R().SetFormData(data).Post("https://lk.mirea.ru/auth.php?login=yes")

		if err != nil {
			return errors.New("Ошибка POST-запроса")
		}

		person.takeFIO(resp)
		person.takePrice(resp)
		person.takePaid(resp)
		person.takeALLDataOfBirthday(resp)
		person.takeEmail(resp)
		person.takeStudentNumber(resp)
		person.takeLogin(resp)
		person.takeWhichYearOfUniversity(resp)
		person.takeGroupNumber(resp)

		*persons = append(*persons, *person)
	}
	return nil
}

func toJSON(persons []Person) {
	file, err := os.Create("information.json")

	if err != nil {
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

func main() {

	LogPass := map[string]string{
		"palagnyuk.a.a@edu.mirea.ru": "Aa19102006.",
		"bakyr.m.y@edu.mirea.ru":     "Mert251326Mert@",
		"andreev.a.r1@edu.mirea.ru":  "123EWQasdD!",
		"gorchakov.a.a@edu.mirea.ru": "gognop-wyzpUp-1zyqru",
		"prygov.k.d@edu.mirea.ru":    "Kirill_200622",
	}

	person := &Person{}
	var persons []Person

	err := Loginned(LogPass, person, &persons)

	if err != nil{
		fmt.Println(err)
	}
	
	person.print(&persons)
	
	toJSON(persons)
}
