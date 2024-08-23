package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		exit(0, fmt.Sprintf("Usage: %s <package-org> <component-name>"))
	}

	componentTemplateRoot := "component-template/component"
	packageOrg := os.Args[1]
	componentName := os.Args[2]
	componentDir := filepath.Join("components", componentName)

	_, err := os.Stat(componentDir)
	if err == nil {
		exit(1, fmt.Sprintf("Target component directory already exists: %s", componentDir))
	}
	if err != nil && !os.IsNotExist(err) {
		exit(1, err.Error())
	}

	err = os.MkdirAll(componentDir, 0755)
	if err != nil {
		exit(1, err.Error())
	}

	err = fs.WalkDir(
		os.DirFS(componentTemplateRoot),
		".",
		func(path string, d fs.DirEntry, err error) error {
			srcFilePath := filepath.Join(componentTemplateRoot, path)
			fileInfo, err := os.Stat(srcFilePath)
			if err != nil {
				return err
			}

			if fileInfo.IsDir() {
				return nil
			}

			switch path {
			case "main.go":
				err = generateFile(packageOrg, componentName, srcFilePath, filepath.Join(componentDir, path))
			case "wit/component.wit":
				err = generateFile(packageOrg, componentName, srcFilePath, filepath.Join(componentDir, "wit", componentName+".wit"))
			default:
				err = copyFile(srcFilePath, filepath.Join(componentDir, path))
			}
			if err != nil {
				return err
			}

			return nil
		})
	if err != nil {
		exit(1, err.Error())
	}
}

func generateFile(packageOrg, componentName, srcFileName, dstFileName string) error {
	pascalPackageOrg := dashToPascal(packageOrg)
	pascalComponentName := dashToPascal(componentName)

	fmt.Printf("Generating from %s to %s\n", srcFileName, dstFileName)

	contentsBs, err := os.ReadFile(srcFileName)
	if err != nil {
		return err
	}

	contents := string(contentsBs)

	contents = strings.ReplaceAll(contents, "component-name", componentName)
	contents = strings.ReplaceAll(contents, "package-org", packageOrg)
	contents = strings.ReplaceAll(contents, "ComponentName", pascalComponentName)
	contents = strings.ReplaceAll(contents, "PackageOrg", pascalPackageOrg)

	err = os.MkdirAll(filepath.Dir(dstFileName), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(dstFileName, []byte(contents), 0644)
	if err != nil {
		return err
	}

	return nil
}

func dashToPascal(s string) string {
	parts := strings.Split(s, "-")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(string(part[0])) + part[1:]
		}
	}
	return strings.Join(parts, "")
}

func copyFile(srcFileName, dstFileName string) error {
	fmt.Printf("Copy %s to %s\n", srcFileName, dstFileName)

	src, err := os.Open(srcFileName)
	if err != nil {
		return err
	}
	defer func() { _ = src.Close() }()

	err = os.MkdirAll(filepath.Dir(dstFileName), 0755)
	if err != nil {
		return err
	}

	dst, err := os.Create(dstFileName)
	if err != nil {
		return err
	}
	defer func() { _ = dst.Close() }()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}

	return nil
}

func exit(code int, message string) {
	fmt.Println(message)
	os.Exit(code)
}
