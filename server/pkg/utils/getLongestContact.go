package utils

import "contacttracing/pkg/models/dto"

func GetLongestContact(contacts []dto.Contact) dto.Contact {
	longestContact := dto.Contact{}

	for _, c := range contacts {
		if c.Duration > longestContact.Duration {
			longestContact = c
		}
	}

	return longestContact
}
