package permission

import "github.com/mikespook/gorbac"

var (
	// GetStamp : スタンプ情報取得権限
	GetStamp = gorbac.NewStdPermission("get_stamp")
	// CreateStamp : スタンプ作成権限
	CreateStamp = gorbac.NewStdPermission("create_stamp")
	// EditStamp : スタンプ編集権限
	EditStamp = gorbac.NewStdPermission("edit_stamp")
	// DeleteStamp : スタンプ削除権限
	DeleteStamp = gorbac.NewStdPermission("delete_stamp")
	// GetMessageStamp : メッセージスタンプ一覧取得権限
	GetMessageStamp = gorbac.NewStdPermission("get_message_stamp")
	// AddMessageStamp : メッセージスタンプ追加権限
	AddMessageStamp = gorbac.NewStdPermission("add_message_stamp")
	// RemoveMessageStamp : メッセージスタンプ削除権限
	RemoveMessageStamp = gorbac.NewStdPermission("remove_message_stamp")
)