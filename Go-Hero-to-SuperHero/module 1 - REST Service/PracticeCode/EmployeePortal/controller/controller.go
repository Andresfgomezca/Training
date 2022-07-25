package controller

import (
	"employeeportal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

var (
	portal             = "/portal/api/v1" 
	postEmployee       = "/employee"
	getEmployees       = "/employee"
	getEmployee        = "/employee/{employeeName}"
	patchEmployee      = "/employee/{employeeName}"
	deleteEmployee     = "/employee/{employeeName}"
	jsonMarshalFunc    = json.Marshal
	allowedSortOptions = []string{"asc", "desc"}
)

var mockDB = []models.Employee{
	{
		ID:       3,
		Name:     "test1",
		DOJ:      time.Now(),
		Skillset: []string{"sql", "casssandra", "aws"},
	},
	{
		ID:       24,
		Name:     "test2",
		DOJ:      time.Now(),
		Skillset: []string{"mongo", "elasticsearch", "azure"},
	},
}

func Handlers() http.Handler {

	r := mux.NewRouter()
	//POST
	r.HandleFunc(portal+postEmployee, CreateEmployee).Methods(http.MethodPost)
	//GET
	r.HandleFunc(portal+getEmployee, GetEmployee).Methods(http.MethodGet)
	r.HandleFunc(portal+getEmployees, GetEmployee).Methods(http.MethodGet)
	//PUT
	r.HandleFunc(portal+patchEmployee, updateEmployee).Methods(http.MethodPut)
	//DELETE
	r.HandleFunc(portal+deleteEmployee, removeEmployee).Methods(http.MethodDelete)

	return r
}

//POST
func CreateEmployee(rw http.ResponseWriter, r *http.Request) {
	emp := models.Employee{}
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("please check the request body :: %v", err)
		rw.Write([]byte(msg))
		return
	}

	// add persistence here
	// for now we will just save it in out memory and
	// return the decoded employee object
	mockDB = append(mockDB, emp)

	fmt.Printf("The Employee is :: \n ID :: %v\n Name:: %v\n DOJ :: %v\n Skillset :: %v\n", emp.ID, emp.Name, emp.DOJ, emp.Skillset)
}

//gets

type GetEmployeesResponse struct {
	Employees []models.Employee `json:"employees"`
}
//employee by name
func GetEmployee(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eName := params["employeeName"]
	//if the get request has no name it will return the entire list
	if eName == "" {
		res, err := jsonMarshalFunc(mockDB)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("unable to serialize response :: %v", err)
			rw.Write([]byte(msg))
			return
		}
		rw.WriteHeader(http.StatusOK)
		rw.Write(res)
		return
	}

	// since we don't have persistence yet
	// we will use mockDB an in memory db to check if the employee exists or not
	for _, e := range mockDB {
		if strings.Compare(e.Name, eName) == 0 {
			rw.Header().Set("Content-Type", "application/json")
			rw.WriteHeader(http.StatusOK)
			rEmp := fmt.Sprintf("The Employee is Present with Details :: \n ID :: %v\n Name:: %v\n DOJ :: %v\n Skillset :: %v\n", e.ID, e.Name, e.DOJ, e.Skillset)
			rw.Write([]byte(rEmp))
			return
		}
	}

	msg := fmt.Sprintf("no employee found with the given name :: %v", eName)
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(msg))
}

// PUT
func updateEmployee(rw http.ResponseWriter, r *http.Request) {
	//employee updated by name
	params := mux.Vars(r)
	eName := params["employeeName"]
	if eName == "" {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("employeeName to update cannot be empty"))
		return
	}

	patchEmployee := models.Employee{}
	//reads the body of the request and patch the information to the var patch employee
	err := json.NewDecoder(r.Body).Decode(&patchEmployee)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		//the body is not correct
		msg := fmt.Sprintf("please check the request body :: %v", err)
		rw.Write([]byte(msg))
		return
	}
	//update of the employee's information
	for i, e := range mockDB {
		if eName == e.Name {
			//in this way the lement employee is correctly updated
			mockDB[i] = patchEmployee
			//mockDB = append(mockDB, patchEmployee)
			//update by element is not persisting the information in the slice
			e.ID = patchEmployee.ID
			e.DOJ = patchEmployee.DOJ
			e.Name = patchEmployee.Name
			e.Skillset = patchEmployee.Skillset

			// this is where we should be persisting the updated data
			// we will just set these values in our mockDB array
			rw.WriteHeader(http.StatusNoContent)
			msg := fmt.Sprintf("Resource with name %v updated successfully.", eName)
			rw.Write([]byte(msg))
			return
		}
	}

	rw.WriteHeader(http.StatusNotFound)
	msg := fmt.Sprintf("Resource with name %v not found.", eName)
	rw.Write([]byte(msg))
}

//DELETE
//This function removes the employee with this name
func remove(empDB []models.Employee, i int) []models.Employee {
	return append(empDB[:i], empDB[i+1:]...)
}

//Delete by name
func removeEmployee(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	eName := params["employeeName"]
	if eName == "" {
		rw.WriteHeader(http.StatusOK)
		rw.Write([]byte("employeeName to delete cannot be empty"))
		return
	}

	for i, e := range mockDB {
		if eName == e.Name {

			// this is where we should be persisting the data (delete request)
			// we will just remove the employee from mockDB array
			mockDB = remove(mockDB, i)
			fmt.Println(mockDB)
			rw.WriteHeader(http.StatusOK)
			msg := fmt.Sprintf("Resource with name %v deleted successfully.", eName)
			rw.Write([]byte(msg))
			return
		}
	}

	rw.WriteHeader(http.StatusNotFound)
	msg := fmt.Sprintf("Resource with name %v not found.", eName)
	rw.Write([]byte(msg))
}

/* code to sort the response according to the arguments of the request, default is asc
sortOrder := r.URL.Query().Get("sortOrder")
if !stringInSlice(sortOrder, allowedSortOptions) {
	rw.WriteHeader(http.StatusBadRequest)
	msg := fmt.Sprintf("wrong option for query parameter sortOrder valid options are :: %v", allowedSortOptions)
	rw.Write([]byte(msg))
	return
}
if sortOrder == "" {
	sortOrder = "asc"
}
if sortOrder == "asc" {
	sort.SliceStable(mockDB[:], func(i, j int) bool {
		return mockDB[i].Name < mockDB[j].Name
	})
} else {
	sort.SliceStable(mockDB[:], func(i, j int) bool {
		return mockDB[i].Name > mockDB[j].Name
	})
}
*/
