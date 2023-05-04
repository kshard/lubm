//
// Copyright (C) 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/lubm
//

package lubm

import "fmt"

//
// See http://swat.cse.lehigh.edu/projects/lubm/queries-sparql.txt
// See http://swat.cse.lehigh.edu/projects/lubm/lubm.jpg
//

// # Query1
// # This query bears large input and high selectivity. It queries about just one class and
// # one property and does not assume any hierarchy information or inference.
//
//	{
//	 ?X rdf:type ub:GraduateStudent .
//	 ?X ub:takesCourse http://www.Department0.University0.edu/GraduateCourse0
//	}
func Query1(course ...string) string {
	c := "edu:University0.Department0/GraduateCourse5"
	if len(course) != 0 {
		c = course[0]
	}

	return fmt.Sprintf(`
		f(s, p, o).

		q(x) :-
			f(x, ub:takesCourse, <%s>),
			f(x, rdf:type, ub:GraduateStudent).
	`, c)
}

// # Query2
// # This query increases in complexity: 3 classes and 3 properties are involved. Additionally,
// # there is a triangular pattern of relationships between the objects involved.
//
//	{
//	  ?X rdf:type ub:GraduateStudent .
//	  ?Y rdf:type ub:University .
//	  ?Z rdf:type ub:Department .
//	  ?X ub:memberOf ?Z .
//	  ?Z ub:subOrganizationOf ?Y .
//	  ?X ub:undergraduateDegreeFrom ?Y
//	}
func Query2() string {
	return `
		f(s, p, o).

		q(y, z, x) :-
			f(y, rdf:type, ub:University),

			f(x, ub:undergraduateDegreeFrom, y),
			f(x, rdf:type, ub:GraduateStudent),
			f(x, ub:memberOf, z),

			f(z, rdf:type, ub:Department),
			f(z, ub:subOrganizationOf, y).
	`
}

// # Query3
// # This query is similar to Query 1 but class Publication has a wide hierarchy.
//
//	{
//		?X rdf:type ub:Publication .
//	  ?X ub:publicationAuthor http://www.Department0.University0.edu/AssistantProfessor0
//	}
func Query3(author ...string) string {
	a := "edu:University0.Department0/AssistantProfessor0"
	if len(author) != 0 {
		a = author[0]
	}

	return fmt.Sprintf(`
		f(s, p, o).

		q(x) :-
			f(x, ub:publicationAuthor, <%s>),
			f(x, rdf:type, ub:Publication).
	`, a)
}

// # Query4
// # This query has small input and high selectivity. It assumes subClassOf relationship
// # between Professor and its subclasses. Class Professor has a wide hierarchy. Another
// # feature is that it queries about multiple properties of a single class.
//
//	{
//		?X rdf:type ub:Professor .
//	  ?X ub:worksFor <http://www.Department0.University0.edu> .
//	  ?X ub:name ?Y1 .
//	  ?X ub:emailAddress ?Y2 .
//	  ?X ub:telephone ?Y3
//	}
func Query4(dept ...string) string {
	d := "edu:University0.Department0"
	if len(dept) != 0 {
		d = dept[0]
	}

	return fmt.Sprintf(`
		f(s, p, o).

		q(x, name, email, phone) :-
			f(x, ub:worksFor, <%s>),
			f(x, ub:name, name),
			f(x, ub:emailAddress, email),
			f(x, ub:telephone, phone).
	`, d)
}

// # Query5
// # This query assumes subClassOf relationship between Person and its subclasses
// # and subPropertyOf relationship between memberOf and its subproperties.
// # Moreover, class Person features a deep and wide hierarchy.
//
//	{
//		?X rdf:type ub:Person .
//	  ?X ub:memberOf <http://www.Department0.University0.edu>
//	}
func Query5(dept ...string) string {
	d := "edu:University0.Department0"
	if len(dept) != 0 {
		d = dept[0]
	}

	return fmt.Sprintf(`
		f(s, p, o).

		q(x) :-
			f(x, ub:memberOf, <%s>),
			f(x, rdf:type, ub:UndergraduateStudent).
	`, d)
}

// # Query6
// # This query queries about only one class. But it assumes both the explicit
// # subClassOf relationship between UndergraduateStudent and Student and the
// # implicit one between GraduateStudent and Student. In addition, it has large
// # input and low selectivity.
//
//	{?X rdf:type ub:Student}
func Query6() string {
	return `
		f(s, p, o).

		q(x) :-
			f(x, rdf:type, ub:UndergraduateStudent).
	`
}

// # Query7
// # This query is similar to Query 6 in terms of class Student but it increases in the
// # number of classes and properties and its selectivity is high.
//
//	{
//		?X rdf:type ub:Student .
//		?Y rdf:type ub:Course .
//		?X ub:takesCourse ?Y .
//		<http://www.Department0.University0.edu/AssociateProfessor0> ub:teacherOf, ?Y
//	}
func Query7(teacher ...string) string {
	t := "edu:University0.Department0/AssistantProfessor0"
	if len(teacher) != 0 {
		t = teacher[0]
	}

	return fmt.Sprintf(`
		f(s, p, o).

		q(x, y) :-
			f(<%s>, ub:teacherOf, y),
			f(y, rdf:type, ub:Course),

			f(x, ub:takesCourse, y),
			f(x, rdf:type, ub:UndergraduateStudent).
	`, t)
}

// # Query8
// # This query is further more complex than Query 7 by including one more property.
//
//	{
//		?X rdf:type ub:Student .
//	  ?Y rdf:type ub:Department .
//	  ?X ub:memberOf ?Y .
//	  ?Y ub:subOrganizationOf <http://www.University0.edu> .
//	  ?X ub:emailAddress ?Z
//	}
func Query8(university ...string) string {
	u := "edu:University0"
	if len(university) != 0 {
		u = university[0]
	}

	return fmt.Sprintf(`
		f(s, p, o).

		q(y, x, email) :-
			f(y, ub:subOrganizationOf, <%s>),
			f(y, rdf:type, ub:Department),

			f(x, ub:memberOf, y),
			f(x, ub:emailAddress, email).
	`, u)
}

// # Query9
// # Besides the aforementioned features of class Student and the wide hierarchy of
// # class Faculty, like Query 2, this query is characterized by the most classes and
// # properties in the query set and there is a triangular pattern of relationships.
//
//	{
//		?X rdf:type ub:Student .
//	  ?Y rdf:type ub:Faculty .
//	  ?Z rdf:type ub:Course .
//	  ?X ub:advisor ?Y .
//	  ?Y ub:teacherOf ?Z .
//	  ?X ub:takesCourse ?Z
//	}
func Query9() string {
	return `
		f(s, p, o).

		q(x) :-
			f(x, ub:advisor, y),
			f(y, ub:teacherOf, z),
			f(x, ub:takesCourse, z).
	`
}
