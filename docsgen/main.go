package main

import (
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/provider"
	"github.com/hashicorp/terraform/helper/schema"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"sort"
	"text/template"
)

const (
	OutputDir            = "docs"
	GenDir               = "docsgen"
	ResourceTemplateFile = "resource.tmpl"
	IndexTemplateFile    = "index.tmpl"
	SourceUrl            = "https://github.com/alexashley/terraform-provider-kong/tree/master/kong/provider"
)

type Attribute struct {
	Type          string
	Description   string
	Name          string
	ConflictsWith []string
	Default       interface{}
	Required      bool
}

type Resource struct {
	Name      string
	CanImport bool

	UserProvidedAttributes []Attribute
	ComputedAttributes     []Attribute
	Src                    string
}

type ResourceMetadata struct {
	Description string `yaml:"description"`
	Example     string `yaml:"example"`
	Import      string `yaml:"import"`
}

type ResourceTemplateModel struct {
	Resource Resource
	Meta     ResourceMetadata
}

type IndexTemplateModel struct {
	Title     string
	Resources []string
}

func getComplexType(attribute *schema.Schema) string {
	if attribute.Elem != nil {
		if elemSchema, ok := attribute.Elem.(*schema.Schema); ok {
			elemType := mapAttributeType(elemSchema)

			return elemType
		} else {
			return "*schema.Resource"
		}
	} else {
		return "string"
	}
}

func mapAttributeType(attribute *schema.Schema) string {
	typeString := "Unknown"

	switch attribute.Type {
	case schema.TypeString:
		typeString = "string"
	case schema.TypeInt:
		typeString = "int"
	case schema.TypeBool:
		typeString = "bool"
	case schema.TypeFloat:
		typeString = "float"
	case schema.TypeList:
		typeString = fmt.Sprintf("[]%s", getComplexType(attribute))
	case schema.TypeSet:
		typeString = fmt.Sprintf("set[%s]", getComplexType(attribute))
	case schema.TypeMap:
		typeString = fmt.Sprintf("map[string][%s]", getComplexType(attribute))
	}

	return typeString
}

func mapDefaultValue(defaultValue interface{}) interface{} {
	if defaultValue == nil {
		return ""
	}

	return defaultValue
}

func attributeSorter(attributes []Attribute) func(i, j int) bool {
	return func(i, j int) bool {
		a1 := attributes[i]
		a2 := attributes[j]

		// float required items to the top,
		// then sort alphabetically, descending
		if a1.Required && a2.Required {
			return a1.Name < a2.Name
		} else if a1.Required {
			return true
		} else if a2.Required {
			return false
		}

		return a1.Name < a2.Name
	}
}

func readResourcesFromProvider() (resources []Resource, warnings []string) {
	terraformProviderResources := provider.KongProvider().ResourcesMap

	for terraformResourceName, terraformResource := range terraformProviderResources {
		resource := Resource{
			Name:      terraformResourceName,
			CanImport: terraformResource.Importer != nil,
		}

		for attributeName, attributeSchema := range terraformResource.Schema {
			attribute := Attribute{
				Name:          attributeName,
				ConflictsWith: attributeSchema.ConflictsWith,
				Description:   attributeSchema.Description,
				Type:          mapAttributeType(attributeSchema),
				Default:       mapDefaultValue(attributeSchema.Default),
				Required:      attributeSchema.Required,
			}

			if attribute.Description == "" {
				warnings = append(warnings, fmt.Sprintf("%s.%s is missing a description.\n", terraformResourceName, attributeName))
			}

			// this condition can't be inverted to !attributeSchema.Computed
			// there are optional fields that may require computation if not provided by the user
			if attributeSchema.Required || attributeSchema.Optional {
				resource.UserProvidedAttributes = append(resource.UserProvidedAttributes, attribute)
			} else {
				resource.ComputedAttributes = append(resource.ComputedAttributes, attribute)
			}
		}
		resource.Src = fmt.Sprintf("%s/resource_%s.go", SourceUrl, terraformResourceName)

		sort.Slice(resource.UserProvidedAttributes, attributeSorter(resource.UserProvidedAttributes))
		sort.Slice(resource.ComputedAttributes, attributeSorter(resource.ComputedAttributes))

		resources = append(resources, resource)
	}

	return resources, warnings
}

func generateDocs(resources []Resource) error {
	if _, err := os.Stat(OutputDir); os.IsNotExist(err) {
		os.Mkdir(OutputDir, os.ModePerm)
	}

	index := IndexTemplateModel{
		Title: "`terraform-provider-kong`",
	}

	for _, r := range resources {
		file, err := os.Create(fmt.Sprintf("%s/%s.md", OutputDir, r.Name))
		if err != nil {
			return err
		}

		metaFile, err := ioutil.ReadFile(fmt.Sprintf("%s/resource-metadata/%s.yml", GenDir, r.Name))

		if err != nil {
			return err
		}

		var resourceMeta ResourceMetadata
		err = yaml.Unmarshal(metaFile, &resourceMeta)

		if err != nil {
			return err
		}

		t := template.Must(template.ParseFiles(fmt.Sprintf("%s/%s", GenDir, ResourceTemplateFile)))
		err = t.ExecuteTemplate(file, "base", ResourceTemplateModel{
			Resource: r,
			Meta:     resourceMeta,
		})

		if err != nil {
			return err
		}

		file.Close()

		index.Resources = append(index.Resources, r.Name)
	}

	sort.Slice(index.Resources, func(i, j int) bool {
		return index.Resources[i] < index.Resources[j]
	})
	indexFile, err := os.Create(fmt.Sprintf("%s/index.md", OutputDir))

	if err != nil {
		return err
	}

	defer indexFile.Close()

	t := template.Must(template.ParseFiles(fmt.Sprintf("%s/%s", GenDir, IndexTemplateFile)))
	if err := t.ExecuteTemplate(indexFile, "index", index); err != nil {
		return err
	}

	return nil
}

func main() {
	resources, warnings := readResourcesFromProvider()

	fmt.Printf("Discovered %d resources from provider: \n", len(resources))
	for _, r := range resources {
		fmt.Printf("- %s\n", r.Name)
	}
	if len(warnings) > 0 {
		fmt.Println("\nThe following issues should be addressed in order to generate high-quality documentation:")
	}

	for _, warning := range warnings {
		fmt.Printf("- %s", warning)
	}

	if len(warnings) > 0 {
		os.Exit(1)
	}

	err := generateDocs(resources)

	if err != nil {
		fmt.Printf("Error occured while generating docs %v\n", err)
	}
}
