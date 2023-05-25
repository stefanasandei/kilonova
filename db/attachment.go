package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/KiloProjects/kilonova"
	"github.com/jackc/pgx/v5"
)

const createAttachmentQuery = "INSERT INTO attachments (visible, private, execable, name, data, last_updated_by) VALUES (?, ?, ?, ?, ?, ?) RETURNING id;"

func (a *DB) CreateAttachment(ctx context.Context, att *kilonova.Attachment, problemID int, data []byte, authorID *int) error {
	if problemID == 0 || data == nil {
		return kilonova.ErrMissingRequired
	}
	if _, err := a.ProblemAttachments(ctx, problemID, &kilonova.AttachmentFilter{Name: &att.Name}); err != nil {
		return kilonova.ErrAttachmentExists
	}

	var id int
	err := a.conn.GetContext(ctx, &id, a.conn.Rebind(createAttachmentQuery), att.Visible, att.Private, att.Exec, att.Name, data, authorID)
	if err != nil {
		return err
	}
	_, err = a.pgconn.Exec(ctx, "INSERT INTO problem_attachments_m2m (problem_id, attachment_id) VALUES ($1, $2)", problemID, id)
	if err != nil {
		return err
	}
	att.ID = id
	return err
}

const selectedAttFields = "id, created_at, last_updated_at, last_updated_by, visible, private, execable, name, data_size" // Make sure to keep this in sync

func (a *DB) Attachment(ctx context.Context, id int) (*kilonova.Attachment, error) {
	var att dbAttachment
	err := a.conn.GetContext(ctx, &att, "SELECT "+selectedAttFields+" FROM attachments WHERE id = $1 LIMIT 1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return internalToAttachment(&att), err
}

func (a *DB) ProblemAttachment(ctx context.Context, problemID, attachmentID int) (*kilonova.Attachment, error) {
	var att dbAttachment
	err := a.conn.GetContext(ctx, &att, "SELECT "+selectedAttFields+" FROM problem_attachments WHERE problem_id = $1 AND id = $2 LIMIT 1", problemID, attachmentID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return internalToAttachment(&att), err
}

func (a *DB) AttachmentByName(ctx context.Context, problemID int, filename string) (*kilonova.Attachment, error) {
	var att dbAttachment
	err := a.conn.GetContext(ctx, &att, "SELECT "+selectedAttFields+" FROM problem_attachments WHERE problem_id = $1 AND name = $2 LIMIT 1", problemID, filename)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return internalToAttachment(&att), err
}

func (a *DB) MarkdownAttachments(ctx context.Context, limit int, offset int) ([][]byte, error) {
	q, _ := a.pgconn.Query(ctx, "SELECT data FROM attachments WHERE name LIKE '%.md%'")
	atts, err := pgx.CollectRows(q, pgx.RowTo[[]byte])
	if err != nil {
		return [][]byte{}, err
	}
	return atts, nil
}

func (a *DB) ProblemAttachments(ctx context.Context, pbid int, filter *kilonova.AttachmentFilter) ([]*kilonova.Attachment, error) {
	var attachments []*dbAttachment
	where, args := attachmentFilterQuery(filter)
	limit, offset := 0, 0
	if filter != nil {
		limit, offset = filter.Limit, filter.Offset
	}
	query := a.conn.Rebind("SELECT " + selectedAttFields + " FROM problem_attachments WHERE problem_id = ? AND " + strings.Join(where, " AND ") + " ORDER BY name ASC " + FormatLimitOffset(limit, offset))
	err := a.conn.SelectContext(ctx, &attachments, query, append([]any{pbid}, args...)...)
	if errors.Is(err, sql.ErrNoRows) {
		return []*kilonova.Attachment{}, nil
	}
	return mapper(attachments, internalToAttachment), err
}

func (a *DB) AttachmentData(ctx context.Context, id int) ([]byte, error) {
	var data []byte
	err := a.conn.GetContext(ctx, &data, "SELECT data FROM attachments WHERE id = $1", id)
	if errors.Is(err, sql.ErrNoRows) {
		return []byte{}, nil
	}
	return data, err
}

func (a *DB) AttachmentDataByName(ctx context.Context, problemID int, name string) ([]byte, error) {
	var data []byte
	err := a.conn.GetContext(ctx, &data, "SELECT data FROM problem_attachments WHERE problem_id = $1 AND name = $2", problemID, name)
	if errors.Is(err, sql.ErrNoRows) {
		return []byte{}, nil
	}
	return data, err
}

func (a *DB) ProblemRawDesc(ctx context.Context, problemID int, name string) ([]byte, bool, error) {
	var rez struct {
		Data    []byte `db:"data"`
		Private bool   `db:"private"`
	}
	err := a.conn.GetContext(ctx, &rez, "SELECT data, private FROM problem_attachments WHERE problem_id = $1 AND name = $2", problemID, name)
	if errors.Is(err, sql.ErrNoRows) {
		return []byte{}, true, nil
	}
	return rez.Data, rez.Private, err
}

const attachmentUpdateStatement = "UPDATE attachments SET %s WHERE id = ?"

func (a *DB) UpdateAttachment(ctx context.Context, id int, upd *kilonova.AttachmentUpdate) error {
	toUpd, args := attachmentUpdateQuery(upd)
	if len(toUpd) == 0 {
		return kilonova.ErrNoUpdates
	}
	args = append(args, id)
	query := a.conn.Rebind(fmt.Sprintf(attachmentUpdateStatement, strings.Join(toUpd, ", ")))
	_, err := a.conn.ExecContext(ctx, query, args...)
	return err
}

func (a *DB) UpdateAttachmentData(ctx context.Context, id int, data []byte, updatedBy *int) error {
	_, err := a.pgconn.Exec(ctx, "UPDATE attachments SET data = $1, last_updated_at = NOW(), last_updated_by = COALESCE($3, last_updated_by) WHERE id = $2", data, id, updatedBy)
	return err
}

func (a *DB) DeleteAttachment(ctx context.Context, attid int) error {
	_, err := a.pgconn.Exec(ctx, "DELETE FROM attachments WHERE id = $1", attid)
	return err
}

func (a *DB) DeleteAttachments(ctx context.Context, pbid int, attIDs []int) (int64, error) {
	result, err := a.pgconn.Exec(ctx, "DELETE FROM attachments WHERE id = ANY($2) AND EXISTS (SELECT 1 FROM problem_attachments_m2m WHERE attachment_id = attachments.id AND problem_id = $1)", pbid, attIDs)
	if err != nil {
		return -1, err
	}

	return result.RowsAffected(), nil
}

func attachmentFilterQuery(filter *kilonova.AttachmentFilter) ([]string, []any) {
	where, args := []string{"1 = 1"}, []any{}
	if filter == nil {
		return where, args
	}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = ?"), append(args, v)
	}
	if v := filter.Name; v != nil {
		where, args = append(where, "name = ?"), append(args, v)
	}
	if v := filter.Visible; v != nil {
		where, args = append(where, "visible = ?"), append(args, v)
	}
	if v := filter.Private; v != nil {
		where, args = append(where, "private = ?"), append(args, v)
	}
	if v := filter.Exec; v != nil {
		where, args = append(where, "execable = ?"), append(args, v)
	}
	return where, args
}

func attachmentUpdateQuery(upd *kilonova.AttachmentUpdate) ([]string, []any) {
	toUpd, args := []string{}, []any{}
	if v := upd.Name; v != nil {
		toUpd, args = append(toUpd, "name = ?"), append(args, v)
	}
	if v := upd.Visible; v != nil {
		toUpd, args = append(toUpd, "visible = ?"), append(args, v)
	}
	if v := upd.Private; v != nil {
		toUpd, args = append(toUpd, "private = ?"), append(args, v)
	}
	if v := upd.Exec; v != nil {
		toUpd, args = append(toUpd, "execable = ?"), append(args, v)
	}
	return toUpd, args
}

type dbAttachment struct {
	ID        int       `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	Visible   bool      `db:"visible"`
	Private   bool      `db:"private"`
	Exec      bool      `db:"execable"`

	LastUpdatedAt time.Time `db:"last_updated_at"`
	LastUpdatedBy *int      `db:"last_updated_by"`

	Name string `db:"name"`
	Size int    `db:"data_size"`
	//Data []byte `db:"data"`
}

func internalToAttachment(att *dbAttachment) *kilonova.Attachment {
	if att == nil {
		return nil
	}
	return &kilonova.Attachment{
		ID:        att.ID,
		CreatedAt: att.CreatedAt,
		Visible:   att.Visible,
		Private:   att.Private,
		Exec:      att.Exec,

		LastUpdatedAt: att.LastUpdatedAt,
		LastUpdatedBy: att.LastUpdatedBy,

		Name: att.Name,
		Size: att.Size,
	}
}
