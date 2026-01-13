package main

import "fmt"

// ---------- ENUMS ----------

type Difficulty int

const (
	EASY Difficulty = iota
	MEDIUM
	HARD
)

type SubmissionStatus int

const (
	PENDING SubmissionStatus = iota
	ACCEPTED
	WRONG_ANSWER
)

type Language int

const (
	JAVA Language = iota
	CPP
	PYTHON
	GO
)

// ---------- USER ----------

type User struct {
	ID          int
	Name        string
	Submissions []*Submission
}

func NewUser(id int, name string) *User {
	return &User{
		ID:   id,
		Name: name,
	}
}

func (u *User) Submit(problem *Problem, code string, lang Language) *Submission {
	sub := &Submission{
		ID:       len(u.Submissions) + 1,
		User:     u,
		Problem:  problem,
		Code:     code,
		Language: lang,
		Status:   PENDING,
	}
	u.Submissions = append(u.Submissions, sub)
	return sub
}

// ---------- PROBLEM ----------

type Problem struct {
	ID          int
	Title       string
	Description string
	Difficulty  Difficulty
	TestCases   []*TestCase
}

func NewProblem(id int, title, desc string, diff Difficulty) *Problem {
	return &Problem{
		ID:          id,
		Title:       title,
		Description: desc,
		Difficulty:  diff,
	}
}

// ---------- TEST CASE ----------

type TestCase struct {
	Input          string
	ExpectedOutput string
}

// ---------- SUBMISSION ----------

type Submission struct {
	ID       int
	User     *User
	Problem  *Problem
	Code     string
	Language Language
	Status   SubmissionStatus
}

// ---------- EXECUTOR ----------

type CodeExecutor struct{}

func (ce *CodeExecutor) Run(sub *Submission) {
	if len(sub.Code) > 10 {
		sub.Status = ACCEPTED
	} else {
		sub.Status = WRONG_ANSWER
	}
}

// ---------- MAIN ----------

func main() {

	executor := &CodeExecutor{}

	problem := NewProblem(1, "Two Sum", "Find two numbers", EASY)
	problem.TestCases = []*TestCase{{"2 7 11 15", "9"}}

	user1 := NewUser(1, "Naman")
	user2 := NewUser(2, "Rahul")

	sub1 := user1.Submit(problem, "valid solution code here", GO)
	sub2 := user2.Submit(problem, "bad", GO)

	executor.Run(sub1)
	executor.Run(sub2)

	fmt.Println("Results:")
	fmt.Println(user1.Name, ":", sub1.Status)
	fmt.Println(user2.Name, ":", sub2.Status)
}
