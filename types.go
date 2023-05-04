//
// Copyright (C) 2023 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/lubm
//

package lubm

import "github.com/fogfish/curie"

// Unique Identity
type UID curie.IRI

func (iri UID) MarshalJSON() ([]byte, error) {
	return []byte(`"` + iri + `"`), nil
}

// Cross-references
type IRI curie.IRI

func (iri IRI) MarshalJSON() ([]byte, error) {
	return []byte(`{"@id": "` + iri + `"}`), nil
}

type University struct {
	ID   UID    `json:"@id"`
	Type UID    `json:"@type"`
	Name string `json:"ub:name"`
}

type Department struct {
	ID                UID    `json:"@id"`
	Type              UID    `json:"@type"`
	Name              string `json:"ub:name"`
	SubOrganizationOf IRI    `json:"ub:subOrganizationOf"`
}

type Faculty struct {
	ID                      UID    `json:"@id"`
	Type                    UID    `json:"@type"`
	Name                    string `json:"ub:name"`
	HeadOf                  *IRI   `json:"ub:headOf,omitempty"`
	TeacherOf               []IRI  `json:"ub:teacherOf"`
	UndergraduateDegreeFrom *IRI   `json:"ub:undergraduateDegreeFrom,omitempty"`
	MastersDegreeFrom       *IRI   `json:"ub:mastersDegreeFrom,omitempty"`
	DoctoralDegreeFrom      *IRI   `json:"ub:doctoralDegreeFrom,omitempty"`
	WorksFor                IRI    `json:"ub:worksFor"`
	EmailAddress            string `json:"ub:emailAddress"`
	Telephone               string `json:"ub:telephone"`
	ResearchInterest        string `json:"ub:researchInterest"`
}

type Student struct {
	ID                      UID    `json:"@id"`
	Type                    UID    `json:"@type"`
	Name                    string `json:"ub:name"`
	MemberOf                IRI    `json:"ub:memberOf"`
	EmailAddress            string `json:"ub:emailAddress"`
	Telephone               string `json:"ub:telephone"`
	TakesCourse             []IRI  `json:"ub:takesCourse"`
	UndergraduateDegreeFrom *IRI   `json:"ub:undergraduateDegreeFrom,omitempty"`
	Advisor                 *IRI   `json:"ub:advisor,omitempty"`
	TeachingAssistantOf     *IRI   `json:"ub:teachingAssistantOf,omitempty"`
}

type Course struct {
	ID   UID    `json:"@id"`
	Type UID    `json:"@type"`
	Name string `json:"ub:name"`
}

type Publication struct {
	ID                UID    `json:"@id"`
	Type              UID    `json:"@type"`
	Name              string `json:"ub:name"`
	PublicationAuthor []IRI  `json:"ub:publicationAuthor"`
}

type ResearchGroup struct {
	ID                UID `json:"@id"`
	Type              UID `json:"@type"`
	SubOrganizationOf IRI `json:"ub:subOrganizationOf"`
}
