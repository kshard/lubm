//
// Copyright (C) 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/lubm
//

package lubm

import (
	"encoding/json"
	"math/rand"
	"strconv"

	"github.com/kshard/spock"
	"github.com/kshard/spock/encoding/jsonld"
)

type DataSet struct {
	writer          chan<- spock.Bag
	rand            *rand.Rand
	maxUniversityID int
}

func NewDataSet(
	seed int64,
	maxUniversityID int,
	writer chan<- spock.Bag,
) *DataSet {
	rnd := rand.New(rand.NewSource(seed))
	rnd.Seed(seed)
	rand.Seed(seed)

	return &DataSet{
		rand:            rnd,
		maxUniversityID: maxUniversityID,
		writer:          writer,
	}
}

func (ds DataSet) Write(obj any) error {
	bin, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	var bag jsonld.Bag
	if err := json.Unmarshal(bin, &bag); err != nil {
		return err
	}

	ds.writer <- spock.Bag(bag)
	return nil
}

//
// See http://swat.cse.lehigh.edu/projects/lubm/profile.htm
//

func (dataset *DataSet) Generate(universityID int) error {
	return dataset.genUniversity(universityID, dataset.maxUniversityID)
}

func (dataset *DataSet) genUniversity(universityID, maxUniversityID int) error {
	university := newUniversity(universityID)
	if err := dataset.Write(university); err != nil {
		return err
	}

	// In each university
	// 15~25 Departments are subOrgnization of the University
	for id := 0; id < 15+dataset.rand.Intn(11); id++ {
		dept := newDepartment(university, id)
		if err := dataset.Write(dept); err != nil {
			return err
		}

		dataset.genDepartment(dept)
	}

	return nil
}

func (dataset *DataSet) genDepartment(dept *Department) error {
	faculties := make([]*Faculty, 0)

	// 7~10 FullProfessors worksFor the Department
	for id := 0; id < 7+dataset.rand.Intn(4); id++ {
		faculty := newProfessor(id, "Full", dept)
		faculties = append(faculties, faculty)
	}
	fullProfessors := len(faculties)

	// one of the FullProfessors is headOf the Department
	id := dataset.rand.Intn(len(faculties))
	faculties[id].HeadOf = (*IRI)(&dept.ID)

	// 10~14 AssociateProfessors worksFor the Department
	for id := 0; id < 10+dataset.rand.Intn(5); id++ {
		faculty := newProfessor(id, "Associate", dept)

		faculties = append(faculties, faculty)
	}
	associateProfessors := len(faculties)

	// 8~11 AssistantProfessors worksFor the Department
	for id := 0; id < 8+dataset.rand.Intn(4); id++ {
		faculty := newProfessor(id, "Assistant", dept)

		faculties = append(faculties, faculty)
	}
	assistantProfessors := len(faculties)

	// 5~7 Lecturers worksFor the Department
	for id := 0; id < 5+dataset.rand.Intn(3); id++ {
		faculty := newLecturer(id, dept)

		faculties = append(faculties, faculty)
	}

	// every Faculty has an undergraduateDegreeFrom a University, a mastersDegreeFrom a University, and a doctoralDegreeFrom a University
	for _, faculty := range faculties {
		faculty.UndergraduateDegreeFrom = dataset.degreeFromUniversity()
		faculty.MastersDegreeFrom = dataset.degreeFromUniversity()
		faculty.DoctoralDegreeFrom = dataset.degreeFromUniversity()
	}

	// UndergraduateStudent : Faculty = 8~14 : 1
	undergraduateStudents := make([]*Student, 0)
	for range faculties {
		for i := 0; i < 8+dataset.rand.Intn(7); i++ {
			student := newUndergraduateStudent(len(undergraduateStudents), dept)

			undergraduateStudents = append(undergraduateStudents, student)
		}
	}

	// 1/5 of the UndergraduateStudents have a Professor as their advisor
	for _, student := range dataset.fractionStudents(5, undergraduateStudents) {
		student.Advisor = dataset.professorID(faculties[:assistantProfessors])
	}

	// GraduateStudent : Faculty = 3~4 : 1
	graduateStudents := make([]*Student, 0)
	for range faculties {
		for i := 0; i < 3+dataset.rand.Intn(2); i++ {
			student := newGraduateStudent(len(graduateStudents), dept)
			// every GraudateStudent has an undergraduateDegreeFrom a University
			student.UndergraduateDegreeFrom = dataset.degreeFromUniversity()
			// every GraduateStudent has a Professor as his advisor
			student.Advisor = dataset.professorID(faculties[:assistantProfessors])

			graduateStudents = append(graduateStudents, student)
		}
	}

	// every Faculty is teacherOf 1~2 Courses
	courses := make([]*Course, 0)
	for _, faculty := range faculties {
		for i := 0; i < 1+rand.Intn(2); i++ {
			course := newCourse(len(courses), dept)
			faculty.TeacherOf = append(faculty.TeacherOf, IRI(course.ID))

			courses = append(courses, course)
		}
	}

	// every UndergraduateStudent takesCourse 2~4 Courses
	for _, student := range undergraduateStudents {
		student.TakesCourse = dataset.takesCourse(2, 4, courses)
	}

	// every Faculty is teacherOf 1~2 GraduateCourses
	graduateCourses := make([]*Course, 0)
	for _, faculty := range faculties {
		for i := 0; i < 1+rand.Intn(2); i++ {
			course := newGraduateCourse(len(graduateCourses), dept)
			faculty.TeacherOf = append(faculty.TeacherOf, IRI(course.ID))

			graduateCourses = append(graduateCourses, course)
		}
	}

	// every GraduateStudent takesCourse 1~3 GraduateCourses
	for _, student := range graduateStudents {
		student.TakesCourse = dataset.takesCourse(1, 3, graduateCourses)
	}

	// 1/5~1/4 of the GraduateStudents are chosen as TeachingAssistant for one Course
	for _, student := range dataset.fractionStudents(5, graduateStudents) {
		student.TeachingAssistantOf = dataset.courseID(courses)
	}

	// 1/4~1/3 of the GraduateStudents are chosen as ResearchAssistant

	publications := make([]*Publication, 0)

	// every FullProfessor is publicationAuthor of 15~20 Publications
	for _, professor := range faculties[:fullProfessors] {
		for i := 0; i < 15+rand.Intn(6); i++ {
			publication := newPublication(len(publications), dept, professor)

			publications = append(publications, publication)
		}
	}

	// every AssociateProfessor is publicationAuthor of 10~18 Publications
	for _, professor := range faculties[fullProfessors:associateProfessors] {
		for i := 0; i < 10+rand.Intn(9); i++ {
			publication := newPublication(len(publications), dept, professor)

			publications = append(publications, publication)
		}
	}

	// every AssistantProfessor is publicationAuthor of 5~10 Publications
	for _, professor := range faculties[associateProfessors:assistantProfessors] {
		for i := 0; i < 5+rand.Intn(6); i++ {
			publication := newPublication(len(publications), dept, professor)

			publications = append(publications, publication)
		}
	}

	// every Lecturer has 0~5 Publications
	for _, professor := range faculties[assistantProfessors:] {
		for i := 0; i < 0+rand.Intn(6); i++ {
			publication := newPublication(len(publications), dept, professor)

			publications = append(publications, publication)
		}
	}

	// every GraduateStudent co-authors 0~5 Publications with some Professors
	for _, student := range graduateStudents {
		for _, publication := range dataset.publications(6, publications) {
			publication.PublicationAuthor = append(publication.PublicationAuthor, IRI(student.ID))
		}
	}

	// 10~20 ResearchGroups are subOrgnization of the Department
	researchGroups := make([]*ResearchGroup, 0)
	for i := 0; i < 10+rand.Intn(21); i++ {
		researchGroup := newResearchGroup(dept, i)

		researchGroups = append(researchGroups, researchGroup)
	}

	if err := dataset.Write(faculties); err != nil {
		return err
	}
	if err := dataset.Write(undergraduateStudents); err != nil {
		return err
	}
	if err := dataset.Write(graduateStudents); err != nil {
		return err
	}
	if err := dataset.Write(courses); err != nil {
		return err
	}
	if err := dataset.Write(graduateCourses); err != nil {
		return err
	}
	if err := dataset.Write(publications); err != nil {
		return err
	}
	if err := dataset.Write(researchGroups); err != nil {
		return err
	}

	return nil
}

func (dataset *DataSet) degreeFromUniversity() *IRI {
	id := dataset.rand.Intn(dataset.maxUniversityID)
	iri := IRI("edu:University" + strconv.Itoa(id))
	return &iri
}

func (dataset *DataSet) takesCourse(min, max int, courses []*Course) []IRI {
	set := map[IRI]struct{}{}

	n := min + dataset.rand.Intn(max-min)
	for i := 0; i < n; i++ {
		course := courses[dataset.rand.Intn(len(courses))]
		set[IRI(course.ID)] = struct{}{}
	}

	seq := make([]IRI, 0)
	for iri := range set {
		seq = append(seq, iri)
	}
	return seq
}

func (dataset *DataSet) professorID(faculties []*Faculty) *IRI {
	id := dataset.rand.Intn(len(faculties))
	iri := IRI(faculties[id].ID)
	return &iri
}

func (dataset *DataSet) courseID(courses []*Course) *IRI {
	id := dataset.rand.Intn(len(courses))
	iri := IRI(courses[id].ID)
	return &iri
}

func (dataset *DataSet) fractionStudents(n int, students []*Student) []*Student {
	set := map[UID]*Student{}

	for i := 0; i < len(students)/5; i++ {
		student := students[dataset.rand.Intn(len(students))]
		set[student.ID] = student
	}

	seq := make([]*Student, 0)
	for _, student := range set {
		seq = append(seq, student)
	}
	return seq
}

func (dataset *DataSet) publications(n int, publications []*Publication) []*Publication {
	set := map[UID]*Publication{}

	for i := 0; i < dataset.rand.Intn(n); i++ {
		publication := publications[dataset.rand.Intn(len(publications))]
		set[publication.ID] = publication
	}

	seq := make([]*Publication, 0)
	for _, publication := range set {
		seq = append(seq, publication)
	}
	return seq
}

func newUniversity(id int) *University {
	name := "University" + strconv.Itoa(id)

	return &University{
		ID:   UID("edu:" + name),
		Type: UID("ub:University"),
		Name: name,
	}
}

func newDepartment(u *University, id int) *Department {
	name := "Department" + strconv.Itoa(id)

	return &Department{
		ID:                u.ID + UID("."+name),
		Type:              UID("ub:Department"),
		Name:              name,
		SubOrganizationOf: IRI(u.ID),
	}
}

func newResearchGroup(dept *Department, id int) *ResearchGroup {
	name := "ResearchGroup" + strconv.Itoa(id)

	return &ResearchGroup{
		ID:                dept.ID + UID("/"+name),
		Type:              UID("ub:ResearchGroup"),
		SubOrganizationOf: IRI(dept.ID),
	}
}

func newProfessor(id int, kind string, dept *Department) *Faculty {
	name := kind + "Professor" + strconv.Itoa(id)

	return &Faculty{
		ID:               dept.ID + UID("/"+name),
		Type:             UID("ub:" + kind + "Professor"),
		Name:             name,
		TeacherOf:        []IRI{},
		WorksFor:         IRI(dept.ID),
		EmailAddress:     string(dept.ID) + "@" + name,
		Telephone:        telephone(),
		ResearchInterest: "Research0",
	}
}

func newLecturer(id int, dept *Department) *Faculty {
	name := "Lecturer" + strconv.Itoa(id)

	return &Faculty{
		ID:               dept.ID + UID("/"+name),
		Type:             UID("ub:Lecturer"),
		Name:             name,
		TeacherOf:        []IRI{},
		WorksFor:         IRI(dept.ID),
		EmailAddress:     string(dept.ID) + "@" + name,
		Telephone:        telephone(),
		ResearchInterest: "Research0",
	}
}

func newCourse(id int, dept *Department) *Course {
	name := "Course" + strconv.Itoa(id)

	return &Course{
		ID:   dept.ID + UID("/"+name),
		Type: UID("ub:Course"),
		Name: name,
	}
}

func newGraduateCourse(id int, dept *Department) *Course {
	name := "GraduateCourse" + strconv.Itoa(id)

	return &Course{
		ID:   dept.ID + UID("/"+name),
		Type: UID("ub:GraduateCourse"),
		Name: name,
	}
}

func newUndergraduateStudent(id int, dept *Department) *Student {
	name := "UndergraduateStudent" + strconv.Itoa(id)

	return &Student{
		ID:           dept.ID + UID("/"+name),
		Type:         UID("ub:UndergraduateStudent"),
		Name:         name,
		MemberOf:     IRI(dept.ID),
		EmailAddress: string(dept.ID) + "@" + name,
		Telephone:    telephone(),
		TakesCourse:  []IRI{},
	}
}

func newGraduateStudent(id int, dept *Department) *Student {
	name := "GraduateStudent" + strconv.Itoa(id)

	return &Student{
		ID:           dept.ID + UID("/"+name),
		Type:         UID("ub:GraduateStudent"),
		Name:         name,
		MemberOf:     IRI(dept.ID),
		EmailAddress: string(dept.ID) + "@" + name,
		Telephone:    telephone(),
		TakesCourse:  []IRI{},
	}
}

func newPublication(id int, dept *Department, faculty *Faculty) *Publication {
	name := "Publication" + strconv.Itoa(id)

	return &Publication{
		ID:                dept.ID + UID("/"+faculty.Name+"/"+name),
		Type:              UID("ub:Publication"),
		Name:              name,
		PublicationAuthor: []IRI{IRI(faculty.ID)},
	}
}

func telephone() string {
	d3 := func() string {
		return strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10)) + strconv.Itoa(rand.Intn(10))
	}

	return d3() + "-" + d3() + "-" + d3()
}
