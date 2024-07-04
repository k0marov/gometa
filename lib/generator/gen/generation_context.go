package gen

type GenerationContext struct {
	ModuleName   string
	PackageName  string
	EntityImport string
	EntityName   string
}

func (g GenerationContext) WithPackageName(pkgName string) GenerationContext {
	return GenerationContext{
		PackageName:  pkgName,
		ModuleName:   g.ModuleName,
		EntityImport: g.EntityImport,
		EntityName:   g.EntityName,
	}
}
