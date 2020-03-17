package db

import (
	"fmt"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/suite"
	"strconv"
	"testing"
)

var tasks *Tasks

type TaskTestSuite struct {
	suite.Suite
	unit   *Tasks
	gomega *GomegaWithT
}

func Test_TaskTestSuite(t *testing.T) {
	var err error
	if tasks, err = Create(); err != nil {
		panic(fmt.Errorf("failed to initialize db: %s", err))
	}

	testSuite := TaskTestSuite{unit: tasks, gomega: NewGomegaWithT(t)}
	suite.Run(t, &testSuite)
}

func (s *TaskTestSuite) Test_InsertedTasksCanBeQueried() {
	var err error
	if err = s.unit.Drop(); err != nil {
		panic(fmt.Errorf("failed to remove existing bucket: %s", err))
	}
	if err = s.unit.Insert("task 1"); err != nil {
		panic(fmt.Errorf("failed to insert given task: %s", err))
	}
	if err = s.unit.Insert("task 2"); err != nil {
		panic(fmt.Errorf("failed to insert given task"))
	}

	var taskList []Task
	if taskList, err = tasks.QueryAll(); err != nil {
		panic(fmt.Errorf("failed to retrieve task associated with the given id: %s", err))
	}

	for i, task := range taskList {
		fmt.Println(strconv.Itoa(int(task.Id)) + " " + task.Description + " " + string(task.Status))

		s.gomega.Expect(task.Id).Should(Equal(int64(i + 1)))
		s.gomega.Expect(task.Description).Should(Equal("task " + strconv.Itoa(i+1)))
		s.gomega.Expect(task.Status).Should(Equal(Created))
	}
}

func (s *TaskTestSuite) Test_TasksAreCorrectlyUpdated() {
	var err error
	if err = s.unit.Drop(); err != nil {
		panic(fmt.Errorf("failed to remove existing bucket: %s", err))
	}
	if err = s.unit.Insert("task 1"); err != nil {
		panic(fmt.Errorf("failed to insert given task: %s", err))
	}
	if err = s.unit.Insert("task 2"); err != nil {
		panic(fmt.Errorf("failed to insert given task: %s", err))
	}
	if _, err = s.unit.Update(1, Completed); err != nil {
		panic(fmt.Errorf("failed to update given task: %s", err))
	}

	var tasks []Task
	if tasks, err = s.unit.QueryAll(); err != nil {
		panic(fmt.Errorf("failed to query all tasks: %s", err))
	}

	for i, task := range tasks {
		s.gomega.Expect(task.Id).Should(Equal(int64(i + 1)))
		s.gomega.Expect(task.Description).Should(Equal("task " + strconv.Itoa(i+1)))
		if i == 0 {
			s.gomega.Expect(task.Status).Should(Equal(Completed))
		} else {
			s.gomega.Expect(task.Status).Should(Equal(Created))
		}
	}
}

func (s *TaskTestSuite) Test_NoErrorsWhenQueryingAndNoTasks() {
	var tasks []Task
	var err error
	if err = s.unit.DeleteAll(); err != nil {
		panic(fmt.Errorf("failed to remove existing bucket: %s", err))
	}
	if tasks, err = s.unit.QueryAll(); err != nil {
		panic(fmt.Errorf("couldn't query tasks: %s", err))
	}
	s.gomega.Expect(len(tasks)).Should(Equal(0))

}
