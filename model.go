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

func (c *Calendar) AddCourse(course *Course) {
	if _, exist := c.CourseList[course.Date]; !exist {
		c.CourseList[course.Date] = make([]*Course, 1)
	}
	c.CourseList[courseList.Date] = append(c.CourseList, course)
}

func NewCalendar() {
	return &Calendar{
		CourseList: make(map[string][]*Course),
	}
}

//------------------------------------ COURSE ------------------------------------
type Course struct {
	Date  string
	Start string
	Name  string
}

func (c *Course) NewCourse(date string, start string, name string) {
	return &Course{
		Date:  date,
		Start: start,
		Name:  name,
	}
}
