package employees

type EmployeeService struct {
	repo *Repository
}

func NewService(repo *Repository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) GetAllEmployees() ([]Employee, error) {
	return s.repo.GetAllEmployees()
}


func (s *EmployeeService) GetEmployeeByID(id string) (*Employee, error) {
	return s.repo.GetEmployeeByID(id)
}