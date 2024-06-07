package generator

import (
	"fmt"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/generator/delivery"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/generator/entity_struct"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/generator/gen"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/generator/repository"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/generator/service"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/generator/setup"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/helpers"
	"gitlab.sch.ocrv.com.rzd/blockchain/platform/gometa.git/lib/schema"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Generate(schemaPath, projectDir string) error {
	goMod, err := os.Open(filepath.Join(projectDir, "go.mod"))
	if err != nil {
		return fmt.Errorf("while opening project go.mod: %w", err)
	}
	defer goMod.Close()
	moduleName, err := GetGoModuleName(goMod)
	if err != nil {
		return fmt.Errorf("getting go module name: %w", err)
	}

	schemaPath, _ = filepath.Abs(schemaPath)

	ent, err := schema.Parse(schemaPath)
	if err != nil {
		return fmt.Errorf("failed parsing schema: %w", err)
	}
	internalDir := filepath.Join(projectDir, "internal")

	entityFile, err := helpers.CreateFileRecursively(filepath.Join(internalDir, "models", fmt.Sprintf("%s.go", ent.JsonName)))
	defer entityFile.Close()
	if err != nil {
		return err
	}
	if err := entity_struct.Generate(entityFile, ent, "models"); err != nil {
		return fmt.Errorf("generating model file: %w", err)
	}

	entityImportPath := fmt.Sprintf("%s/internal/models", moduleName)

	baseGenCtx := gen.GenerationContext{
		ModuleName:   moduleName,
		EntityImport: entityImportPath,
		EntityName:   ent.Name,
	}
	goPackageName := ent.JsonName

	repoModelsFile, err := helpers.CreateFileRecursively(filepath.Join(internalDir, "repository", goPackageName, "models.go"))
	defer repoModelsFile.Close()
	if err != nil {
		return err
	}
	if err := repository.GenerateModel(repoModelsFile, ent, baseGenCtx.WithPackageName(goPackageName)); err != nil {
		return fmt.Errorf("generating repo models.go file: %w", err)
	}
	repoFile, err := helpers.CreateFileRecursively(filepath.Join(internalDir, "repository", goPackageName, "repository.go"))
	defer repoFile.Close()
	if err != nil {
		return err
	}
	if err := repository.Generate(repoFile, baseGenCtx.WithPackageName(goPackageName)); err != nil {
		return fmt.Errorf("generating model file: %w", err)
	}

	serviceFile, err := helpers.CreateFileRecursively(filepath.Join(internalDir, "services", goPackageName, "service.go"))
	defer serviceFile.Close()
	if err != nil {
		return err
	}
	if err := service.Generate(serviceFile, baseGenCtx.WithPackageName(goPackageName)); err != nil {
		return fmt.Errorf("generating service file: %w", err)
	}

	handlersMappersFile, err := helpers.CreateFileRecursively(filepath.Join(internalDir, "web", "controllers", "apiv1", goPackageName, "mappers.go"))
	defer handlersMappersFile.Close()
	if err != nil {
		return err
	}
	if err := delivery.GenerateMappers(handlersMappersFile, ent, baseGenCtx.WithPackageName(goPackageName)); err != nil {
		return fmt.Errorf("generating mappers for json: %w", err)
	}
	handlersFile, err := helpers.CreateFileRecursively(filepath.Join(internalDir, "web", "controllers", "apiv1", goPackageName, "controller.go"))
	defer handlersFile.Close()
	if err != nil {
		return err
	}
	if err := delivery.GenerateHandlers(handlersFile, ent, baseGenCtx.WithPackageName(goPackageName)); err != nil {
		return fmt.Errorf("generating handlers layer: %w", err)
	}

	if err := setup.AddToDI(goPackageName, ent.Name, projectDir, moduleName); err != nil {
		return fmt.Errorf("adding new crud to DI: %w", err)
	}
	return nil
}

func GetGoModuleName(goModContents io.Reader) (string, error) {
	contents, err := io.ReadAll(goModContents)
	if err != nil {
		return "", fmt.Errorf("while reading from go.mod: %w", err)
	}
	lines := strings.Split(string(contents), "\n")
	if len(lines) == 0 {
		return "", fmt.Errorf("got empty contents")
	}
	firstLine := strings.Split(lines[0], " ")
	if len(firstLine) != 2 {
		return "", fmt.Errorf("invalid first go.mod line: %q", lines[0])
	}
	return firstLine[1], nil
}
