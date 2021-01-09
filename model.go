package main

//------------------------------------- MAIL -------------------------------------
type Mail struct {
	Header  map[string]string
	Content string
	Links   []string
}

//------------------------------------ COURSES -----------------------------------
type Calendar struct {
	CourseList map[string][]*Course
}

func (c *Calendar) AddCourse(date string, course *Course) {
	if _, exist := c.CourseList[date]; !exist {
		c.CourseList[date] = make([]*Course, 0)
	}
	c.CourseList[date] = append(c.CourseList[date], course)
}

func NewCalendar() *Calendar {
	return &Calendar{
		CourseList: make(map[string][]*Course),
	}
}

//------------------------------------ COURSE ------------------------------------
type Course struct {
	Start string
	Name  string
}

func NewCourse(start string, name string) *Course {
	return &Course{
		Start: start,
		Name:  name,
	}
}
