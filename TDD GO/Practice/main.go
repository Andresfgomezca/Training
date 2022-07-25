package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type task struct {
	ID      int    `json:ID`
	Name    string `json:Name`
	Content string `json:Content`
}

type allTasks []task

var tasks = allTasks{
	{
		ID:      1,
		Name:    "task one",
		Content: "Some Content",
	},
}

//handlefunc provides 2 parameters, a response to the client and a request of the client
func indexRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my API")
}

func getTasks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)

}
func getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "invalid id")
		return
	}

	for _, task := range tasks {
		if task.ID == taskID {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func createTasks(w http.ResponseWriter, r *http.Request) {
	var newTask task
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "insert a valid task")
	}

	json.Unmarshal(reqBody, &newTask)
	newTask.ID = len(tasks) + 1 //changing task ID
	tasks = append(tasks, newTask)
	//response to the client with the new task info and the http status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)

}

func deleteTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		fmt.Fprintf(w, "invalid id")
		return
	}

	for i, task := range tasks {
		if task.ID == taskID {
			//this function will replace the current slice with all the elements before and after
			//the element found with this ID
			tasks = append(tasks[:i], tasks[i+1:]...)
			fmt.Fprintf(w, "The task with the ID %v has been remove succesfully", taskID)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	var updatedTask task
	if err != nil {
		fmt.Fprintf(w, "invalid id")
		return
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "invalid data")
		return
	}
	for i, task := range tasks {
		if task.ID == taskID {
			//we remove the task that will be updated
			tasks = append(tasks[:i], tasks[i+1:]...)
			//then we add the task with the id of the removed task
			json.Unmarshal(reqBody, &updatedTask)
			updatedTask.ID = len(tasks) + 1 //changing task ID
			tasks = append(tasks, updatedTask)
			fmt.Fprintf(w, "The task with the ID %v has been updated succesfully", taskID)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
		}
	}
}

func main() {
	//enrutador
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexRoute)
	router.HandleFunc("/tasks", getTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	router.HandleFunc("/tasks", createTasks).Methods("POST")
	router.HandleFunc("/tasks/{id}", deleteTasks).Methods("DELETE")
	router.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")

	//to start the server http we need to use a method called listenandserve
	log.Fatal(http.ListenAndServe(":3000", router))
	//compile daemon allow us to keep the server up and reset the conection every time we add changes
	//to the api
}
