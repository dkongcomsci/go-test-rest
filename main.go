package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	ora "gopkg.in/rana/ora.v4"
)

type Topic struct {
	TopicID int
	Name    string
}
type Page struct {
	Title string
	Body  []byte
}
type Person struct {
	ID        string   `json:"id,omitempty`
	Firstname string   `json:"firstname,omitempty`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}
type Subscriber struct {
	PricePlanName      string `json:"pricePlanName,omitempty"`
	Sock               string `json:"sock,omitempty"`
	PackageDescription string `json:"packageDescription,omitempty"`
	OutBundleRate      string `json:"outBundleRate,omitempty"`
	OtherOutBundleRate string `json:"otherOutBundleRate,omitempty"`
}

var people []Person

func GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}
func searchPriceplan(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	msisdn := params["msisdn"]
	fmt.Println(params["msisdn"])
	var subNo, price_plan, price_plan_code, soc, package_desc interface{}

	env, srv, ses, err := ora.NewEnvSrvSes("THOTSAPORN_SAK_B/thot#0718@172.16.12.119:1548/CUST01")
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()
	defer srv.Close()
	defer ses.Close()

	sql := "select agreement_no sub_no from AGREEMENT_RESOURCE  where RESOURCE_VALUE ='" + msisdn + "'  and expiration_date is null "
	fmt.Println("SQL :" + sql)
	stmtCount, err := ses.Prep(sql)
	defer stmtCount.Close()
	if err != nil {
		panic(err)
	}
	rset, err := stmtCount.Qry()
	if err != nil {
		panic(err)
	}
	row := rset.NextRow()
	if row != nil {
		subNo = row[0]
	}
	if err := rset.Err(); err != nil {
		panic(err)
	}
	fmt.Println(subNo)

	env, srv, ses, err = ora.NewEnvSrvSes("pipeline/P@ssw0rd@172.19.216.129:1556/LGSTST")
	if err != nil {
		log.Fatal(err)
	}
	defer env.Close()
	defer srv.Close()
	defer ses.Close()

	sql = fmt.Sprintf("select PRICE_PLAN,PRICE_PLAN_CODE,SOC_LIST,PACKAGE_DESC  from v_package_desc where  agreement_no = %v", subNo)
	fmt.Println("SQL :" + sql)
	stmtCount, err = ses.Prep(sql)
	defer stmtCount.Close()
	if err != nil {
		panic(err)
	}
	rset, err = stmtCount.Qry()
	if err != nil {
		panic(err)
	}
	row = rset.NextRow()
	if row != nil {
		fmt.Println("Row :")
		fmt.Println(row[0])
		price_plan = row[0]
		price_plan_code = row[1]
		soc = row[2]
		package_desc = row[3]
	}
	if err := rset.Err(); err != nil {
		panic(err)
	}
	fmt.Println(price_plan)
	fmt.Println(price_plan_code)
	fmt.Println(soc)
	fmt.Println(package_desc)

	var voice_rate, other_rate, buffet_desc string
	if _, err = ses.PrepAndExe("CALL GET_RATE_DESC_TXT_2(11126012,:1,:2,:3)", &voice_rate, &other_rate, &buffet_desc); err != nil {
		log.Fatal(err)
	}
	fmt.Println("voice_rate :" + voice_rate)
	fmt.Println("other_rate :" + other_rate)
	fmt.Println("buffet_desc :" + buffet_desc)

	var subscriber Subscriber
	subscriber = Subscriber{
		PricePlanName:      fmt.Sprintf("%v", price_plan),
		Sock:               fmt.Sprintf("%v", soc),
		PackageDescription: fmt.Sprintf("%v", package_desc) + buffet_desc,
		OutBundleRate:      voice_rate,
		OtherOutBundleRate: other_rate}
	json.NewEncoder(w).Encode(subscriber)
}
func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}
func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}
func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
func handler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "demo.html")
}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("demo.html")

	var results []Topic
	results = []Topic{
		Topic{1, "Title1"},
		Topic{2, "Title2"},
	}

	t.Execute(w, map[string]interface{}{
		"project_name": "MY DATA",
		"results":      results,
	})
}

func main() {

	// var message string
	// message = "Hello, Go"
	// message2 := "Hello, Gog"
	// fmt.Println("Hello, Go")
	// fmt.Println(message)
	// fmt.Println(message2)
	// var names []string
	// // เพิ่ม element เข้าไป
	// names = append(names, "Somchai")
	// names = append(names, "Somsree")
	// names = append(names, "Somset")
	// // names2 := []string{"Somchai", "Somsree", "Somset"}
	// i := 10
	// if i%2 == 0 {
	// 	// even
	// } else {
	// 	// odd
	// }
	// for i := 1; i <= 10; i++ {
	// 	fmt.Println(i)
	// }
	// switch i {
	// case 0:
	// 	fmt.Println("Zero")
	// case 1:
	// 	fmt.Println("One")
	// case 2:
	// 	fmt.Println("Two")
	// case 3:
	// 	fmt.Println("Three")
	// case 4:
	// 	fmt.Println("Four")
	// case 5:
	// 	fmt.Println("Five")
	// default:
	// 	fmt.Println("Unknown Number")
	// }

	// var user string
	// var user1 string
	// var user2 string
	// if _, err = ses.PrepAndExe("CALL GET_RATE_DESC_TXT_2(11126012,:1,:2,:3)", &user, &user1, &user2); err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("user :" + user)
	// fmt.Println("user1 :" + user1)
	// fmt.Println("user2 :" + user2)

	// db, err := sql.Open("ora", "pipeline/P@ssw0rd@172.19.216.129:1556/LGSTST")

	// if err != nil {
	// 	panic(err)
	// }

	// stmt, err := db.Prepare("CALL GET_RATE_DESC_TXT_2(%v,:1,:2,:3)", "15359321")
	// stmt, err := db.Prepare("CALL GET_RATE_DESC_TXT_2(:i_user_guid, :rate_desc_txt,:other_rate_txt,:buffet_desc)")

	// if err != nil {
	// 	panic(err)
	// }

	// guid := "15359321"

	// var str, str1, str2 string
	// res, err := stmt.Exec(guid,
	// 	sql.Named("rate_desc_txt", sql.Out{Dest: &str}),
	// 	sql.Named("other_rate_txt", sql.Out{Dest: &str1}),
	// 	sql.Named("buffet_desc", sql.Out{Dest: &str2}))

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println(res)

	// env, srv, ses, err := ora.NewEnvSrvSes("pipeline/P@ssw0rd@172.19.216.129:1556/LGSTST")

	// // Set timeout (Go 1.8)
	// ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	// // Set prefetch count (Go 1.8)
	// ctx = ora.WithStmtCfg(ctx, ora.Cfg().StmtCfg.SetPrefetchRowCount(50000))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer env.Close()
	// defer srv.Close()
	// defer ses.Close()

	// sql := fmt.Sprintf("SELECT * FROM LEGO_PPSOC where SOC_CD = %v", "13145825")
	// fmt.Println("SQL :")
	// fmt.Println(sql)
	// stmtCount, err := ses.Prep(sql)
	// defer stmtCount.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// rset, err := stmtCount.Qry()
	// if err != nil {
	// 	panic(err)
	// }
	// row := rset.NextRow()
	// if row != nil {
	// 	fmt.Println("Row :")
	// 	fmt.Println(row[0], row[1], row[2], row[3])
	// }
	// if err := rset.Err(); err != nil {
	// 	panic(err)
	// }

	// call stored procedure
	// pass *Rset to Exe to receive the results of a sys_refcursor
	// sql := fmt.Sprintf("CALL GET_RATE_DESC_TXT_2(%v,:1,:2,:3)", "15359321")
	// stmtProcCall, err := ses.Prep(sql)
	// defer stmtProcCall.Close()
	// if err != nil {
	// 	panic(err)
	// }
	// var str, str2, str3 string
	// rowsAffected, err := stmtProcCall.Exe(&str, &str2, &str3)
	// fmt.Println("str : " + str)
	// fmt.Println("str2 : " + str2)
	// fmt.Println("str3 : " + str3)
	// fmt.Println(rowsAffected)
	// if err != nil {
	// 	panic(err)
	// }
	// procRset := &ora.Rset{}
	// rowsAffected, err = stmtProcCall.Exe(procRset)
	// if procRset.IsOpen() {
	// 	for procRset.Next() {
	// 		fmt.Println(procRset.Row[0], procRset.Row[1])
	// 	}
	// 	if err := procRset.Err(); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(procRset.Len())
	// }

	// db, err := sql.Open("ora", "THOTSAPORN_SAK_B/thot#0718@172.16.12.119:1548/CUST01")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/", handler)
	router.HandleFunc("/homeHandler", homeHandler).Methods("GET")
	router.HandleFunc("/search/{msisdn}", searchPriceplan).Methods("GET")
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")

	srver := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srver.ListenAndServe())

	// log.Fatal(http.ListenAndServe(":8000", router))
}
