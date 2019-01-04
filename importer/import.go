package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/alexashley/terraform-provider-kong/kong/kong"
	"github.com/google/subcommands"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type importCommand struct {
	adminApiUrl    string
	rbacToken      string
	isDryRun       bool
	importFileName string
	tfConfigPath   string
}

var pluginsToResourceImplementations = map[string]string{
	"openid-connect":               "kong_plugin_openid_connect",
	"request-transformer-advanced": "kong_plugin_request_transformer_advanced",
}

func (*importCommand) Name() string {
	return "import"
}

func (*importCommand) Synopsis() string {
	return "Import consumers, services, and routes from Kong."
}

func (*importCommand) Usage() string {
	return `import -admin-api-url=https://kong-admin.foo.com`
}

func (cmd *importCommand) SetFlags(flags *flag.FlagSet) {
	flags.StringVar(
		&cmd.adminApiUrl,
		"admin-api-url",
		"http://localhost:8001",
		"Kong's admin api url. Usually listening on port 8001.",
	)
	flags.StringVar(
		&cmd.rbacToken,
		"rbac-token",
		"",
		"Kong EE RBAC token. Only necessary if your Kong Enterprise installation is secured with RBAC.",
	)
	flags.BoolVar(
		&cmd.isDryRun,
		"dry-run",
		false,
		"List the resources that will be imported, but do not actually import them.",
	)
	flags.StringVar(
		&cmd.importFileName,
		"state",
		"import-state.json",
		"Holds the current import state and any exclusions",
	)
	flags.StringVar(
		&cmd.tfConfigPath,
		"tf-config",
		"",
		"Path to Terraform config directory",
	)
}

func (cmd *importCommand) Execute(_ context.Context, flags *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	fmt.Println("Importing resources from: " + cmd.adminApiUrl)

	client, err := kong.NewKongClient(kong.KongConfig{
		AdminApiUrl: cmd.adminApiUrl,
		RbacToken:   cmd.rbacToken,
	})

	if err != nil {
		fmt.Printf("error initializing Kong client: %v\n", err)
		return subcommands.ExitFailure
	}

	state := &kongState{
		client: client,
	}

	if err := state.loadState(cmd.importFileName); err != nil {
		fmt.Printf("error loading import state file %v\n", err)
		return subcommands.ExitFailure
	}

	if err := state.discover(); err != nil {
		fmt.Printf("error while discovering resources %v\n", err)
		return subcommands.ExitFailure
	}

	fmt.Println("\nDiscovery:")
	fmt.Println(state.discoveryReport())

	if !cmd.isDryRun {
		if err := state.importResources(cmd); err != nil {
			fmt.Printf("error occurred while importing resources: %v\n", err)
			err := state.finish(cmd.importFileName)

			if err != nil {
				fmt.Println("Additional error saving progress ", err)
			}

			return subcommands.ExitFailure
		} else {
			if err := state.finish(cmd.importFileName); err != nil {
				fmt.Printf("Error occurred saving import file %v\n", err)
			}
		}
	}

	return subcommands.ExitSuccess
}

type kongState struct {
	services  []kong.KongService
	routes    []kong.KongRoute
	plugins   []kong.KongPlugin
	consumers []kong.KongConsumer

	imports map[string][]string // { services: [<uuid>], routes: [<uuid>,] }

	client *kong.KongClient
}

func (s *kongState) loadState(fileName string) error {
	if stateFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755); err != nil {
		return err
	} else {
		defer stateFile.Close()
		s.imports = make(map[string][]string)
		raw, err := ioutil.ReadAll(stateFile)

		if err != nil {
			return err
		}

		if len(raw) == 0 {
			s.imports["services"] = make([]string, 0)
			s.imports["routes"] = make([]string, 0)
			s.imports["plugins"] = make([]string, 0)
			s.imports["consumers"] = make([]string, 0)
		} else if err := json.Unmarshal(raw, &s.imports); err != nil {
			return err
		}
	}

	return nil
}

func (s *kongState) discover() error {
	if consumers, err := s.client.GetConsumers(); err != nil {
		return err
	} else {
		s.consumers = consumers
	}

	if services, err := s.client.GetServices(); err != nil {
		return err
	} else {
		s.services = services
	}

	if routes, err := s.client.GetRoutes(); err != nil {
		return err
	} else {
		s.routes = routes
	}

	if plugins, err := s.client.GetPlugins(); err != nil {
		return err
	} else {
		s.plugins = plugins
	}

	return nil
}

func (s *kongState) discoveryReport() string {
	lines := make([]string, 0)

	if len(s.consumers) > 0 {
		lines = append(lines, fmt.Sprintf("Discovered %d consumers", len(s.consumers)))
	} else {
		lines = append(lines, "No consumers discovered.")
	}

	if len(s.services) > 0 {
		lines = append(lines, fmt.Sprintf("Discovered %d services", len(s.services)))
	} else {
		lines = append(lines, "No services discovered.")
	}

	if len(s.routes) > 0 {
		lines = append(lines, fmt.Sprintf("Discovered %d routes", len(s.routes)))
	} else {
		lines = append(lines, "No routes discovered.")
	}

	if len(s.plugins) > 0 {
		lines = append(lines, fmt.Sprintf("Discovered %d plugins", len(s.plugins)))
	} else {
		lines = append(lines, "No plugins discovered.")
	}

	for index, line := range lines {
		lines[index] = fmt.Sprintf("- %s", line)
	}

	return strings.Join(lines, "\n")
}

func createHclSafeName(name string) string {
	invalid := []string{"-", "/", " ", "."}
	hclName := name

	for _, c := range invalid {
		hclName = strings.Replace(hclName, c, "_", -1)
	}

	return hclName
}

type resourceImport struct {
	kongResourceType      string // plugin, route, service, consumer
	terraformResourceType string // kong_plugin, kong_route, kong_plugin_openid_connect
	resourceName          string // what's to the right of the terraformResourceType in the HCL
	resourceId            string
	dryRun                bool
	configPath            string // HCL config
}

func (s *kongState) hasResourceBeenImported(resource *resourceImport) bool {
	resourceTypePluralized := resource.kongResourceType + "s"
	if importedIds, ok := s.imports[resourceTypePluralized]; ok {
		for _, id := range importedIds {
			if id == resource.resourceId {
				return true
			}
		}
	}

	return false
}

func (s *kongState) importResource(resourceImport *resourceImport) error {
	if s.hasResourceBeenImported(resourceImport) {
		return nil
	}
	terraformResourceName := fmt.Sprintf("%s.%s", resourceImport.terraformResourceType, createHclSafeName(resourceImport.resourceName))

	if !resourceImport.dryRun {
		// ex: terraform import -config=examples/import kong_service.service_to_import e86f981e-a580-4bd6-aef3-1324adfcc12c
		cmd := exec.Command(
			"terraform",
			"import",
			fmt.Sprintf("-config=%s", resourceImport.configPath),
			terraformResourceName,
			resourceImport.resourceId,
		)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(stderr.String())
			return err
		}
	}

	fmt.Println("Imported:", terraformResourceName)

	return nil
}

func (s *kongState) importResources(cmd *importCommand) error {
	if len(s.consumers) > 0 {
		fmt.Println("\nImporting consumers:")

		for _, consumer := range s.consumers {
			resource := &resourceImport{
				kongResourceType:      "consumer",
				terraformResourceType: "kong_consumer",
				resourceName:          getResourceNameForConsumer(&consumer),
				resourceId:            consumer.Id,
				dryRun:                cmd.isDryRun,
				configPath:            cmd.tfConfigPath,
			}
			if err := s.importResource(resource); err != nil {
				return err
			} else {
				s.imports["consumers"] = append(s.imports["consumers"], consumer.Id)
			}
		}

		if len(s.services) > 0 {
			fmt.Println("\nImporting services:")

			for _, service := range s.services {
				resource := &resourceImport{
					kongResourceType:      "service",
					terraformResourceType: "kong_service",
					resourceName:          service.Name,
					resourceId:            service.Id,
					dryRun:                cmd.isDryRun,
					configPath:            cmd.tfConfigPath,
				}

				if err := s.importResource(resource); err != nil {
					return err
				} else {
					s.imports["services"] = append(s.imports["services"], service.Id)
				}
			}
		}

		if len(s.routes) > 0 {
			fmt.Println("\nImporting routes:")

			for _, route := range s.routes {
				resource := &resourceImport{
					kongResourceType:      "route",
					terraformResourceType: "kong_route",
					resourceName:          getResourceNameForRoute(s, &route),
					resourceId:            route.Id,
					dryRun:                cmd.isDryRun,
					configPath:            cmd.tfConfigPath,
				}

				if err := s.importResource(resource); err != nil {
					return err
				} else {
					s.imports["routes"] = append(s.imports["routes"], route.Id)
				}
			}
		}

		if len(s.plugins) > 0 {
			fmt.Println("\nImporting plugins:")

			for _, plugin := range s.plugins {
				terraformResourceType := "kong_plugin"
				if pluginResourceImplementation, ok := pluginsToResourceImplementations[plugin.Name]; ok {
					terraformResourceType = pluginResourceImplementation
				}

				resource := &resourceImport{
					kongResourceType:      "plugin",
					terraformResourceType: terraformResourceType,
					resourceName:          getResourceNameForPlugin(s, &plugin),
					resourceId:            plugin.Id,
					dryRun:                cmd.isDryRun,
					configPath:            cmd.tfConfigPath,
				}

				if err := s.importResource(resource); err != nil {
					return err
				} else {
					s.imports["plugins"] = append(s.imports["plugins"], plugin.Id)
				}
			}
		}
	}

	return nil
}

func (s *kongState) finish(fileName string) error {
	if data, err := json.Marshal(s.imports); err != nil {
		return err
	} else {
		return ioutil.WriteFile(fileName, data, 0644)
	}
}

func getResourceNameForConsumer(consumer *kong.KongConsumer) string {
	if len(consumer.Username) > 0 {
		return consumer.Username
	} else {
		return consumer.CustomId
	}
}

func getResourceNameForRoute(s *kongState, route *kong.KongRoute) string {
	var service kong.KongService
	for _, s := range s.services {
		if s.Id == route.Service.Id {
			service = s
			break
		}
	}
	name := service.Name

	// TODO: the path/host slices should probably be sorted...
	if len(route.Paths) > 0 {
		path := strings.Split(route.Paths[0], "/")[1:] // need to remove the trailing space from splitting /path
		for index, p := range path {
			// if the path was prefixed with the service name, we don't want to repeat it
			// e.g., service name: products, route path: /products
			// the result should be products, not products_products
			if index == 0 && p == name {
				continue
			}

			name = name + "_" + p
		}
	} else {
		name = name + route.Hosts[0]
	}

	return name
}

func getResourceNameForPlugin(s *kongState, plugin *kong.KongPlugin) string {
	namePrefix := ""

	if plugin.ServiceId != "" {
		for _, service := range s.services {
			if service.Id == plugin.ServiceId {
				namePrefix = service.Name
				break
			}
		}
	} else if plugin.RouteId != "" {
		for _, route := range s.routes {
			if route.Id == plugin.RouteId {
				namePrefix = getResourceNameForRoute(s, &route)
				break
			}
		}
	} else if plugin.ConsumerId != "" {
		for _, consumer := range s.consumers {
			if consumer.Id == plugin.ConsumerId {
				namePrefix = getResourceNameForConsumer(&consumer)
				break
			}
		}
	} else {
		namePrefix = "global"
	}

	// for plugins with specific resource implementations (like openid-connect) we don't want to add the plugin name at the end
	// it's redundant. compare:
	// kong_plugin_openid_connect.foo_openid_connect
	// vs
	// kong_plugin_openid_connect.foo
	if pluginHasSpecificResourceImplementation(plugin) {
		return namePrefix
	}

	return namePrefix + "_" + plugin.Name
}

func pluginHasSpecificResourceImplementation(plugin *kong.KongPlugin) bool {
	_, ok := pluginsToResourceImplementations[plugin.Name]

	return ok
}
