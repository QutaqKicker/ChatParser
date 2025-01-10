package constants

const (
	AuditPortEnvName       = "ChatParser_Audit_Port"
	BackupPortEnvName      = "ChatParser_Backup_Port"
	ChatPortEnvName        = "ChatParser_Chat_Port"
	UserPortEnvName        = "ChatParser_User_Port"
	RouterPortEnvName      = "ChatParser_Router_Port"
	KafkaBroker1UrlEnvName = "ChatParser_Kafka_1_Url"
)

const (
	KafkaAuditCreateLogTopicName     = "audit-create-log"
	KafkaUserMessageCounterTopicName = "user-message-counter"
)
