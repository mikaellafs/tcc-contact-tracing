package repositories

import (
	"context"
	"database/sql"
	"log"
	"time"

	"contacttracing/src/models/db"
	"contacttracing/src/models/dto"
)

const (
	contactRepositoryLog                 = "[Contact Repository]"
	minDistance                          = 2
	maxDiffTimeToConsiderConstantContact = 20 * time.Minute
)

type PostgreSQLContactRepository struct {
	db *sql.DB
}

func NewPostgreSQLContactRepository(db *sql.DB) *PostgreSQLContactRepository {
	repo := new(PostgreSQLContactRepository)
	repo.db = db
	return repo
}

func (r *PostgreSQLContactRepository) Migrate(ctx context.Context) error {
	log.Println(contactRepositoryLog, "Create contacts table")
	query := `
    CREATE TABLE IF NOT EXISTS contacts(
		id SERIAL PRIMARY KEY,
		userId TEXT NOT NULL,
		anotherUser TEXT NOT NULL,
        firstContactTimestamp TIMESTAMP WITH TIME ZONE NOT NULL,
		lastContactTimestamp TIMESTAMP WITH TIME ZONE NOT NULL,
		distance FLOAT(4) NOT NULL,
		rssi FLOAT(4) NOT NULL,
		batteryLevel FLOAT(4) NOT NULL,
		UNIQUE(userId, anotherUser, firstContactTimestamp, lastContactTimestamp)
    );
    `
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *PostgreSQLContactRepository) Create(ctx context.Context, contact db.Contact) (*db.Contact, error) {
	log.Println(contactRepositoryLog, "Create contact: ", contact)

	var id int64
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO contacts(userId, anotherUser, firstContactTimestamp, lastContactTimestamp, distance, rssi, batteryLevel) 
		VALUES($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT(userId, anotherUser, firstContactTimestamp, lastContactTimestamp)
		DO NOTHING
		RETURNING id`,
		contact.User, contact.AnotherUser, contact.FirstContactTimestamp, contact.LastContactTimestamp,
		contact.Distance, contact.RSSI, contact.BatteryLevel).Scan(&id)

	err = parsePostgreSQLError(err)

	contact.ID = id
	return &contact, err
}

func (r *PostgreSQLContactRepository) GetContactsWithin(ctx context.Context, days int, from time.Time, userId string) ([]dto.Contact, error) {
	log.Println(contactRepositoryLog, "Get contacts from", userId, "within", days, "days, starting from", from)

	rows, err := r.db.QueryContext(ctx, `
	SELECT id, userId, anotherUser, 
		firstContactTimestamp, lastContactTimestamp
	FROM contacts
	WHERE userId = $1
			AND distance <= $2
			AND EXTRACT(EPOCH FROM ($3 - firstContactTimeStamp))/3600 <= $4*24
	ORDER BY (anotherUser, firstContactTimestamp, lastContactTimestamp) ASC
	`, userId, minDistance, from, days)
	if err != nil {
		return nil, err
	}

	contacts := aggregateContactsResult(rows)

	return contacts, nil
}

func (r *PostgreSQLContactRepository) GetContactsBetweenUsers(ctx context.Context, user1, user2 string, from, to time.Time) ([]dto.Contact, error) {
	log.Println(contactRepositoryLog, "Get contacts between", user1, "and", user2, "from", from, "to", to)

	query := `
		SELECT *
		FROM contacts
		WHERE userId = $1 AND anotherUser = $2 AND
			lastContactTimestamp >= $3 AND lastContactTimestamp <= $4
		ORDER BY (anotherUser, firstContactTimestamp, lastContactTimestamp) ASC
	`

	rows, err := r.db.QueryContext(ctx, query, user1, user2, from, to)
	if err != nil {
		return nil, err
	}

	contacts := aggregateContactsResult(rows)

	return contacts, nil
}

func aggregateContactsResult(row *sql.Rows) []dto.Contact {
	var aggregatedContacts []dto.Contact

	var currentContact *db.Contact
	for row.Next() {
		contact := db.Contact{}
		row.Scan(&contact.ID, &contact.User, &contact.AnotherUser,
			&contact.FirstContactTimestamp, &contact.LastContactTimestamp)

		if currentContact == nil {
			currentContact = &contact
			continue
		}

		diffTime := time.Time.Sub(contact.FirstContactTimestamp, currentContact.LastContactTimestamp)

		if currentContact.AnotherUser != contact.AnotherUser || diffTime >= maxDiffTimeToConsiderConstantContact {
			aggregatedContacts = append(aggregatedContacts, dto.Contact{
				User:            currentContact.User,
				DateLastContact: currentContact.LastContactTimestamp,
				AnotherUser:     currentContact.AnotherUser,
				Duration:        time.Time.Sub(currentContact.LastContactTimestamp, currentContact.FirstContactTimestamp),
			})

			currentContact = &contact
			continue
		}

		currentContact.LastContactTimestamp = contact.LastContactTimestamp
	}

	if currentContact != nil {
		aggregatedContacts = append(aggregatedContacts, dto.Contact{
			User:            currentContact.User,
			DateLastContact: currentContact.LastContactTimestamp,
			AnotherUser:     currentContact.AnotherUser,
			Duration:        time.Time.Sub(currentContact.LastContactTimestamp, currentContact.FirstContactTimestamp),
		})
	}

	return aggregatedContacts
}
