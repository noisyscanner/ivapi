package integration

import "github.com/noisyscanner/gofly/gofly"

var frenchLang = &gofly.Language{
	Id:            1,
	Code:          "fr",
	Lang:          "French",
	Locale:        "fr_FR",
	HasHelpers:    true,
	HasReflexives: true,
}

var frenchTense = &gofly.Tense{
	Id:          1,
	Identifier:  "present",
	DisplayName: "Present",
	Order:       0,
}

var frenchPronoun = &gofly.Pronoun{
	Id:          1,
	Identifier:  "je",
	DisplayName: "Je",
	Order:       0,
}

var frenchVerb = &gofly.Verb{
	Id:                   1,
	Infinitive:           "jour",
	NormalisedInfinitive: "jour",
	English:              "to play",
}

var germanLang = &gofly.Language{
	Id:            2,
	Code:          "de",
	Lang:          "German",
	Locale:        "de_DE",
	HasHelpers:    true,
	HasReflexives: true,
}

var germanTense = &gofly.Tense{
	Id:          2,
	Identifier:  "present",
	DisplayName: "Present",
	Order:       0,
}

var germanPronoun = &gofly.Pronoun{
	Id:          2,
	Identifier:  "ich",
	DisplayName: "Ich",
	Order:       0,
}

var germanVerb = &gofly.Verb{
	Id:                   2,
	Infinitive:           "spielen",
	NormalisedInfinitive: "spielen",
	English:              "to play",
}
