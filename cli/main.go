package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Status string
 
const (
	TODO Status = "TODO"
	INPROGRESS Status = "INPROGRESS"
	PENDING Status = "PENDING"
	COMPLETED Status = "COMPLETED"
)


type Task struct {
	ID string
	Title string
	Desc string
	Priority int
	Status Status
	ProjectID string
	CreatedAt time.Time
	ETA int
}

type Project struct {
	ID string
	Name string
	Tasks map[string]*Task
}

func NewProject(id, name string) *Project{
	return &Project{
		ID: id,
		Name: name,
		Tasks: map[string]*Task{},

	}
}

func (p *Project) AddTask(t *Task) {
	p.Tasks[t.ID] = t
}
func (p *Project) RemoveTask(TaskID string) {
	delete(p.Tasks, TaskID)
}
func (p *Project)TotalETA() int{
	sum := 0
	 for _,t := range p.Tasks {
		if t.Status != COMPLETED {
			sum += t.ETA
		}
	 }

	return sum
} 
func (t *Task) UpdateStatus(s Status) {
	t.Status =s
}

func(t *Task)UpdatePriority(p int) {
	t.Priority =p
}

type ProjectRepo interface{
	Save(p *Project) error
	Get(id string) (*Project, error)
}

type TaskRepo interface {
	Save(t *Task) error
	Get(id string)(*Task, error)
	Delete(id string) error
}

type InMemoryProjectRepo struct{ //Singleton
	projects map[string]*Project
	mu sync.RWMutex
}

func NewInMemoryProjectRepo() *InMemoryProjectRepo{
	return &InMemoryProjectRepo{
		projects: map[string]*Project{},
	}
}

func (r *InMemoryProjectRepo) Save(p *Project) error{
	r.mu.Lock()
	defer r.mu.Unlock()
	r.projects[p.ID] = p
	return nil
}

func (r *InMemoryProjectRepo) Get(id string) (*Project, error){
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	p, ok := r.projects[id]
	if !ok {
		return nil, errors.New("Project Not Found")
	}
	return p, nil
}

type InMemoryTaskRepo struct{
	mu sync.RWMutex
	tasks map[string]*Task
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
	return &InMemoryTaskRepo{tasks: map[string]*Task{}}
}

func (r *InMemoryTaskRepo) Save(t *Task) error{
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[t.ID] = t
	return nil
}

func(r *InMemoryTaskRepo) Get( id string) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.tasks[id]

	if !ok {
		return nil, errors.New("task Not found")
	}
	return t, nil
}

func(r *InMemoryTaskRepo) Delete(id string) error{
	r.mu.Lock()
	r.mu.Unlock()
	delete(r.tasks, id)
	return nil
}

type TaskManager struct{ //singleton
	projectRepo ProjectRepo
	taskRepo TaskRepo
	mu sync.Mutex
	projectSeq int
	taskSeq int
}
func NewTaskManager(pr ProjectRepo, tr TaskRepo) *TaskManager{
	return &TaskManager{
		projectRepo: pr,
		taskRepo: tr,
	}
}

func (m *TaskManager) nextProjectID() string{
	m.mu.Lock()
	defer m.mu.Unlock()
	m.projectSeq++
	return fmt.Sprintf("P%d", m.projectSeq)
}

func (m *TaskManager) nextTaskID() string{
	m.mu.Lock()
	defer m.mu.Unlock()
	m.taskSeq++
	return fmt.Sprintf("T%d", m.taskSeq)
}

func (m *TaskManager) CreateProject(name string) (*Project, error) {
	id:= m.nextProjectID()
	p := NewProject(id, name)

	if err := m.projectRepo.Save(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (m *TaskManager) Createtask(projectID, title, desc string, priority, eta int) (*Task, error) {
	p, err := m.projectRepo.Get(projectID)
	if err != nil {
		return nil, err
	}

	task := &Task{
		ID: m.nextTaskID(),
		Title: title,
		Desc: desc,
		Priority: priority,
		ETA: eta,
		Status: TODO,
		ProjectID: projectID,
		CreatedAt: time.Now(),
	}
	p.AddTask(task)
	_ = m.projectRepo.Save(p)
	_ = m.taskRepo.Save(task)
	return task, nil
}

func (m *TaskManager) MoveTask(taskID, toProjectID string) error {
	task, err := m.taskRepo.Get(taskID) 
	if err != nil {
		return err
	}

	fromProject, err := m.projectRepo.Get(task.ProjectID)
	if err != nil {
		return err
	}

	toProject, err := m.projectRepo.Get(toProjectID)
	if err != nil { return nil }

	fromProject.RemoveTask(taskID)
	_ = m.projectRepo.Save(fromProject)

	task.ProjectID = toProjectID
	toProject.AddTask(task)
	_ = m.projectRepo.Save(toProject)

	_ = m.taskRepo.Save(task)
	return nil
}



func main() {
	fmt.Println("cli")

	projectRepo := NewInMemoryProjectRepo()
	taskRepo := NewInMemoryTaskRepo()
	manager := NewTaskManager(projectRepo, taskRepo)

	p1, _ := manager.CreateProject("Roopam")
	p2, _ := manager.CreateProject("DBaaS")

	t1, _ := manager.Createtask(p1.ID, "Design DB Schema", "Table", 3, 10)
	// t2, _ := manager.Createtask(p1.ID, "Add auth", "RBAC", 5, 6)
	// t3, _ := manager.Createtask(p1.ID, "Write APIs", "CRUD", 2, 12)

	fmt.Println("Project: ", p1.ID, p1.Name)

	_ = manager.MoveTask(t1.ID, p2.ID)

}