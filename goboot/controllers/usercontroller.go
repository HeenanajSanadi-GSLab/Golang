package controllers

import (

    "os"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"../config"
	"github.com/julienschmidt/httprouter"
	
	"../daos"
	"../models"
	"../errorhandling"
	"github.com/google/logger"
)

// MyUserController controller
type MyUserController struct {
	myuserDao daos.UserDao
}

type Response errorhandling.Response

// NewMyUserController creating controller
func NewMyUserController() *MyUserController {
	myconfig, err := config.GetConfiguration()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	myController := &MyUserController{}
	myController.myuserDao = daos.UserFactoryDao(myconfig.Engine)
	return myController
}

/*
GetUsers curl -GET http://localhost:8002/users
curl -GET http://localhost:8002/users
*/
func (uc MyUserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	myconfig, err := config.GetConfiguration()
	lf, err := os.OpenFile(myconfig.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
    	logger.Fatalf("Failed to open log file: %v", err)
  	}
  	defer lf.Close()
  	logstmt := logger.Init("LoggerFile", false, true, lf)

	logstmt.Info("List all Users  >>")
	log.Printf("List all Users  >> ")
	us, err := uc.myuserDao.GetAll()
	if err != nil {
		log.Fatal(err)
		logstmt.Error("No records in DB")
		return
	}
	//jsonUs, _ := json.Marshal(us)
	response := Response{StatusCode:200, ErrorMessage:"SUCCESS", Payload:us}
	res, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", res)
}

/*
CreateUser curl -XPOST -H 'Content-Type: application/json' -d '{"name": "L John Mammen", "gender": "male", "age": 15}' http://localhost:8002/user
curl -XPOST -H 'Content-Type: application/json' -d '{"name": "L John Mammen", "gender": "male", "age": 15}' http://localhost:8002/user
*/
func (uc MyUserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	myconfig, err := config.GetConfiguration()
	lf, err := os.OpenFile(myconfig.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
    	logger.Fatalf("Failed to open log file: %v", err)
  	}
  	defer lf.Close()
  	logstmt := logger.Init("LoggerFile", false, true, lf)

	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)

	if flag := uc.myuserDao.CheckDuplicateUser(u.Name); flag == true {
		response := Response{StatusCode:409, ErrorMessage:"Duplicate User Name"}
		jsonU, _ := json.Marshal(response)
        w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(409)
		fmt.Fprintf(w, "%s", jsonU)
        return
    }

	err = uc.myuserDao.Create(&u)
	if err != nil {
		log.Fatal(err)
	}
	logstmt.Info("Create User ID of user is >> ", u.Id)
	log.Printf("Create User ID of user is >> %d", u.Id)

	users := []models.User{}
	users = append(users,u)

	response := Response{StatusCode:200, ErrorMessage:"SUCCESS", Payload:users}
	jsonU, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", jsonU)
}


// UpdateUser curl -H 'Content-Type: application/json' -H 'Accept: application/json' -X PUT -d '{"name": "L John Mammen", "gender": "male", "age": 15, "id":5}' http://localhost:8002/user
/*
curl -H 'Content-Type: application/json' -H 'Accept: application/json' -X PUT -d '{"name": "L John Mammen", "gender": "male", "age": 15, "id":5}' http://localhost:8002/user
*/

func (uc MyUserController) UpdateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	myconfig, err := config.GetConfiguration()
	lf, err := os.OpenFile(myconfig.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
    	logger.Fatalf("Failed to open log file: %v", err)
  	}
  	defer lf.Close()
  	logstmt := logger.Init("LoggerFile", false, true, lf)

	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)
	err = uc.myuserDao.Update(&u)
	if err != nil {
		log.Fatal(err)
		return
	}
	//jsonU, _ := json.Marshal(u)
	log.Printf("Update User is >>>>> %d", u.Id)
	logstmt.Info("Update User is >>>>> ", u.Id)

	user, err := uc.myuserDao.Get(u.Id)
	if user.Id == 0{
		response := Response{StatusCode:404, ErrorMessage:"Doesn't exist in DB"}
		jsonU, _ := json.Marshal(response)
        w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		fmt.Fprintf(w, "%s", jsonU)
        return
	}	

	users := []models.User{}
	users = append(users,user)

	response := Response{StatusCode:200, ErrorMessage:"SUCCESS", Payload:users}
	jsonU, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", jsonU)
}

/*
RemoveUser curl -XDELETE http://localhost:8002/user/id
curl -XDELETE http://localhost:8002/user/id
*/
func (uc MyUserController) RemoveUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	myconfig, err := config.GetConfiguration()
	lf, err := os.OpenFile(myconfig.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
    	logger.Fatalf("Failed to open log file: %v", err)
  	}
  	defer lf.Close()
  	logstmt := logger.Init("LoggerFile", false, true, lf)

	id, err := strconv.Atoi(p.ByName("id"))
	log.Printf("RemoveUser ID of user is >>>>> %d", id)
	logstmt.Info("RemoveUser ID of user is >>>>> ", id)

	user, err := uc.myuserDao.Get(id)

	if user.Id == 0{
		response := Response{StatusCode:404, ErrorMessage:"FAILED"}
		jsonU, _ := json.Marshal(response)
        w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		fmt.Fprintf(w, "%s", jsonU)
        return
	}	

	err = uc.myuserDao.Delete(id)
	if err != nil {
		log.Fatal(err)
		return
	}
	/*dept_err := uc.myuserDao.DeleteDepartmentByUserId(id)
		if dept_err != nil {
		log.Fatal(dept_err)
	}*/
	users := []models.User{}
	users = append(users,user)

	response := Response{StatusCode:200, ErrorMessage:"SUCCESS", Payload:users}
	jsonU, _ := json.Marshal(response)
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", jsonU)

}

/*
GetUser curl -GET http://localhost:8002/user/id
curl -GET http://localhost:8002/user/id
*/


func (uc MyUserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	myconfig, err := config.GetConfiguration()
	lf, err := os.OpenFile(myconfig.LogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
    	logger.Fatalf("Failed to open log file: %v", err)
  	}
  	defer lf.Close()
  	logstmt := logger.Init("LoggerFile", false, true, lf)

	id, err := strconv.Atoi(p.ByName("id"))
	log.Printf("GET user ID is >>>>> %d", id)
	logstmt.Info("GET user ID is >>>>> ", id)
	user, err := uc.myuserDao.Get(id)
	
	if user.Id == 0{
		er := Response{StatusCode:404, ErrorMessage:"User does not exist in DB"}
		jsonU, _ := json.Marshal(er)
        w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		fmt.Fprintf(w, "%s", jsonU)
        return
	}	

	//jsonU, _ := json.Marshal(user)
	users := []models.User{}
	users = append(users,user)
	response := Response{StatusCode:200, ErrorMessage:"SUCCESS", Payload:users}
	res, _ := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", res)
}

