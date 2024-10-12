package queues

type Job struct {
	JobId  string   `sql:"job_id"`
	Domain string   `sql:"domain"`
	Pages  []string `sql:"pages"`
}
