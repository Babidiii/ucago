package main

import "fmt"

//------------------------------------- MAIL -------------------------------------
type Mail struct {
	Header  map[string]string
	Content string
	Links   []string
}

//------------------------------------ COURSES -----------------------------------
type Calendar struct {
	CourseList map[string]map[string]*Course
}

func (c *Calendar) AddCourse(date string, course *Course) {
	if _, exist := c.CourseList[date]; !exist {
		c.CourseList[date] = make(map[string]*Course)
	}
	c.CourseList[date][course.Start] = course
}

func NewCalendar() *Calendar {
	return &Calendar{
		CourseList: make(map[string]map[string]*Course),
	}
}

//------------------------------------ COURSE ------------------------------------
type Course struct {
	Start string
	Name  string
	Link  string
	Info  map[string]string
	End   string
}

func NewCourse(start string, name string) *Course {
	return &Course{
		Start: start,
		Name:  name,
	}
}

func (c *Course) Display() {
	fmt.Printf("\t%s\n", c.Name)
	fmt.Printf("\tStart: %s End: %s\n", c.Start, c.End)
	fmt.Printf("\tLink:\n\t%s\n", c.Link)
	//	fmt.Println("\tInfo:")
}
