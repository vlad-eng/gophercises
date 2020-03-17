package db

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/boltdb/bolt"
	"strings"
)

type Status string

const (
	Completed Status = "Completed"
	Created   Status = "Created"
)

type Tasks struct {
	boltDB *bolt.DB
}

type Task struct {
	Id          int64
	Description string
	Status      Status
}

const TaskList = "TaskList"
const columnSeparator = "|"

func Create() (*Tasks, error) {
	var boltDB *bolt.DB
	var err error
	if boltDB, err = bolt.Open("tasks.db", 0600, nil); err != nil {
		return nil, err
	}

	tx, err := boltDB.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.CreateBucketIfNotExists([]byte(TaskList)); err != nil {
		return nil, err
	}

	return &Tasks{boltDB: boltDB}, tx.Commit()
}

func (db *Tasks) Insert(taskDescription string) error {
	task := Task{
		Id:          -1,
		Description: taskDescription,
		Status:      Created,
	}

	return db.updateBucket(task)
}

func (db *Tasks) Update(id int64, fieldToUpdate interface{}) (*Task, error) {
	var status Status
	var ok bool
	status, ok = fieldToUpdate.(Status)
	if !ok {
		return nil, fmt.Errorf("unrecognized field type")
	}

	var task *Task
	var err error
	if task, err = db.Query(id); err != nil {
		return nil, err
	}
	task.Status = status
	if err = db.updateBucket(*task); err != nil {
		return nil, err
	}

	return task, nil
}

func (db *Tasks) updateBucket(task Task) error {
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		var bucket *bolt.Bucket
		var err error
		if bucket, err = tx.CreateBucketIfNotExists([]byte(TaskList)); err != nil {
			return err
		}

		var taskId uint64
		if task.Id == -1 {
			taskId, err = bucket.NextSequence()
		} else {
			taskId = uint64(task.Id)
		}
		taskContent := append([]byte(task.Description), []byte(columnSeparator)...)
		taskContent = append(taskContent, []byte(string(task.Status))...)

		taskIdBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(taskIdBytes, taskId)

		if err = bucket.Put(taskIdBytes, taskContent); err != nil {
			return err
		}

		return nil
	})

	return err
}

func (db *Tasks) QueryAll() ([]Task, error) {
	tasks := make([]Task, 0)
	var taskId uint64
	var taskContent []string

	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(TaskList))
		err := bucket.ForEach(func(key, value []byte) error {
			var e error
			taskId, e = binary.ReadUvarint(bytes.NewBuffer(key))
			taskContent = strings.Split(string(value), columnSeparator)
			task := Task{Id: int64(taskId), Description: taskContent[0], Status: Status(taskContent[1])}
			tasks = append(tasks, task)
			return e
		})

		return err
	})
	return tasks, err
}

func (db *Tasks) Query(id int64) (*Task, error) {
	var task Task
	var taskContent []string
	err := db.boltDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(TaskList))
		if bucket != nil {
			key := make([]byte, 8)
			binary.LittleEndian.PutUint64(key, uint64(id))
			value := bucket.Get(key)
			if value != nil {
				taskContent = strings.Split(string(value), columnSeparator)
				task = Task{Id: id, Description: taskContent[0], Status: Status(taskContent[1])}
				return nil
			} else {
				return fmt.Errorf("couldn't find task with id: %d", id)
			}
		}
		return fmt.Errorf("bucket %s hasn't been created", TaskList)
	})

	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (db *Tasks) Drop() error {
	var err error
	err = db.boltDB.Update(func(tx *bolt.Tx) error {
		err = tx.DeleteBucket([]byte(TaskList))
		if err != nil && strings.Compare(err.Error(), "bucket not found") == 0 {
			return nil
		}
		return err
	})

	return err
}

func (db *Tasks) Delete(id int64) error {
	err := db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(TaskList))
		if bucket != nil {
			key := make([]byte, 8)
			binary.LittleEndian.PutUint64(key, uint64(id))
			value := bucket.Get(key)
			if value != nil {
				return bucket.Delete(key)
			} else {
				return fmt.Errorf("couldn't find task with id: %d", id)
			}
		}
		return fmt.Errorf("bucket %s hasn't been created", TaskList)
	})
	return err
}

func (db *Tasks) DeleteAll() error {
	var err error
	err = db.boltDB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(TaskList))
		if bucket != nil {
			err := bucket.ForEach(func(key, value []byte) error {
				return bucket.Delete(key)
			})
			return err
		}
		return fmt.Errorf("bucket %s hasn't been created", TaskList)
	})
	return err
}
