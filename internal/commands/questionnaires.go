package commands

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/rzolkos/basecamp-cli/internal/client"
)

// Questionnaire represents a Basecamp questionnaire (automatic check-ins)
type Questionnaire struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	QuestionsURL string `json:"questions_url"`
}

// Question represents a check-in question
type Question struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	Schedule      string  `json:"schedule"`
	Paused        bool    `json:"paused"`
	AnswersURL    string  `json:"answers_url"`
	CommentsCount int     `json:"comments_count"`
	CommentsURL   string  `json:"comments_url"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Creator       Creator `json:"creator"`
	URL           string  `json:"app_url"`
}

// QuestionAnswer represents an answer to a question
type QuestionAnswer struct {
	ID            int     `json:"id"`
	Content       string  `json:"content"`
	GroupOn       string  `json:"group_on"`
	CommentsCount int     `json:"comments_count"`
	CommentsURL   string  `json:"comments_url"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	Creator       Creator `json:"creator"`
	URL           string  `json:"app_url"`
}

// fetchQuestionnaire gets the questionnaire for a project
func fetchQuestionnaire(cl *client.Client, projectID string) (ProjectDetail, Questionnaire, error) {
	project, err := fetchProject(cl, projectID)
	if err != nil {
		return ProjectDetail{}, Questionnaire{}, err
	}

	questionnaireURL, err := getDockURL(project, "questionnaire")
	if err != nil {
		return ProjectDetail{}, Questionnaire{}, err
	}

	data, err := cl.Get(questionnaireURL)
	if err != nil {
		return ProjectDetail{}, Questionnaire{}, err
	}

	var questionnaire Questionnaire
	if err := json.Unmarshal(data, &questionnaire); err != nil {
		return ProjectDetail{}, Questionnaire{}, err
	}

	return project, questionnaire, nil
}

// QuestionnaireCmd gets the questionnaire for a project
type QuestionnaireCmd struct{}

type QuestionnaireOutput struct {
	ProjectID       int    `json:"project_id"`
	ProjectName     string `json:"project_name"`
	QuestionnaireID int    `json:"questionnaire_id"`
	Title           string `json:"title"`
}

func (c *QuestionnaireCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	project, questionnaire, err := fetchQuestionnaire(cl, projectID)
	if err != nil {
		return err
	}

	return PrintJSON(QuestionnaireOutput{
		ProjectID:       project.ID,
		ProjectName:     project.Name,
		QuestionnaireID: questionnaire.ID,
		Title:           questionnaire.Title,
	})
}

// QuestionsCmd lists questions in a questionnaire
type QuestionsCmd struct{}

type QuestionsOutput struct {
	ProjectID       int             `json:"project_id"`
	QuestionnaireID int             `json:"questionnaire_id"`
	Questions       []QuestionBrief `json:"questions"`
}

type QuestionBrief struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Schedule string `json:"schedule"`
	Paused   bool   `json:"paused"`
}

func (c *QuestionsCmd) Run(args []string) error {
	projectID, _, err := getProjectID(args)
	if err != nil {
		return err
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	project, questionnaire, err := fetchQuestionnaire(cl, projectID)
	if err != nil {
		return err
	}

	data, err := cl.Get(questionnaire.QuestionsURL)
	if err != nil {
		return err
	}

	var questions []Question
	if err := json.Unmarshal(data, &questions); err != nil {
		return err
	}

	output := QuestionsOutput{
		ProjectID:       project.ID,
		QuestionnaireID: questionnaire.ID,
		Questions:       make([]QuestionBrief, len(questions)),
	}

	for i, q := range questions {
		output.Questions[i] = QuestionBrief{
			ID:       q.ID,
			Title:    q.Title,
			Schedule: q.Schedule,
			Paused:   q.Paused,
		}
	}

	return PrintJSON(output)
}

// QuestionCmd views a single question
type QuestionCmd struct{}

type QuestionDetailOutput struct {
	ID            int             `json:"id"`
	Title         string          `json:"title"`
	Schedule      string          `json:"schedule"`
	Paused        bool            `json:"paused"`
	Creator       string          `json:"creator"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	CommentsCount int             `json:"comments_count"`
	URL           string          `json:"url"`
	Comments      []CommentOutput `json:"comments,omitempty"`
}

func (c *QuestionCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse question_id and flags
	var questionID string
	showComments := false

	for i := 0; i < len(remaining); i++ {
		if remaining[i] == "--comments" {
			showComments = true
		} else if questionID == "" {
			questionID = remaining[i]
		}
	}

	if questionID == "" {
		return errors.New("question_id required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/buckets/" + projectID + "/questions/" + questionID + ".json")
	if err != nil {
		return err
	}

	var question Question
	if err := json.Unmarshal(data, &question); err != nil {
		return err
	}

	output := QuestionDetailOutput{
		ID:            question.ID,
		Title:         question.Title,
		Schedule:      question.Schedule,
		Paused:        question.Paused,
		Creator:       question.Creator.Name,
		CreatedAt:     question.CreatedAt,
		UpdatedAt:     question.UpdatedAt,
		CommentsCount: question.CommentsCount,
		URL:           question.URL,
	}

	if showComments && question.CommentsURL != "" {
		comments, err := fetchComments(cl, question.CommentsURL)
		if err != nil {
			return err
		}
		output.Comments = comments
	}

	return PrintJSON(output)
}

// QuestionAnswersCmd lists answers to a question
type QuestionAnswersCmd struct{}

type QuestionAnswersOutput struct {
	QuestionID int                   `json:"question_id"`
	Answers    []QuestionAnswerBrief `json:"answers"`
}

type QuestionAnswerBrief struct {
	ID        int    `json:"id"`
	Content   string `json:"content"`
	GroupOn   string `json:"group_on"`
	Creator   string `json:"creator"`
	CreatedAt string `json:"created_at"`
}

func (c *QuestionAnswersCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	if len(remaining) < 1 {
		return errors.New("question_id required")
	}
	questionID := remaining[0]

	cl, err := client.New()
	if err != nil {
		return err
	}

	// First get the question to get its answers URL
	questionData, err := cl.Get("/buckets/" + projectID + "/questions/" + questionID + ".json")
	if err != nil {
		return err
	}

	var question Question
	if err := json.Unmarshal(questionData, &question); err != nil {
		return err
	}

	// Fetch answers
	answersData, err := cl.GetAll(question.AnswersURL)
	if err != nil {
		return err
	}

	answers := make([]QuestionAnswerBrief, len(answersData))
	for i, answerJSON := range answersData {
		var answer QuestionAnswer
		if err := json.Unmarshal(answerJSON, &answer); err != nil {
			return err
		}
		answers[i] = QuestionAnswerBrief{
			ID:        answer.ID,
			Content:   stripHTML(answer.Content),
			GroupOn:   answer.GroupOn,
			Creator:   answer.Creator.Name,
			CreatedAt: answer.CreatedAt,
		}
	}

	// Parse question ID for output
	var qID int
	fmt.Sscanf(questionID, "%d", &qID)

	return PrintJSON(QuestionAnswersOutput{
		QuestionID: qID,
		Answers:    answers,
	})
}

// QuestionAnswerCmd views a single answer
type QuestionAnswerCmd struct{}

type QuestionAnswerDetailOutput struct {
	ID            int             `json:"id"`
	Content       string          `json:"content"`
	GroupOn       string          `json:"group_on"`
	Creator       string          `json:"creator"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
	CommentsCount int             `json:"comments_count"`
	URL           string          `json:"url"`
	Comments      []CommentOutput `json:"comments,omitempty"`
}

func (c *QuestionAnswerCmd) Run(args []string) error {
	projectID, remaining, err := getProjectID(args)
	if err != nil {
		return err
	}

	// Parse answer_id and flags
	var answerID string
	showComments := false

	for i := 0; i < len(remaining); i++ {
		if remaining[i] == "--comments" {
			showComments = true
		} else if answerID == "" {
			answerID = remaining[i]
		}
	}

	if answerID == "" {
		return errors.New("answer_id required")
	}

	cl, err := client.New()
	if err != nil {
		return err
	}

	data, err := cl.Get("/buckets/" + projectID + "/question_answers/" + answerID + ".json")
	if err != nil {
		return err
	}

	var answer QuestionAnswer
	if err := json.Unmarshal(data, &answer); err != nil {
		return err
	}

	output := QuestionAnswerDetailOutput{
		ID:            answer.ID,
		Content:       stripHTML(answer.Content),
		GroupOn:       answer.GroupOn,
		Creator:       answer.Creator.Name,
		CreatedAt:     answer.CreatedAt,
		UpdatedAt:     answer.UpdatedAt,
		CommentsCount: answer.CommentsCount,
		URL:           answer.URL,
	}

	if showComments && answer.CommentsURL != "" {
		comments, err := fetchComments(cl, answer.CommentsURL)
		if err != nil {
			return err
		}
		output.Comments = comments
	}

	return PrintJSON(output)
}
