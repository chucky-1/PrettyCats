package configs

type Config struct {
	Port string `env:"PORT" envDefault:":8000"`

	PsUsername string `env:"POSTGRES_USERNAME" envDefault:"user"`
	PsPassword string `env:"DB_PASSWORD" envDefault:"testpassword"`
	PsHost     string `env:"POSTGRES_HOST" envDefault:"localhost"`
	PsPort     string `env:"POSTGRES_PORT" envDefault:"5432"`
	PsDBName   string `env:"POSTGRES_DATABASE" envDefault:"postgres"`

	Mongo           string `env:"MONGODB_CONNSTRING" envDefault:"mongodb://root:example@mongo:27017/"`
	MongoUsername   string `env:"MONGO_USERNAME" envDefault:"root"`
	MongoPort       string `env:"MONGODB_PORT" envDefault:"27017"`
	MongoDBName     string `env:"MONGODB_DBNAME" envDefault:"mongodb"`
	MongoCollection string `env:"MONGODB_COLLECTION" envDefault:"mongocl"`

	RedisHost string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort string `env:"REDIS_PORT" envDefault:"6379"`

	CacheHost string `env:"CACHE_HOST" envDefault:"server1"`
	CachePort string `env:"CACHE_PORT" envDefault:"6379"`

	KeyForSignatureJwt string `env:"KEY_FOR_SIGNATURE_JWT" envDefault:"tey3em0n"`
	Salt               string `env:"SALT_FOR_GENERATE_PASSWORD" envDefault:"DSA61SC61CSCS6615"`
}
