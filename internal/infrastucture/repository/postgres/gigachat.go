package postgres

import (
	"bot/internal/application/dto"
	pu "bot/pkg/postgres_utils"
	"context"
)

func (repo *PostgresRepository) AddNewDialogMessage(ctx context.Context, telegram_id, project_id int64, content []byte) error {
	sql := `
		UPDATE telegram_dialogs
		SET dialog = telegram_dialogs.dialog || $1::jsonb
		WHERE telegram_id = $2::bigint;
	`
	if _, err := repo.db.ExecContext(ctx, sql, content, telegram_id); err != nil {
		return err
	}
	return nil
}

func (repo *PostgresRepository) GetAllDialogMessages(ctx context.Context, telegram_id, project_id int64) (string, error) {
	sql := `
		SELECT t.dialog::text
		FROM telegram_dialogs t
		WHERE t.telegram_id = $1::bigint;
	`
	rows, err := pu.Dispatch[dto.GigaChatDbo](
		ctx,
		repo.db,
		sql,
		telegram_id,
	)
	if err != nil {
		return "", err
	}
	row := rows[0]
	return row.Dialog.String, nil
}
