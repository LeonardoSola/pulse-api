package config

import "testing"

func TestBuildRabbitMQURL(t *testing.T) {
	tests := []struct {
		name  string
		user  string
		pass  string
		vhost string
		want  string
	}{
		{name: "default vhost", user: "guest", pass: "guest", vhost: "/", want: "amqp://guest:guest@localhost:5672/"},
		{name: "named vhost", user: "guest", pass: "guest", vhost: "prod", want: "amqp://guest:guest@localhost:5672/prod"},
		{name: "escaped credentials", user: "us@r", pass: "p:ss/w", vhost: "/", want: "amqp://us%40r:p%3Ass%2Fw@localhost:5672/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildRabbitMQURL("localhost", "5672", tt.user, tt.pass, tt.vhost)
			if got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLoadRabbitMQEnvironmentDefaults(t *testing.T) {
	for _, key := range []string{"RABBITMQ_ENABLED", "RABBITMQ_URL", "RABBITMQ_HOST", "RABBITMQ_PORT", "RABBITMQ_USER", "RABBITMQ_PASSWORD", "RABBITMQ_VHOST", "RABBITMQ_EXCHANGE"} {
		t.Setenv(key, "")
	}

	cfg, err := LoadRabbitMQEnvironment()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Enabled {
		t.Error("RabbitMQ should be disabled by default")
	}
	if cfg.URL != "amqp://guest:guest@localhost:5672/" {
		t.Errorf("URL = %q", cfg.URL)
	}
	if cfg.Exchange != "pulse-api" {
		t.Errorf("Exchange = %q", cfg.Exchange)
	}
}

func TestLoadRabbitMQEnvironmentURLTakesPrecedence(t *testing.T) {
	t.Setenv("RABBITMQ_ENABLED", "TRUE")
	t.Setenv("RABBITMQ_URL", "amqp://u:p@broker:5672/")
	t.Setenv("RABBITMQ_HOST", "ignored")

	cfg, err := LoadRabbitMQEnvironment()
	if err != nil {
		t.Fatal(err)
	}

	if !cfg.Enabled {
		t.Error("RABBITMQ_ENABLED=TRUE should enable RabbitMQ")
	}
	if cfg.URL != "amqp://u:p@broker:5672/" {
		t.Errorf("URL = %q", cfg.URL)
	}
}

func TestLoadServerEnvironment(t *testing.T) {
	t.Setenv("API_KEY", "")
	if _, err := LoadServerEnvironment(); err == nil {
		t.Error("expected an error when API_KEY is empty")
	}

	t.Setenv("API_KEY", "secret")
	t.Setenv("SERVER_PORT", "")
	cfg, err := LoadServerEnvironment()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != 8080 {
		t.Errorf("default port = %d, want 8080", cfg.Port)
	}

	t.Setenv("SERVER_PORT", "9000")
	cfg, _ = LoadServerEnvironment()
	if cfg.Port != 9000 {
		t.Errorf("port = %d, want 9000", cfg.Port)
	}
}

func TestLoadDatabaseEnvironment(t *testing.T) {
	t.Setenv("DB_CONNECTION_STRING", "")
	t.Setenv("DB_HOST", "db")
	t.Setenv("DB_PORT", "5432")
	t.Setenv("DB_USER", "pulse")
	t.Setenv("DB_PASS", "pw")
	t.Setenv("DB_NAME", "pulse_api")
	t.Setenv("DB_SSL_MODE", "disable")

	cfg, _ := LoadDatabaseEnvironment()
	want := "host=db port=5432 user=pulse password=pw dbname=pulse_api sslmode=disable"
	if cfg.ConnectionString != want {
		t.Errorf("got %q, want %q", cfg.ConnectionString, want)
	}

	t.Setenv("DB_CONNECTION_STRING", "postgres://x")
	cfg, _ = LoadDatabaseEnvironment()
	if cfg.ConnectionString != "postgres://x" {
		t.Errorf("DB_CONNECTION_STRING should take precedence, got %q", cfg.ConnectionString)
	}
}
