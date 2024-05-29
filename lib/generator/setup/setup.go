package setup

import (
	"fmt"
	"path/filepath"
)

func AddToDI(goPackageName, entityName, projectDir string, moduleName string) error {
	if err := addToContainerStruct(goPackageName, entityName, moduleName, filepath.Join(projectDir, "internal", "app", "dependencies", "container.go")); err != nil {
		return fmt.Errorf("while adding new entity service to container: %w", err)
	}
	if err := addCtrlToRouter(goPackageName, entityName, moduleName, filepath.Join(projectDir, "internal", "app", "initializers", "router.go")); err != nil {
		return fmt.Errorf("while adding new entity controller to router: %w", err)
	}
	if err := addToApplication(goPackageName, entityName, moduleName, filepath.Join(projectDir, "internal", "app", "application.go")); err != nil {
		return fmt.Errorf("while adding new entity to application: %w", err)
	}
	return nil
}
