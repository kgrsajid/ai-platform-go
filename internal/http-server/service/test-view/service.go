package testviewservice

type TestViewRepo interface {
	AddTestView(testId uint, userId uint) error
}

type Service struct {
	TestViewRepo TestViewRepo
}

func New(TestViewRepo TestViewRepo) *Service {
	return &Service{
		TestViewRepo: TestViewRepo,
	}
}

func (s *Service) AddTestView(testId uint, userId uint) error {
	err := s.TestViewRepo.AddTestView(testId, userId)
	return err
}
