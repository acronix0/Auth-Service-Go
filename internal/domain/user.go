package domain

type User struct {
	ID       int    
	Email    string 
	Password string 
	Name     string 
	Phone    string 
	Blocked  bool   
	Role     string 
}
