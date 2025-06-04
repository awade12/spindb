package cmd

import (
	"fmt"

	"github.com/awade12/spindb/internal/db"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new database instance",
	Long:  `Create and configure a new database instance (PostgreSQL, MySQL, or SQLite)`,
}

var createPostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "Create a PostgreSQL database",
	Long:  `Create and start a PostgreSQL database instance using Docker`,
	RunE:  createPostgres,
}

var createMysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "Create a MySQL database",
	Long:  `Create and start a MySQL database instance using Docker`,
	RunE:  createMysql,
}

var createSqliteCmd = &cobra.Command{
	Use:   "sqlite",
	Short: "Create a SQLite database",
	Long:  `Create a SQLite database file`,
	RunE:  createSqlite,
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createPostgresCmd)
	createCmd.AddCommand(createMysqlCmd)
	createCmd.AddCommand(createSqliteCmd)

	createPostgresCmd.Flags().StringP("name", "n", "", "Database name (required)")
	createPostgresCmd.Flags().StringP("user", "u", "postgres", "Database user")
	createPostgresCmd.Flags().StringP("password", "p", "", "Database password (required)")
	createPostgresCmd.Flags().IntP("port", "", 5432, "Database port")
	createPostgresCmd.Flags().StringP("version", "v", "15", "PostgreSQL version")
	createPostgresCmd.Flags().Bool("public", false, "Make database publicly accessible")
	createPostgresCmd.MarkFlagRequired("name")
	createPostgresCmd.MarkFlagRequired("password")

	createMysqlCmd.Flags().StringP("name", "n", "", "Database name (required)")
	createMysqlCmd.Flags().StringP("user", "u", "root", "Database user")
	createMysqlCmd.Flags().StringP("password", "p", "", "Database password (required)")
	createMysqlCmd.Flags().IntP("port", "", 3306, "Database port")
	createMysqlCmd.Flags().StringP("version", "v", "8.0", "MySQL version")
	createMysqlCmd.MarkFlagRequired("name")
	createMysqlCmd.MarkFlagRequired("password")

	createSqliteCmd.Flags().StringP("file", "f", "", "SQLite database file path (required)")
	createSqliteCmd.MarkFlagRequired("file")
}

func createPostgres(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	port, _ := cmd.Flags().GetInt("port")
	version, _ := cmd.Flags().GetString("version")
	public, _ := cmd.Flags().GetBool("public")

	manager := db.NewManager()
	config := &db.PostgresConfig{
		Name:     name,
		User:     user,
		Password: password,
		Port:     port,
		Version:  version,
		Public:   public,
	}

	fmt.Printf("Creating PostgreSQL database '%s'...\n", name)
	return manager.CreatePostgres(config)
}

func createMysql(cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	user, _ := cmd.Flags().GetString("user")
	password, _ := cmd.Flags().GetString("password")
	port, _ := cmd.Flags().GetInt("port")
	version, _ := cmd.Flags().GetString("version")

	manager := db.NewManager()
	config := &db.MySQLConfig{
		Name:     name,
		User:     user,
		Password: password,
		Port:     port,
		Version:  version,
	}

	fmt.Printf("Creating MySQL database '%s'...\n", name)
	return manager.CreateMySQL(config)
}

func createSqlite(cmd *cobra.Command, args []string) error {
	file, _ := cmd.Flags().GetString("file")

	manager := db.NewManager()
	config := &db.SQLiteConfig{
		FilePath: file,
	}

	fmt.Printf("Creating SQLite database at '%s'...\n", file)
	return manager.CreateSQLite(config)
}
