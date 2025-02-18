package dto

import (
	"bot/internal/core/models"
	"database/sql"
)

type AttachmentDbo struct {
	Id          int64          `db:"id"`
	Name        string         `db:"name"`
	Description sql.NullString `db:"description"`
	ObjectId    string         `db:"object_id"`
	Mimetype    string         `db:"mimetype"`
	UserID      int64          `db:"user_id,omitempty"`
	CreatedAt   sql.NullTime   `db:"created_at,omitempty"`
	UpdateAt    sql.NullTime   `db:"updated_at,omitempty"`
	DeleteAt    sql.NullTime   `db:"deleted_at,omitempty"`
}

func (a *AttachmentDbo) ToValue() models.Attachment {
	return models.Attachment{
		Id:          a.Id,
		Name:        a.Name,
		Description: a.Description.String,
		UserId:      a.UserID,
		ObjectId:    a.ObjectId,
		Mimetype:    a.Mimetype,
	}
}
