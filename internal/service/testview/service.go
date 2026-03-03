package testviewservice

type TestViewRepository interface {
	AddTestView(testId uint, userId uint) error
}

type Service struct {
	testViewRepo TestViewRepository
}

func New(testViewRepo TestViewRepository) *Service {
	return &Service{testViewRepo: testViewRepo}
}

func (s *Service) AddTestView(testId uint, userId uint) error {
	return s.testViewRepo.AddTestView(testId, userId)
}
