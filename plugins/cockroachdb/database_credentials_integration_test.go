package cockroachdb

import (
	"context"
	"maps"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/1Password/shell-plugins/sdk"
	"github.com/1Password/shell-plugins/sdk/schema/fieldname"
)

// Run with COCKROACHDB_INTEGRATION=1 make test. Docker pulls the pinned image if needed.
func TestDatabaseCredentialsIntegration(t *testing.T) {
	if os.Getenv("COCKROACHDB_INTEGRATION") != "1" {
		t.Skip("Set COCKROACHDB_INTEGRATION=1 to test against CockroachDB in Docker")
	}

	for _, secure := range []bool{false, true} {
		name := "insecure"
		if secure {
			name = "TLS"
		}
		t.Run(name, func(t *testing.T) {
			start := "exec cockroach start-single-node --listen-addr=localhost:26258 --store=type=mem,size=1GiB --cache=64MiB --max-sql-memory=64MiB "
			adminArgs := []string{"--host=localhost:26258"}
			if secure {
				start = "mkdir /certs && cockroach cert create-ca --certs-dir=/certs --ca-key=/certs/ca.key && " +
					"cockroach cert create-node localhost --certs-dir=/certs --ca-key=/certs/ca.key && " +
					"cockroach cert create-client root --certs-dir=/certs --ca-key=/certs/ca.key && " + start + "--certs-dir=/certs"
				adminArgs = append(adminArgs, "--certs-dir=/certs")
			} else {
				start += "--insecure"
				adminArgs = append(adminArgs, "--insecure")
			}
			container, err := cockroachDocker("run", "--rm", "--quiet", "--detach", "--entrypoint=sh", "cockroachdb/cockroach:v25.2.4", "-c", start)
			if err != nil {
				t.Fatalf("Start CockroachDB: %v\n%s", err, container)
			}
			container = strings.TrimSpace(container)
			t.Cleanup(func() {
				if output, err := cockroachDocker("rm", "--force", container); err != nil {
					t.Errorf("Remove test container: %v\n%s", err, output)
				}
			})
			admin := append([]string{"exec", container, "cockroach", "sql"}, adminArgs...)
			deadline := time.Now().Add(90 * time.Second)
			for {
				output, err := cockroachDocker(append(admin, "--execute=SELECT 1")...)
				if err == nil {
					break
				}
				running, inspectErr := cockroachDocker("inspect", "--format={{.State.Running}}", container)
				if inspectErr != nil || strings.TrimSpace(running) != "true" || time.Now().After(deadline) {
					logs, _ := cockroachDocker("logs", container)
					t.Fatalf("CockroachDB did not become ready: %v\n%s\n%s", err, output, logs)
				}
				time.Sleep(time.Second)
			}
			setup := "CREATE DATABASE plugin_test; CREATE USER shell_plugin; GRANT ALL ON DATABASE plugin_test TO shell_plugin;"
			if secure {
				setup += "ALTER USER shell_plugin WITH PASSWORD 'integration-password';"
			}
			if output, err := cockroachDocker(append(admin, "--execute="+setup)...); err != nil {
				t.Fatalf("Create test credentials: %v\n%s", err, output)
			}

			fields := map[sdk.FieldName]string{
				fieldname.Host: "localhost", fieldname.Port: "26258",
				fieldname.User: "shell_plugin", fieldname.Database: "plugin_test",
				"Insecure": "1",
			}
			if secure {
				fields["Insecure"] = "0"
				fields[fieldname.Password] = "integration-password"
			}
			for name := range defaultEnvVarMapping {
				t.Setenv(name, "")
			}
			for name, field := range defaultEnvVarMapping {
				t.Setenv(name, fields[field])
			}
			var imported sdk.ImportOutput
			DatabaseCredentials().Importer(context.Background(), sdk.ImportInput{}, &imported)
			candidates := imported.AllCandidates()
			if len(imported.Errors()) != 0 || len(candidates) != 1 {
				t.Fatalf("Expected one imported credential without errors, got %d candidates and %d errors", len(candidates), len(imported.Errors()))
			}
			provision := func(t *testing.T, fields map[sdk.FieldName]string) map[string]string {
				t.Helper()
				out := sdk.ProvisionOutput{Environment: make(map[string]string)}
				DatabaseCredentials().DefaultProvisioner.Provision(context.Background(), sdk.ProvisionInput{ItemFields: fields}, &out)
				if len(out.Diagnostics.Errors) != 0 {
					t.Fatal("Credential provisioning failed")
				}
				return out.Environment
			}
			query := func(environment map[string]string) (string, error) {
				args := []string{"exec"}
				for name, value := range environment {
					args = append(args, "--env", name+"="+value)
				}
				args = append(args, container, "cockroach", "sql", "--format=csv", "--execute=SELECT current_user, current_database()")
				if secure {
					args = append(args, "--certs-dir=/certs")
				}
				return cockroachDocker(args...)
			}
			checkQuery := func(t *testing.T, fields map[sdk.FieldName]string) {
				t.Helper()
				output, err := query(provision(t, fields))
				if err != nil || !strings.Contains(output, "\nshell_plugin,plugin_test\n") {
					t.Fatalf("Query with imported and provisioned credentials: %v\n%s", err, output)
				}
			}
			t.Run("imported_credentials_connect", func(t *testing.T) {
				checkQuery(t, candidates[0].Fields)
			})
			if secure {
				t.Run("TLS_is_the_default", func(t *testing.T) {
					fields := maps.Clone(candidates[0].Fields)
					delete(fields, "Insecure")
					checkQuery(t, fields)
				})
				checkRejected := func(t *testing.T, environment map[string]string) {
					t.Helper()
					output, err := query(environment)
					if err == nil || !strings.Contains(output, "password authentication failed") {
						t.Fatalf("Expected password authentication failure: %v\n%s", err, output)
					}
				}
				t.Run("missing_password_is_rejected", func(t *testing.T) {
					fields := maps.Clone(candidates[0].Fields)
					delete(fields, fieldname.Password)
					checkRejected(t, provision(t, fields))
				})
				t.Run("incorrect_password_is_rejected", func(t *testing.T) {
					fields := maps.Clone(candidates[0].Fields)
					fields[fieldname.Password] = "wrong-password"
					checkRejected(t, provision(t, fields))
				})
				t.Run("COCKROACH_PASSWORD_is_ignored", func(t *testing.T) {
					environment := provision(t, candidates[0].Fields)
					delete(environment, "PGPASSWORD")
					environment["COCKROACH_PASSWORD"] = "integration-password"
					checkRejected(t, environment)
				})
			}
		})
	}
}

func cockroachDocker(args ...string) (string, error) {
	timeout := 30 * time.Second
	if args[0] == "run" {
		timeout = 2 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	output, err := exec.CommandContext(ctx, "docker", args...).CombinedOutput()
	return string(output), err
}
