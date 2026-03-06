package testservice

import (
	"context"
	client "project-go/internal/client/chat"
	"project-go/internal/dto/request"
	res "project-go/internal/dto/response"
	"project-go/internal/models"
)

type TestRepository interface {
	CreateTest(test *models.Test) (*models.Test, error)
	UpdateTest(test *models.Test) (*models.Test, error)
	GetAllTest(filter request.TestFilter) ([]res.TestWithQuestionsCount, int64, error)
	GetTestById(testId uint64) (*models.Test, error)
	AddTestResult(testReq request.TestResultReq) (*models.TestResult, error)
	GetAllUserTestResults(filter request.GetAllTestResultsFilter) ([]models.TestResult, error)
}

type QuestionRepository interface {
	CreateQuestion(question *models.TestQuestion) (*models.TestQuestion, error)
}

type Service struct {
	testRepo     TestRepository
	questionRepo QuestionRepository
	aiAPI        client.AIClient
}

func New(testRepo TestRepository, questionRepo QuestionRepository, aiApi client.AIClient) *Service {
	return &Service{
		testRepo:     testRepo,
		questionRepo: questionRepo,
		aiAPI:        aiApi,
	}
}

func (s *Service) TestCreate(req request.TestRequest) (*models.Test, error) {
	testModel := mapTestRequestToModel(req)
	return s.testRepo.CreateTest(testModel)
}

func (s *Service) TestUpdate(req request.TestRequest, testId uint) (*models.Test, error) {
	testModel := mapTestRequestToModel(req)
	testModel.ID = testId
	return s.testRepo.UpdateTest(testModel)
}

func (s *Service) AddTestResult(req request.TestResultReq) (*res.TestResultResponse, error) {
	testResult, err := s.testRepo.AddTestResult(req)
	if err != nil {
		return nil, err
	}
	return &res.TestResultResponse{
		ID:         testResult.ID,
		Score:      testResult.Score,
		MaxScore:   testResult.MaxScore,
		Percentage: testResult.Percentage,
		Attempt:    testResult.Attempt,
	}, nil
}

func (s *Service) GetAllUserTestResults(filter request.GetAllTestResultsFilter) ([]res.TestResultResponse, error) {
	testResults, err := s.testRepo.GetAllUserTestResults(filter)
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

func (s *Service) GetAllTest(filter request.TestFilter) ([]res.TestWithQuestionsCount, int64, error) {
	return s.testRepo.GetAllTest(filter)
}

func (s *Service) GetTestById(testId uint64) (*res.TestDetailsResponse, error) {
	test, err := s.testRepo.GetTestById(testId)
	if err != nil {
		return nil, err
	}
	return res.ToTestDetailsResponse(test), nil
}

func mapTestRequestToModel(req request.TestRequest) *models.Test {
	test := &models.Test{
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		AuthorID:    req.AuthorId,
		IsPrivate:   req.IsPrivate,
		Difficulty:  models.Difficulty(req.Difficulty),
	}
	for _, catID := range req.Categories {
		test.Categories = append(test.Categories, models.Category{ID: catID})
	}
	for _, q := range req.Questions {
		question := models.TestQuestion{Question: q.Question}
		for _, o := range q.Options {
			question.Options = append(question.Options, models.TestOption{
				OptionText: o.OptionText,
				IsCorrect:  o.IsCorrect,
			})
		}
		test.Questions = append(test.Questions, question)
	}
	return test
}

func (s *Service) GenerateQuiz(ctx context.Context, req request.GenerateQuizReq) (*res.GeneratedTestResponse, error) {
	resp, err := s.aiAPI.GenerateQuiz(ctx, req, "ru")
	if err != nil {
		return nil, err
	}
	return resp, nil
}
