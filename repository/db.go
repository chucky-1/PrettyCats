package repository

//func RequestDB(pool pgxpool.Pool) *pgxpool.Pool {
//	if err := initConfig(); err != nil {
//		log.Fatal("error config files")
//	}
//
//	if err := godotenv.Load(); err != nil {
//		log.Fatal("error loading env variables")
//	}
//
//	url := fmt.Sprintf("%s://%s:%s@%s/%s",
//		viper.GetString("db.pos"),
//		viper.GetString("db.username"),
//		os.Getenv("DB_PASSWORD"),
//		viper.GetString("db.host"),
//		viper.GetString("db.dbase"))
//
//	conn, err := pool.Connect(context.Background(), url)
//	if err != nil {
//		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
//		os.Exit(1)
//	}
//	return conn
//}

//func initConfig() error {
//	viper.AddConfigPath("configs")
//	viper.SetConfigName("config")
//	return viper.ReadInConfig()
//}
