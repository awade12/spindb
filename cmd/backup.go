package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/awade12/spindb/internal/backup"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Backup and restore database operations",
	Long:  `Create, list, restore and manage database backups`,
}

var backupCreateCmd = &cobra.Command{
	Use:   "create [database-name]",
	Short: "Create a backup of a database",
	Long:  `Create a backup of the specified database`,
	Args:  cobra.ExactArgs(1),
	RunE:  createBackup,
}

var backupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available backups",
	Long:  `List all available database backups`,
	RunE:  listBackups,
}

var backupRestoreCmd = &cobra.Command{
	Use:   "restore [backup-name] [target-database]",
	Short: "Restore a backup to a database",
	Long:  `Restore a backup to the specified target database`,
	Args:  cobra.ExactArgs(2),
	RunE:  restoreBackup,
}

var backupDeleteCmd = &cobra.Command{
	Use:   "delete [backup-name]",
	Short: "Delete a backup",
	Long:  `Delete the specified backup file`,
	Args:  cobra.ExactArgs(1),
	RunE:  deleteBackup,
}

func init() {
	rootCmd.AddCommand(backupCmd)
	backupCmd.AddCommand(backupCreateCmd)
	backupCmd.AddCommand(backupListCmd)
	backupCmd.AddCommand(backupRestoreCmd)
	backupCmd.AddCommand(backupDeleteCmd)

	backupCreateCmd.Flags().Bool("compress", false, "Compress the backup file")
	backupCreateCmd.Flags().Bool("schema-only", false, "Backup schema only (no data)")
	backupCreateCmd.Flags().Bool("data-only", false, "Backup data only (no schema)")
}

func createBackup(cmd *cobra.Command, args []string) error {
	dbName := args[0]

	compress, _ := cmd.Flags().GetBool("compress")
	schemaOnly, _ := cmd.Flags().GetBool("schema-only")
	dataOnly, _ := cmd.Flags().GetBool("data-only")

	if schemaOnly && dataOnly {
		return fmt.Errorf("cannot specify both --schema-only and --data-only")
	}

	options := &backup.BackupOptions{
		Compress:   compress,
		SchemaOnly: schemaOnly,
		DataOnly:   dataOnly,
	}

	manager := backup.NewBackupManager()

	fmt.Printf("Creating backup for database '%s'...\n", dbName)
	backupInfo, err := manager.CreateBackup(dbName, options)
	if err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	fmt.Printf("✅ Backup created successfully!\n")
	fmt.Printf("   Name: %s\n", backupInfo.Name)
	fmt.Printf("   Size: %.2f MB\n", float64(backupInfo.Size)/(1024*1024))
	fmt.Printf("   Path: %s\n", backupInfo.FilePath)
	if backupInfo.Compressed {
		fmt.Printf("   Compressed: Yes\n")
	}

	return nil
}

func listBackups(cmd *cobra.Command, args []string) error {
	manager := backup.NewBackupManager()

	backups, err := manager.ListBackups()
	if err != nil {
		return fmt.Errorf("failed to list backups: %w", err)
	}

	if len(backups) == 0 {
		fmt.Println("No backups found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "NAME\tDATABASE\tTYPE\tSIZE\tCREATED\tCOMPRESSED")
	fmt.Fprintln(w, "----\t--------\t----\t----\t-------\t----------")

	for _, backup := range backups {
		size := fmt.Sprintf("%.2f MB", float64(backup.Size)/(1024*1024))
		compressed := "No"
		if backup.Compressed {
			compressed = "Yes"
		}

		created := backup.CreatedAt.Format("2006-01-02 15:04")

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n",
			backup.Name,
			backup.Database,
			backup.Type,
			size,
			created,
			compressed,
		)
	}

	return w.Flush()
}

func restoreBackup(cmd *cobra.Command, args []string) error {
	backupName := args[0]
	targetDb := args[1]

	if !strings.Contains(backupName, ".") {
		return fmt.Errorf("backup name must include file extension (e.g., backup_name.sql)")
	}

	manager := backup.NewBackupManager()

	fmt.Printf("Restoring backup '%s' to database '%s'...\n", backupName, targetDb)

	if err := manager.RestoreBackup(backupName, targetDb); err != nil {
		return fmt.Errorf("failed to restore backup: %w", err)
	}

	fmt.Printf("✅ Backup restored successfully to '%s'!\n", targetDb)
	return nil
}

func deleteBackup(cmd *cobra.Command, args []string) error {
	backupName := args[0]

	manager := backup.NewBackupManager()

	if err := manager.DeleteBackup(backupName); err != nil {
		return fmt.Errorf("failed to delete backup: %w", err)
	}

	fmt.Printf("✅ Backup '%s' deleted successfully!\n", backupName)
	return nil
}
