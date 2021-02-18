package global

const (
	dburi       = "mongodb+srv://saskara:saskara@cluster0.509xu.mongodb.net/<dbname>?retryWrites=true&w=majority"
	dbname      = "rad-blog"
	performance = 100
)

var (
	jwtSecret = []byte("radsecret")
)
