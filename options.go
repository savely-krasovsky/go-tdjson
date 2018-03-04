package tdjson

// The Telegram test environment will be used instead of the production environment
func WithTestDC() Option {
	return func(options *options) {
		options.useTestDC = true
	}
}

// The path to the directory for the persistent database; if empty, the current working directory will be used
func WithDatabaseDir(path string) Option {
	return func(options *options) {
		options.databaseDirectory = path
	}
}

// The path to the directory for storing files; if empty, database_directory will be used
func WithFilesDir(path string) Option {
	return func(options *options) {
		options.filesDirectory = path
	}
}

// If set to true, information about downloaded and uploaded files will be saved between application restarts
func WithFileDatabase() Option {
	return func(options *options) {
		options.useFileDatabase = true
	}
}

// If set to true, the library will maintain a cache of users, basic groups, supergroups, channels and secret chats. Implies use WithFileDatabase()
func WithChatInfoDatabase() Option {
	return func(options *options) {
		options.useChatInfoDatabase = true
	}
}

// If set to true, the library will maintain a cache of chats and messages. Implies use WithChatInfoDatabase()
func WithMessageDatabase() Option {
	return func(options *options) {
		options.useMessageDatabase = true
	}
}

// If set to true, support for secret chats will be enabled
func WithSecretChats() Option {
	return func(options *options) {
		options.useSecretChats = true
	}
}

// Application identifier for Telegram API access, which can be obtained at https://my.telegram.org
func WithID(id string) Option {
	return func(options *options) {
		options.apiID = id
	}
}

// Application identifier hash for Telegram API access, which can be obtained at https://my.telegram.org
func WithHash(hash string) Option {
	return func(options *options) {
		options.apiHash = hash
	}
}

// IETF language tag of the user's operating system language
func WithSystemLanguage(lang string) Option {
	return func(options *options) {
		options.systemLanguageCode = lang
	}
}

// Model of the device the application is being run on
func WithDeviceModel(model string) Option {
	return func(options *options) {
		options.deviceModel = model
	}
}

// Version of the operating system the application is being run on
func WithSystemVersion(system string) Option {
	return func(options *options) {
		options.systemVersion = system
	}
}

// Application version
func WithApplicationVersion(version string) Option {
	return func(options *options) {
		options.applicationVersion = version
	}
}

// If set to true, old files will automatically be deleted
func WithStorageOptimizer() Option {
	return func(options *options) {
		options.enableStorageOptimizer = true
	}
}

// If set to true, original file names will be ignored. Otherwise, downloaded files will be saved under names as close as possible to the original name
func WithIgnoreFileNames() Option {
	return func(options *options) {
		options.ignoreFileNames = true
	}
}

// Sets phone number for authorization
func WithPhone(phone string) Option {
	return func(options *options) {
		options.phone = phone
	}
}

// Changes parameters which will be used during execution Auth method with state authorizationStateWaitTdlibParameters.
type Option func(*options)

type options struct {
	useTestDC              bool
	databaseDirectory      string
	filesDirectory         string
	useFileDatabase        bool
	useChatInfoDatabase    bool
	useMessageDatabase     bool
	useSecretChats         bool
	apiID                  string
	apiHash                string
	systemLanguageCode     string
	deviceModel            string
	systemVersion          string
	applicationVersion     string
	enableStorageOptimizer bool
	ignoreFileNames        bool
	phone                  string
}

func (o options) toTdlibParameters() Update {
	return Update{
		"@type":                    "tdlibParameters",
		"use_test_dc":              o.useTestDC,
		"database_directory":       o.databaseDirectory,
		"files_directory":          o.filesDirectory,
		"use_file_database":        o.useFileDatabase,
		"use_chat_info_database":   o.useChatInfoDatabase,
		"use_message_database":     o.useMessageDatabase,
		"use_secret_chats":         o.useSecretChats,
		"api_id":                   o.apiID,
		"api_hash":                 o.apiHash,
		"system_language_code":     o.systemLanguageCode,
		"device_model":             o.deviceModel,
		"system_version":           o.systemVersion,
		"application_version":      o.applicationVersion,
		"enable_storage_optimizer": o.enableStorageOptimizer,
		"ignore_file_names":        o.ignoreFileNames,
	}
}
