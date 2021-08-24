package krong

type Config struct {
	LogFile string `hcl:"log_file"`
	PidFile string `hcl:"pid_file"`
	DBURL   string `hcl:"db_url"`
}
