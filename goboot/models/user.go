package models

type (
	User struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Gender string `json:"gender"`
		Age    int    `json:"age"`
		Department *Department `json:"department"`
	}

	Department struct{
	Department_id int `json:"department_id"`
	Department_name string `json:"department_name"`
}
)
