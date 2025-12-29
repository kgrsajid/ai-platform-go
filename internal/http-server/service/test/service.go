package testservice

import (
	req "project-go/internal/http-server/dto/request"
	res "project-go/internal/http-server/dto/response"
	"project-go/internal/models"
)

type TestRepo interface {
	CreateTest(test *models.Test) (*models.Test, error)
	GetAllTest(testFilter req.TestFilter) ([]res.TestWithQuestionsCount, int64, error)
	GetTestById(testId uint64) (*models.Test, error)
	AddTestResult(testReq req.TestResultReq) (*models.TestResult, error)
	GetAllUserTestResults(filter req.GetALlTestResultsFilter) ([]models.TestResult, error)
}

type QuestionAdd interface {
	CreateQuestion(question *models.TestQuestion) (*models.TestQuestion, error)
}

type Service struct {
	TestRepo     TestRepo
	QuestionRepo QuestionAdd
}

func New(TestRepo TestRepo, QuestionRepo QuestionAdd) *Service {
	return &Service{
		TestRepo:     TestRepo,
		QuestionRepo: QuestionRepo,
	}
}

func (s *Service) TestCreate(test req.TestRequest) (*models.Test, error) {
	testModel := mapTestRequestToModel(test)
	createdTest, err := s.TestRepo.CreateTest(testModel)
	if err != nil {
		return nil, err
	}

	return createdTest, nil
}

func (s *Service) AddTestResult(test req.TestResultReq) (*res.TestResultResponse, error) {
	testResult, err := s.TestRepo.AddTestResult(test)
	if err != nil {
		return nil, err
	}
	return &res.TestResultResponse{
		ID:         testResult.ID,
		Score:      testResult.Score,
		MaxScore:   testResult.MaxScore,
		Percentage: testResult.Percentage,
		Attempt:    testResult.Attempt,
	}, err
}

func (s *Service) GetAllUserTestResults(test req.GetALlTestResultsFilter) ([]res.TestResultResponse, error) {
	testResults, err := s.TestRepo.GetAllUserTestResults(test)
	if err != nil {
		return nil, err
	}
	result := make([]res.TestResultResponse, 0, len(testResults))
	for _, val := range testResults {
		result = append(result, res.TestResultResponse{
			ID:         val.ID,
			Score:      val.Score,
			MaxScore:   val.MaxScore,
			Percentage: val.Percentage,
			Attempt:    val.Attempt,
			CreatedAt:  val.CreatedAt,
		})
	}
	return result, nil
}

func (s *Service) GetAllTest(testFilter req.TestFilter) ([]res.TestWithQuestionsCount, int64, error) {
	tests, total, err := s.TestRepo.GetAllTest(testFilter)
	if err != nil {
		return nil, 0, err
	}
	return tests, total, nil
}

func (s *Service) GetTestById(testId uint64) (*res.TestDetailsResponse, error) {
	tests, err := s.TestRepo.GetTestById(testId)
	if err != nil {
		return nil, err
	}
	return res.ToTestDetailsResponse(tests), nil
}

func mapTestRequestToModel(req req.TestRequest) *models.Test {
	test := &models.Test{
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		AuthorID:    req.AuthorId,
		IsPrivate:   req.IsPrivate,
		Difficulty:  models.Difficulty(req.Difficulty),
	}
	for _, catID := range req.Categories {
		test.Categories = append(test.Categories, models.Category{
			ID: catID,
		})
	}

	for _, q := range req.Questions {
		question := models.TestQuestion{
			Question: q.Question,
		}
		for _, o := range q.Options {
			option := models.TestOption{
				OptionText: o.OptionText,
				IsCorrect:  o.IsCorrect,
			}
			question.Options = append(question.Options, option)
		}
		test.Questions = append(test.Questions, question)
	}

	return test
}
